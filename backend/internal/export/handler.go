package export

import (
	"errors"
	"fmt"

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

// VolumeExcel godoc
// @id          exportVolumeExcel
// @Summary     Export project volume take-off as Excel
// @Tags        export
// @Produce     application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param       id   path      string  true  "Project ID"
// @Success     200  {file}    binary
// @Failure     400  {object}  httperr.ErrorResponse
// @Failure     404  {object}  httperr.ErrorResponse
// @Failure     500  {object}  httperr.ErrorResponse
// @Router      /projects/{id}/exports/volume.xlsx [get]
func (h *Handler) VolumeExcel(c fiber.Ctx) error {
	projectID, err := utils.ParseID(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid project ID", nil)
	}

	filename, data, err := h.svc.GenerateVolumeExcel(projectID)
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			return httperr.NotFound(c, "Project not found")
		}
		return httperr.Internal(c, "Failed to generate Excel file")
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	return c.Send(data)
}

// RabExcel godoc
// @id          exportRabExcel
// @Summary     Export project RAB as Excel
// @Tags        export
// @Produce     application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param       id   path      string  true  "Project ID"
// @Success     200  {file}    binary
// @Failure     400  {object}  httperr.ErrorResponse
// @Failure     404  {object}  httperr.ErrorResponse
// @Failure     500  {object}  httperr.ErrorResponse
// @Router      /projects/{id}/exports/rab.xlsx [get]
func (h *Handler) RabExcel(c fiber.Ctx) error {
	projectID, err := utils.ParseID(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid project ID", nil)
	}

	filename, data, err := h.svc.GenerateRabExcel(projectID)
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			return httperr.NotFound(c, "Project not found")
		}
		return httperr.Internal(c, "Failed to generate Excel file")
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	return c.Send(data)
}

// RabPdf godoc
// @id          exportRabPdf
// @Summary     Export project RAB as PDF
// @Tags        export
// @Produce     application/pdf
// @Param       id   path      string  true  "Project ID"
// @Success     200  {file}    binary
// @Failure     400  {object}  httperr.ErrorResponse
// @Failure     404  {object}  httperr.ErrorResponse
// @Failure     500  {object}  httperr.ErrorResponse
// @Router      /projects/{id}/exports/rab.pdf [get]
func (h *Handler) RabPdf(c fiber.Ctx) error {
	projectID, err := utils.ParseID(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid project ID", nil)
	}

	filename, data, err := h.svc.GenerateRabPdf(projectID)
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			return httperr.NotFound(c, "Project not found")
		}
		return httperr.Internal(c, "Failed to generate PDF file")
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	return c.Send(data)
}

// Backward-compatible helpers
func VolumeExcelHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).VolumeExcel
}

func RabExcelHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).RabExcel
}

func RabPdfHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).RabPdf
}
