package httperr

import "github.com/gofiber/fiber/v3"

type ErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

func Respond(c fiber.Ctx, status int, code, message string, fields map[string]string) error {
	return c.Status(status).JSON(ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Fields:  fields,
		},
	})
}

func BadRequest(c fiber.Ctx, message string, fields map[string]string) error {
	return Respond(c, fiber.StatusBadRequest, "BAD_REQUEST", message, fields)
}

func Unauthorized(c fiber.Ctx, message string) error {
	return Respond(c, fiber.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

func Forbidden(c fiber.Ctx, message string) error {
	return Respond(c, fiber.StatusForbidden, "FORBIDDEN", message, nil)
}

func NotFound(c fiber.Ctx, message string) error {
	return Respond(c, fiber.StatusNotFound, "NOT_FOUND", message, nil)
}

func Conflict(c fiber.Ctx, message string) error {
	return Respond(c, fiber.StatusConflict, "CONFLICT", message, nil)
}

func Internal(c fiber.Ctx, message string) error {
	return Respond(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", message, nil)
}

