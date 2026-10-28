package rab

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/platform/httperr"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetProjectRab(c fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid project ID", nil)
	}

	q := c.Query("q")
	dto, err := h.svc.GetProjectRab(projectID, q)
	if err != nil {
		return httperr.Internal(c, "Failed to fetch RAB items")
	}

	return c.JSON(dto)
}

func (h *Handler) Create(c fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid project ID", nil)
	}

	var input CreateRabItemInput
	if err := c.Bind().Body(&input); err != nil {
		return httperr.BadRequest(c, "Invalid request body", nil)
	}

	dto, err := h.svc.Create(projectID, input)
	if err != nil {
		if errors.Is(err, ErrRequiredRabField) {
			return httperr.BadRequest(c, err.Error(), nil)
		}
		return httperr.Internal(c, "Failed to create RAB item")
	}

	return c.Status(fiber.StatusCreated).JSON(dto)
}

func (h *Handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid RAB item ID", nil)
	}

	var input UpdateRabItemInput
	if err := c.Bind().Body(&input); err != nil {
		return httperr.BadRequest(c, "Invalid request body", nil)
	}

	dto, err := h.svc.Update(id, input)
	if err != nil {
		if errors.Is(err, ErrRabItemNotFound) {
			return httperr.NotFound(c, "RAB item not found")
		}
		return httperr.Internal(c, "Failed to update RAB item")
	}

	return c.JSON(dto)
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid RAB item ID", nil)
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, ErrRabItemNotFound) {
			return httperr.NotFound(c, "RAB item not found")
		}
		return httperr.Internal(c, "Failed to delete RAB item")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Backward-compatible helpers
func GetProjectRabHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).GetProjectRab
}

func CreateHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).Create
}

func UpdateHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).Update
}

func DeleteHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).Delete
}
