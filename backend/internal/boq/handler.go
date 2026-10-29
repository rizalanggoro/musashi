package boq

import (
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

func (h *Handler) ListSummary(c fiber.Ctx) error {
	q := c.Query("q")
	data, err := h.svc.ListSummary(q)
	if err != nil {
		return httperr.Internal(c, "Failed to fetch projects")
	}
	return c.JSON(data)
}

func (h *Handler) GetProjectBoq(c fiber.Ctx) error {
	projectID, err := utils.ParseID(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid project ID", nil)
	}

	dto, err := h.svc.GetProjectBoq(projectID)
	if err != nil {
		return httperr.Internal(c, "Failed to fetch BOQ items")
	}

	return c.JSON(dto)
}

// Backward-compatible helpers
func ListSummaryHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).ListSummary
}

func GetProjectBoqHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).GetProjectBoq
}
