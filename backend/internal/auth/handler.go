package auth

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"musashi/backend/internal/platform/httperr"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// Login godoc
// @id          login
// @Summary     Login
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body LoginInput true "Login credentials"
// @Success     200 {object} User
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     401 {object} httperr.ErrorResponse
// @Router      /auth/login [post]
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
		ID:    strconv.FormatUint(uint64(user.ID), 10),
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
		Token: token,
	})
}

// Logout godoc
// @id          logout
// @Summary     Logout
// @Tags        auth
// @Produce     json
// @Success     200 {object} map[string]bool
// @Router      /auth/logout [post]
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

// Me godoc
// @id          getMe
// @Summary     Get current authenticated user
// @Tags        auth
// @Produce     json
// @Success     200 {object} User
// @Failure     401 {object} httperr.ErrorResponse
// @Router      /auth/me [get]
func (h *Handler) Me(c fiber.Ctx) error {
	user, ok := c.Locals("user").(*SessionPayload)
	if !ok || user == nil {
		return httperr.Unauthorized(c, "Unauthenticated")
	}
	return c.JSON(UserResponse{
		ID:    strconv.FormatUint(uint64(user.UserID), 10),
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	})
}
