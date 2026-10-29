package rab

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

// GetProjectRab godoc
// @id          getProjectRab
// @Summary     Get project RAB
// @Tags        rab
// @Produce     json
// @Param       id path string true "Project ID"
// @Param       q query string false "Search query"
// @Success     200 {object} RabResponseDTO
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /projects/{id}/rab [get]
func (h *Handler) GetProjectRab(c fiber.Ctx) error {
	projectID, err := utils.ParseID(c.Params("id"))
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

// Create godoc
// @id          createRabItem
// @Summary     Create RAB item
// @Tags        rab
// @Accept      json
// @Produce     json
// @Param       id path string true "Project ID"
// @Param       body body CreateRabItemInput true "RAB item input"
// @Success     201 {object} RabItemDTO
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /projects/{id}/rab [post]
func (h *Handler) Create(c fiber.Ctx) error {
	projectID, err := utils.ParseID(c.Params("id"))
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

// Update godoc
// @id          updateRabItem
// @Summary     Update RAB item
// @Tags        rab
// @Accept      json
// @Produce     json
// @Param       id path string true "RAB item ID"
// @Param       body body UpdateRabItemInput true "RAB item update data"
// @Success     200 {object} RabItemDTO
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     404 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /rab/{id} [put]
func (h *Handler) Update(c fiber.Ctx) error {
	id, err := utils.ParseID(c.Params("id"))
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

// Delete godoc
// @id          deleteRabItem
// @Summary     Delete RAB item
// @Tags        rab
// @Param       id path string true "RAB item ID"
// @Success     204 "No Content"
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     404 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /rab/{id} [delete]
func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := utils.ParseID(c.Params("id"))
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
