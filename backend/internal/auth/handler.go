package auth

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/platform/httperr"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Login(c fiber.Ctx) error {
	var input LoginInput
	if err := c.Bind().Body(&input); err != nil {
		return httperr.BadRequest(c, "Invalid request body", nil)
	}
	if input.Email == "" || input.Password == "" {
		return httperr.BadRequest(c, "Email and password are required", nil)
	}

	user, token, err := h.svc.Login(input)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return httperr.Unauthorized(c, "Invalid email or password")
		}
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

func (h *Handler) Logout(c fiber.Ctx) error {
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

func (h *Handler) Me(c fiber.Ctx) error {
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

func (h *Handler) AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Cookies(CookieName)
		if token == "" {
			return httperr.Unauthorized(c, "Authentication required")
		}
		payload, err := h.svc.ParseToken(token)
		if err != nil {
			return httperr.Unauthorized(c, "Session invalid or expired")
		}
		c.Locals("user", payload)
		return c.Next()
	}
}

func (h *Handler) RequireRole(roles ...domain.UserRole) fiber.Handler {
	return RequireRole(roles...)
}

// Backward-compatible standalone helpers

func LoginHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).Login
}

func LogoutHandler() fiber.Handler {
	return (&Handler{}).Logout
}

func MeHandler() fiber.Handler {
	return (&Handler{}).Me
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
