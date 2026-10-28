package costestimate

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

func (h *Handler) GetSummary(c fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid project ID", nil)
	}

	dto, err := h.svc.GetSummary(projectID)
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			return httperr.NotFound(c, "Project not found")
		}
		return httperr.Internal(c, "Failed to calculate cost estimate")
	}

	return c.JSON(dto)
}

// Backward-compatible helper
func GetSummaryHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).GetSummary
}
