package schedule

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"musashi/backend/internal/platform/httperr"
	"musashi/backend/utils"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// ListByProject godoc
// @id          listSchedules
// @Summary     List project schedules
// @Tags        schedules
// @Produce     json
// @Param       id path string true "Project ID"
// @Success     200 {array} ScheduleDTO
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /projects/{id}/schedules [get]
func (h *Handler) ListByProject(c fiber.Ctx) error {
	projectID, err := utils.ParseID(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid project ID", nil)
	}

	data, err := h.svc.ListByProject(projectID)
	if err != nil {
		return httperr.Internal(c, "Failed to fetch schedules")
	}

	return c.JSON(data)
}

// Create godoc
// @id          createSchedule
// @Summary     Create schedule item
// @Tags        schedules
// @Accept      json
// @Produce     json
// @Param       id path string true "Project ID"
// @Param       body body CreateScheduleInput true "Schedule input"
// @Success     201 {object} ScheduleDTO
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /projects/{id}/schedules [post]
func (h *Handler) Create(c fiber.Ctx) error {
	projectID, err := utils.ParseID(c.Params("id"))
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

// Update godoc
// @id          updateSchedule
// @Summary     Update schedule item
// @Tags        schedules
// @Accept      json
// @Produce     json
// @Param       id path string true "Schedule ID"
// @Param       body body UpdateScheduleInput true "Schedule update data"
// @Success     200 {object} ScheduleDTO
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     404 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /schedules/{id} [put]
func (h *Handler) Update(c fiber.Ctx) error {
	id, err := utils.ParseID(c.Params("id"))
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

// Delete godoc
// @id          deleteSchedule
// @Summary     Delete schedule item
// @Tags        schedules
// @Param       id path string true "Schedule ID"
// @Success     204 "No Content"
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     404 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /schedules/{id} [delete]
func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := utils.ParseID(c.Params("id"))
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
