package export

import (
	"errors"
	"fmt"

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

func (h *Handler) VolumeExcel(c fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
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

func (h *Handler) RabExcel(c fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
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

func (h *Handler) RabPdf(c fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
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
