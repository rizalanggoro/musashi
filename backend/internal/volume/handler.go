package volume

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

// List godoc
// @id          listVolumes
// @Summary     List project volume rows
// @Tags        volume
// @Produce     json
// @Param       id path string true "Project ID"
// @Param       q query string false "Search query"
// @Success     200 {array} VolumeRowDTO
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /projects/{id}/volumes [get]
func (h *Handler) List(c fiber.Ctx) error {
	projectID, err := utils.ParseID(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid project ID", nil)
	}

	q := c.Query("q")
	data, err := h.svc.List(projectID, q)
	if err != nil {
		return httperr.Internal(c, "Failed to fetch volume take-offs")
	}

	return c.JSON(data)
}

// Create godoc
// @id          createVolume
// @Summary     Create volume entry
// @Tags        volume
// @Accept      json
// @Produce     json
// @Param       id path string true "Project ID"
// @Param       body body CreateVolumeInput true "Volume input"
// @Success     201 {object} VolumeRowDTO
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /projects/{id}/volumes [post]
func (h *Handler) Create(c fiber.Ctx) error {
	projectID, err := utils.ParseID(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid project ID", nil)
	}

	var input CreateVolumeInput
	if err := c.Bind().Body(&input); err != nil {
		return httperr.BadRequest(c, "Invalid request body", nil)
	}

	dto, err := h.svc.Create(projectID, input)
	if err != nil {
		if errors.Is(err, ErrBoqItemNotFound) {
			return httperr.NotFound(c, "BOQ item not found")
		}
		if errors.Is(err, ErrRequiredBoqField) {
			return httperr.BadRequest(c, err.Error(), nil)
		}
		return httperr.Internal(c, "Failed to create volume take-off")
	}

	return c.Status(fiber.StatusCreated).JSON(dto)
}

// Update godoc
// @id          updateVolume
// @Summary     Update volume entry
// @Tags        volume
// @Accept      json
// @Produce     json
// @Param       id path string true "Volume row ID"
// @Param       body body UpdateVolumeInput true "Volume update data"
// @Success     200 {object} VolumeRowDTO
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     404 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /volumes/{id} [put]
func (h *Handler) Update(c fiber.Ctx) error {
	id, err := utils.ParseID(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid volume ID", nil)
	}

	var input UpdateVolumeInput
	if err := c.Bind().Body(&input); err != nil {
		return httperr.BadRequest(c, "Invalid request body", nil)
	}

	dto, err := h.svc.Update(id, input)
	if err != nil {
		if errors.Is(err, ErrVolumeNotFound) {
			return httperr.NotFound(c, "Volume row not found")
		}
		if errors.Is(err, ErrBoqItemNotFound) {
			return httperr.NotFound(c, "Associated BOQ item not found")
		}
		return httperr.Internal(c, "Failed to update volume row")
	}

	return c.JSON(dto)
}

// Delete godoc
// @id          deleteVolume
// @Summary     Delete volume entry
// @Tags        volume
// @Param       id path string true "Volume row ID"
// @Success     204 "No Content"
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     404 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /volumes/{id} [delete]
func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := utils.ParseID(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid volume ID", nil)
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, ErrVolumeNotFound) {
			return httperr.NotFound(c, "Volume row not found")
		}
		return httperr.Internal(c, "Failed to delete volume row")
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
