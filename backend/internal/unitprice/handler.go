package unitprice

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

func (h *Handler) List(c fiber.Ctx) error {
	q := c.Query("q")
	category := c.Query("category")

	data, err := h.svc.List(q, category)
	if err != nil {
		return httperr.Internal(c, "Failed to fetch unit prices")
	}

	return c.JSON(data)
}

func (h *Handler) Create(c fiber.Ctx) error {
	var input CreateUnitPriceInput
	if err := c.Bind().Body(&input); err != nil {
		return httperr.BadRequest(c, "Invalid request body", nil)
	}

	dto, err := h.svc.Create(input)
	if err != nil {
		if errors.Is(err, ErrRequiredFields) {
			return httperr.BadRequest(c, err.Error(), nil)
		}
		if errors.Is(err, ErrCodeAlreadyExists) {
			return httperr.Conflict(c, err.Error())
		}
		return httperr.Internal(c, "Failed to create unit price")
	}

	return c.Status(fiber.StatusCreated).JSON(dto)
}

func (h *Handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid unit price ID", nil)
	}

	var input UpdateUnitPriceInput
	if err := c.Bind().Body(&input); err != nil {
		return httperr.BadRequest(c, "Invalid request body", nil)
	}

	dto, err := h.svc.Update(id, input)
	if err != nil {
		if errors.Is(err, ErrUnitPriceNotFound) {
			return httperr.NotFound(c, "Unit price not found")
		}
		if errors.Is(err, ErrCodeAlreadyExists) {
			return httperr.Conflict(c, err.Error())
		}
		return httperr.Internal(c, "Failed to update unit price")
	}

	return c.JSON(dto)
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid unit price ID", nil)
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, ErrReferencedByAhspItem) {
			return httperr.Conflict(c, err.Error())
		}
		if errors.Is(err, ErrUnitPriceNotFound) {
			return httperr.NotFound(c, "Unit price not found")
		}
		return httperr.Internal(c, "Failed to delete unit price")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Backward-compatible helpers
func ListHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).List
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
