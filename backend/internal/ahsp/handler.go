package ahsp

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
	data, err := h.svc.List(q)
	if err != nil {
		return httperr.Internal(c, "Failed to fetch AHSPs")
	}
	return c.JSON(data)
}

func (h *Handler) Get(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid AHSP ID", nil)
	}

	dto, err := h.svc.Get(id)
	if err != nil {
		if errors.Is(err, ErrAhspNotFound) {
			return httperr.NotFound(c, "AHSP not found")
		}
		return httperr.Internal(c, "Failed to fetch AHSP")
	}

	return c.JSON(dto)
}

func (h *Handler) Create(c fiber.Ctx) error {
	var input CreateAhspInput
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
		return httperr.Internal(c, "Failed to create AHSP")
	}

	return c.Status(fiber.StatusCreated).JSON(dto)
}

func (h *Handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid AHSP ID", nil)
	}

	var input CreateAhspInput
	if err := c.Bind().Body(&input); err != nil {
		return httperr.BadRequest(c, "Invalid request body", nil)
	}

	dto, err := h.svc.Update(id, input)
	if err != nil {
		if errors.Is(err, ErrAhspNotFound) {
			return httperr.NotFound(c, "AHSP not found")
		}
		if errors.Is(err, ErrRequiredFields) {
			return httperr.BadRequest(c, err.Error(), nil)
		}
		if errors.Is(err, ErrCodeAlreadyExists) {
			return httperr.Conflict(c, err.Error())
		}
		return httperr.Internal(c, "Failed to update AHSP")
	}

	return c.JSON(dto)
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid AHSP ID", nil)
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, ErrAhspNotFound) {
			return httperr.NotFound(c, "AHSP not found")
		}
		return httperr.Internal(c, "Failed to delete AHSP")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Backward-compatible helpers
func ListHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).List
}

func GetHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).Get
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
