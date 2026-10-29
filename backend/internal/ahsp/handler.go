package ahsp

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
// @id          listAhsps
// @Summary     List AHSP
// @Tags        ahsp
// @Produce     json
// @Param       q query string false "Search query"
// @Success     200 {array} AhspDTO
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /ahsps [get]
func (h *Handler) List(c fiber.Ctx) error {
	q := c.Query("q")
	data, err := h.svc.List(q)
	if err != nil {
		return httperr.Internal(c, "Failed to fetch AHSPs")
	}
	return c.JSON(data)
}

// Get godoc
// @id          getAhsp
// @Summary     Get AHSP by ID
// @Tags        ahsp
// @Produce     json
// @Param       id path string true "AHSP ID"
// @Success     200 {object} AhspDetailDTO
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     404 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /ahsps/{id} [get]
func (h *Handler) Get(c fiber.Ctx) error {
	id, err := utils.ParseID(c.Params("id"))
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

// Create godoc
// @id          createAhsp
// @Summary     Create AHSP
// @Tags        ahsp
// @Accept      json
// @Produce     json
// @Param       body body CreateAhspInput true "AHSP input"
// @Success     201 {object} AhspDetailDTO
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     409 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /ahsps [post]
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

// Update godoc
// @id          updateAhsp
// @Summary     Update AHSP
// @Tags        ahsp
// @Accept      json
// @Produce     json
// @Param       id path string true "AHSP ID"
// @Param       body body CreateAhspInput true "AHSP update data"
// @Success     200 {object} AhspDetailDTO
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     404 {object} httperr.ErrorResponse
// @Failure     409 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /ahsps/{id} [put]
func (h *Handler) Update(c fiber.Ctx) error {
	id, err := utils.ParseID(c.Params("id"))
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

// Delete godoc
// @id          deleteAhsp
// @Summary     Delete AHSP
// @Tags        ahsp
// @Param       id path string true "AHSP ID"
// @Success     204 "No Content"
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     404 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /ahsps/{id} [delete]
func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := utils.ParseID(c.Params("id"))
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
