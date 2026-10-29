package unitprice

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
// @id          listUnitPrices
// @Summary     List unit prices
// @Tags        unit-prices
// @Produce     json
// @Param       q query string false "Search query"
// @Param       category query string false "Filter by category (LABOR, MATERIALS, EQUIPMENT, OTHER)"
// @Success     200 {array} UnitPriceDTO
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /unit-prices [get]
func (h *Handler) List(c fiber.Ctx) error {
	q := c.Query("q")
	category := c.Query("category")

	data, err := h.svc.List(q, category)
	if err != nil {
		return httperr.Internal(c, "Failed to fetch unit prices")
	}

	return c.JSON(data)
}

// Create godoc
// @id          createUnitPrice
// @Summary     Create unit price
// @Tags        unit-prices
// @Accept      json
// @Produce     json
// @Param       body body CreateUnitPriceInput true "Unit price input"
// @Success     201 {object} UnitPriceDTO
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     409 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /unit-prices [post]
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

// Update godoc
// @id          updateUnitPrice
// @Summary     Update unit price
// @Tags        unit-prices
// @Accept      json
// @Produce     json
// @Param       id path string true "Unit Price ID"
// @Param       body body UpdateUnitPriceInput true "Unit price update data"
// @Success     200 {object} UnitPriceDTO
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     404 {object} httperr.ErrorResponse
// @Failure     409 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /unit-prices/{id} [put]
func (h *Handler) Update(c fiber.Ctx) error {
	id, err := utils.ParseID(c.Params("id"))
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

// Delete godoc
// @id          deleteUnitPrice
// @Summary     Delete unit price
// @Tags        unit-prices
// @Param       id path string true "Unit Price ID"
// @Success     204 "No Content"
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     404 {object} httperr.ErrorResponse
// @Failure     409 {object} httperr.ErrorResponse
// @Router      /unit-prices/{id} [delete]
func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := utils.ParseID(c.Params("id"))
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
