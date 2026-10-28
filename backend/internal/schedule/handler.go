package schedule

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

func (h *Handler) ListByProject(c fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid project ID", nil)
	}

	data, err := h.svc.ListByProject(projectID)
	if err != nil {
		return httperr.Internal(c, "Failed to fetch schedules")
	}

	return c.JSON(data)
}

func (h *Handler) Create(c fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid project ID", nil)
	}

	var input CreateScheduleInput
	if err := c.Bind().Body(&input); err != nil {
		return httperr.BadRequest(c, "Invalid request body", nil)
	}

	dto, err := h.svc.Create(projectID, input)
	if err != nil {
		if errors.Is(err, ErrRequiredFields) || errors.Is(err, ErrInvalidDate) {
			return httperr.BadRequest(c, err.Error(), nil)
		}
		return httperr.Internal(c, "Failed to create schedule")
	}

	return c.Status(fiber.StatusCreated).JSON(dto)
}

func (h *Handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid schedule ID", nil)
	}

	var input UpdateScheduleInput
	if err := c.Bind().Body(&input); err != nil {
		return httperr.BadRequest(c, "Invalid request body", nil)
	}

	dto, err := h.svc.Update(id, input)
	if err != nil {
		if errors.Is(err, ErrScheduleNotFound) {
			return httperr.NotFound(c, "Schedule not found")
		}
		return httperr.Internal(c, "Failed to update schedule")
	}

	return c.JSON(dto)
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid schedule ID", nil)
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, ErrScheduleNotFound) {
			return httperr.NotFound(c, "Schedule not found")
		}
		return httperr.Internal(c, "Failed to delete schedule")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Backward-compatible helpers
func ListByProjectHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).ListByProject
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
