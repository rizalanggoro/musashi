package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/platform/httperr"
)

const (
	CookieName = "musashi_session"
	CookieMaxAge = 7 * 24 * 3600 // 7 days
)

var sessionSecret = []byte("musashi-super-secret-session-key-change-in-prod")

type SessionPayload struct {
	UserID    uuid.UUID       `json:"sub"`
	Email     string          `json:"email"`
	Name      string          `json:"name"`
	Role      domain.UserRole `json:"role"`
	ExpiresAt int64           `json:"exp"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    string          `json:"id"`
	Email string          `json:"email"`
	Name  string          `json:"name"`
	Role  domain.UserRole `json:"role"`
}

func sign(data []byte) string {
	mac := hmac.New(sha256.New, sessionSecret)
	mac.Write(data)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func createToken(user *domain.User) (string, error) {
	payload := SessionPayload{
		UserID:    user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	encoded := base64.RawURLEncoding.EncodeToString(bytes)
	signature := sign([]byte(encoded))
	return encoded + "." + signature, nil
}

func parseToken(token string) (*SessionPayload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid token format")
	}
	expectedSig := sign([]byte(parts[0]))
	if !hmac.Equal([]byte(expectedSig), []byte(parts[1])) {
		return nil, errors.New("invalid signature")
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	var payload SessionPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	if time.Now().Unix() > payload.ExpiresAt {
		return nil, errors.New("token expired")
	}
	return &payload, nil
}

func LoginHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var input LoginInput
		if err := c.Bind().Body(&input); err != nil {
			return httperr.BadRequest(c, "Invalid request body", nil)
		}
		if input.Email == "" || input.Password == "" {
			return httperr.BadRequest(c, "Email and password are required", nil)
		}

		var user domain.User
		if err := db.Where("lower(email) = ?", strings.ToLower(input.Email)).First(&user).Error; err != nil {
			return httperr.Unauthorized(c, "Invalid email or password")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
			return httperr.Unauthorized(c, "Invalid email or password")
		}

		token, err := createToken(&user)
		if err != nil {
			return httperr.Internal(c, "Failed to create session")
		}

		c.Cookie(&fiber.Cookie{
			Name:     CookieName,
			Value:    token,
			Path:     "/",
			HTTPOnly: true,
			SameSite: "Lax",
			MaxAge:   CookieMaxAge,
		})

		return c.JSON(UserResponse{
			ID:    user.ID.String(),
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		})
	}
}

func LogoutHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:     CookieName,
			Value:    "",
			Path:     "/",
			HTTPOnly: true,
			SameSite: "Lax",
			MaxAge:   -1,
		})
		return c.JSON(fiber.Map{"ok": true})
	}
}

func MeHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		user, ok := c.Locals("user").(*SessionPayload)
		if !ok || user == nil {
			return httperr.Unauthorized(c, "Unauthenticated")
		}
		return c.JSON(UserResponse{
			ID:    user.UserID.String(),
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		})
	}
}

func AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Cookies(CookieName)
		if token == "" {
			return httperr.Unauthorized(c, "Authentication required")
		}
		payload, err := parseToken(token)
		if err != nil {
			return httperr.Unauthorized(c, "Session invalid or expired")
		}
		c.Locals("user", payload)
		return c.Next()
	}
}

func RequireRole(roles ...domain.UserRole) fiber.Handler {
	return func(c fiber.Ctx) error {
		user, ok := c.Locals("user").(*SessionPayload)
		if !ok || user == nil {
			return httperr.Unauthorized(c, "Authentication required")
		}
		for _, r := range roles {
			if user.Role == r {
				return c.Next()
			}
		}
		return httperr.Forbidden(c, "Access denied: insufficient permissions")
	}
}

