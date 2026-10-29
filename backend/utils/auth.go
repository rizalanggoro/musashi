package utils

import (
	"github.com/gofiber/fiber/v3"

	"musashi/backend/internal/auth"
	"musashi/backend/internal/domain"
	"musashi/backend/internal/platform/httperr"
)

// GetCurrentUser retrieves the authenticated session payload from fiber context locals.
func GetCurrentUser(c fiber.Ctx) *auth.SessionPayload {
	user, ok := c.Locals("user").(*auth.SessionPayload)
	if !ok || user == nil {
		return nil
	}
	return user
}

// AuthMiddleware creates a Fiber middleware that validates session cookies via auth.Service.
func AuthMiddleware(svc auth.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Cookies(auth.CookieName)
		if token == "" {
			return httperr.Unauthorized(c, "Authentication required")
		}
		payload, err := svc.ParseToken(token)
		if err != nil {
			return httperr.Unauthorized(c, "Session invalid or expired")
		}
		c.Locals("user", payload)
		return c.Next()
	}
}

// RequireRole checks if the authenticated user has one of the allowed roles.
func RequireRole(roles ...domain.UserRole) fiber.Handler {
	return func(c fiber.Ctx) error {
		user := GetCurrentUser(c)
		if user == nil {
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

