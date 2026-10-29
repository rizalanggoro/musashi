package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"

	"musashi/backend/internal/domain"
)

type mockAuthRepo struct {
	findByEmailFn func(email string) (*domain.User, error)
	findByIDFn    func(id uint) (*domain.User, error)
}

func (m *mockAuthRepo) FindByEmail(email string) (*domain.User, error) {
	if m.findByEmailFn != nil {
		return m.findByEmailFn(email)
	}
	return nil, errors.New("not found")
}

func (m *mockAuthRepo) FindByID(id uint) (*domain.User, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, errors.New("not found")
}

func TestAuthService_Login(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	mockUser := &domain.User{
		ID:           1,
		Email:        "admin@musashi.local",
		Name:         "Admin",
		PasswordHash: string(hashedPassword),
		Role:         domain.RoleAdmin,
	}

	t.Run("successful login", func(t *testing.T) {
		repo := &mockAuthRepo{
			findByEmailFn: func(email string) (*domain.User, error) {
				if email == "admin@musashi.local" {
					return mockUser, nil
				}
				return nil, errors.New("not found")
			},
		}

		svc := NewService(repo)
		user, token, err := svc.Login(LoginInput{
			Email:    "admin@musashi.local",
			Password: "secret123",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.ID != 1 || token == "" {
			t.Fatalf("expected user and token, got user: %+v, token: %s", user, token)
		}

		// Token should be parsable
		payload, err := svc.ParseToken(token)
		if err != nil {
			t.Fatalf("failed to parse generated token: %v", err)
		}
		if payload.UserID != 1 || payload.Email != "admin@musashi.local" || payload.Role != domain.RoleAdmin {
			t.Fatalf("mismatched payload: %+v", payload)
		}
	})

	t.Run("user not found returns invalid credentials", func(t *testing.T) {
		repo := &mockAuthRepo{
			findByEmailFn: func(email string) (*domain.User, error) {
				return nil, errors.New("not found")
			},
		}

		svc := NewService(repo)
		_, _, err := svc.Login(LoginInput{
			Email:    "unknown@musashi.local",
			Password: "secret123",
		})
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Fatalf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("wrong password returns invalid credentials", func(t *testing.T) {
		repo := &mockAuthRepo{
			findByEmailFn: func(email string) (*domain.User, error) {
				return mockUser, nil
			},
		}

		svc := NewService(repo)
		_, _, err := svc.Login(LoginInput{
			Email:    "admin@musashi.local",
			Password: "wrongpassword",
		})
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Fatalf("expected ErrInvalidCredentials, got %v", err)
		}
	})
}

func TestAuthService_ParseToken(t *testing.T) {
	svc := NewService(&mockAuthRepo{})

	t.Run("malformed token", func(t *testing.T) {
		_, err := svc.ParseToken("not-a-valid-token")
		if !errors.Is(err, ErrInvalidToken) {
			t.Errorf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("tampered signature", func(t *testing.T) {
		user := &domain.User{ID: 2, Email: "test@local", Role: domain.RoleStaff}
		token, err := createToken(user)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}
		parts := strings.Split(token, ".")
		tampered := parts[0] + ".invalidsignature"

		_, err = svc.ParseToken(tampered)
		if !errors.Is(err, ErrInvalidSignature) {
			t.Errorf("expected ErrInvalidSignature, got %v", err)
		}
	})

	t.Run("expired token", func(t *testing.T) {
		expiredPayload := SessionPayload{
			UserID:    3,
			Email:     "expired@local",
			ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
		}
		// manually encode and sign
		bytes, _ := jsonMarshal(expiredPayload)
		encoded := base64URLEncode(bytes)
		signed := encoded + "." + sign([]byte(encoded))

		_, err := svc.ParseToken(signed)
		if !errors.Is(err, ErrTokenExpired) {
			t.Errorf("expected ErrTokenExpired, got %v", err)
		}
	})
}

func TestAuthHandler(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	mockUser := &domain.User{
		ID:           1,
		Email:        "admin@musashi.local",
		Name:         "Admin",
		PasswordHash: string(hashedPassword),
		Role:         domain.RoleAdmin,
	}

	repo := &mockAuthRepo{
		findByEmailFn: func(email string) (*domain.User, error) {
			if email == "admin@musashi.local" {
				return mockUser, nil
			}
			return nil, errors.New("not found")
		},
	}

	svc := NewService(repo)
	handler := NewHandler(svc)

	app := fiber.New()
	app.Post("/login", handler.Login)
	app.Post("/logout", handler.Logout)
	app.Get("/me", func(c fiber.Ctx) error {
		// Mock locals as set by AuthMiddleware
		c.Locals("user", &SessionPayload{
			UserID: 1,
			Email:  "admin@musashi.local",
			Name:   "Admin",
			Role:   domain.RoleAdmin,
		})
		return handler.Me(c)
	})

	t.Run("HTTP login success", func(t *testing.T) {
		body := strings.NewReader(`{"email":"admin@musashi.local","password":"admin123"}`)
		req := httptest.NewRequest(http.MethodPost, "/login", body)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		cookies := resp.Cookies()
		found := false
		for _, c := range cookies {
			if c.Name == CookieName && c.Value != "" {
				found = true
			}
		}
		if !found {
			t.Errorf("expected session cookie to be set")
		}
	})

	t.Run("HTTP logout clears cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
		for _, c := range resp.Cookies() {
			if c.Name == CookieName && c.MaxAge < 0 {
				return // passed
			}
		}
		t.Errorf("expected cleared cookie with max age < 0")
	})

	t.Run("HTTP me returns current user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})
}

// Helpers for test
func jsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func base64URLEncode(src []byte) string {
	return base64.RawURLEncoding.EncodeToString(src)
}
