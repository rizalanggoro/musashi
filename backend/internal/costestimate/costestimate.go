package costestimate

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/platform/httperr"
)

type CostEstimateDTO struct {
	ProjectID string `json:"project_id"`
	Budget    string `json:"budget"`
	Rab       string `json:"rab"`
	Spb       string `json:"spb"`
	Variance  string `json:"variance"`
	Sisa      string `json:"sisa"`
}

func GetSummaryHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		projectID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid project ID", nil)
		}

		var project domain.Project
		if err := db.Where("id = ? AND deleted_at IS NULL", projectID).First(&project).Error; err != nil {
			return httperr.NotFound(c, "Project not found")
		}

		var totalRab int64
		db.Model(&domain.RabItem{}).Where("project_id = ?", projectID).Select("COALESCE(SUM(total_cost), 0)").Scan(&totalRab)

		var spb int64
		if project.ProgressSPB == domain.ProgressDone {
			spb = totalRab
			if spb == 0 {
				spb = project.Budget
			}
		}

		variance := project.Budget - totalRab
		sisa := project.Budget - spb

		return c.JSON(CostEstimateDTO{
			ProjectID: projectID.String(),
			Budget:    strconv.FormatInt(project.Budget, 10),
			Rab:       strconv.FormatInt(totalRab, 10),
			Spb:       strconv.FormatInt(spb, 10),
			Variance:  strconv.FormatInt(variance, 10),
			Sisa:      strconv.FormatInt(sisa, 10),
		})
	}
}

