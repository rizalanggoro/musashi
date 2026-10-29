package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"

	"musashi/backend/internal/auth"
	"musashi/backend/internal/domain"
)

type mockAuthService struct {
	parseTokenFn func(token string) (*auth.SessionPayload, error)
}

func (m *mockAuthService) Login(input auth.LoginInput) (*domain.User, string, error) {
	return nil, "", nil
}

func (m *mockAuthService) ParseToken(token string) (*auth.SessionPayload, error) {
	if m.parseTokenFn != nil {
		return m.parseTokenFn(token)
	}
	return nil, errors.New("invalid token")
}

func TestAuthMiddleware(t *testing.T) {
	svc := &mockAuthService{
		parseTokenFn: func(token string) (*auth.SessionPayload, error) {
			if token == "valid-token" {
				return &auth.SessionPayload{
					UserID: 1,
					Email:  "test@example.com",
					Name:   "Tester",
					Role:   domain.RoleAdmin,
				}, nil
			}
			return nil, errors.New("bad token")
		},
	}

	app := fiber.New()
	app.Get("/protected", AuthMiddleware(svc), func(c fiber.Ctx) error {
		user := GetCurrentUser(c)
		return c.SendString("Hello " + user.Name)
	})

	t.Run("missing cookie returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("invalid token returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: "invalid-token"})
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("valid token succeeds", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: "valid-token"})
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})
}

func TestRequireRole(t *testing.T) {
	app := fiber.New()
	app.Get("/admin-only", func(c fiber.Ctx) error {
		role := c.Query("role")
		if role != "" {
			c.Locals("user", &auth.SessionPayload{
				UserID: 1,
				Role:   domain.UserRole(role),
			})
		}
		return c.Next()
	}, RequireRole(domain.RoleAdmin), func(c fiber.Ctx) error {
		return c.SendString("admin ok")
	})

	t.Run("unauthenticated user returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}
	})

	t.Run("staff role returns 403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin-only?role=staff", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("expected 403, got %d", resp.StatusCode)
		}
	})

	t.Run("admin role returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin-only?role=admin", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
	})
}

