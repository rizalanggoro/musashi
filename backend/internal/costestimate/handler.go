package costestimate

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

// GetSummary godoc
// @id          getCostEstimate
// @Summary     Get project cost estimate summary
// @Tags        cost-estimate
// @Produce     json
// @Param       id   path      string  true  "Project ID"
// @Success     200  {object}  CostEstimateDTO
// @Failure     400  {object}  httperr.ErrorResponse
// @Failure     404  {object}  httperr.ErrorResponse
// @Failure     500  {object}  httperr.ErrorResponse
// @Router      /projects/{id}/cost-estimate [get]
func (h *Handler) GetSummary(c fiber.Ctx) error {
	projectID, err := utils.ParseID(c.Params("id"))
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
