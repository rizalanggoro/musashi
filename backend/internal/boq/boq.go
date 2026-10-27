package boq

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/platform/httperr"
)

type BoqProjectSummaryDTO struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Section      string `json:"section"`
	Plant        string `json:"plant"`
	VolumeStatus string `json:"volume_status"`
	RabStatus    string `json:"rab_status"`
}

type BoqItemDTO struct {
	ID            string  `json:"id"`
	ProjectID     string  `json:"project_id"`
	AhspID        *string `json:"ahsp_id"`
	WorkItem      string  `json:"work_item"`
	Unit          string  `json:"unit"`
	PlannedVolume string  `json:"planned_volume"`
}

func ListSummaryHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		q := strings.TrimSpace(c.Query("q"))

		query := db.Model(&domain.Project{}).Where("deleted_at IS NULL")
		if q != "" {
			like := "%" + strings.ToLower(q) + "%"
			query = query.Where("lower(name) LIKE ? OR lower(section) LIKE ? OR lower(plant) LIKE ?", like, like, like)
		}

		var projects []domain.Project
		if err := query.Order("created_at DESC").Find(&projects).Error; err != nil {
			return httperr.Internal(c, "Failed to fetch projects")
		}

		results := make([]BoqProjectSummaryDTO, len(projects))
		for i, p := range projects {
			var boqCount int64
			db.Model(&domain.BoqItem{}).Where("project_id = ?", p.ID).Count(&boqCount)

			var volCount int64
			db.Table("volumes").
				Joins("JOIN boq_items ON boq_items.id = volumes.boq_item_id").
				Where("boq_items.project_id = ?", p.ID).
				Count(&volCount)

			var rabCount int64
			db.Model(&domain.RabItem{}).Where("project_id = ?", p.ID).Count(&rabCount)

			volumeStatus := "empty"
			if boqCount > 0 {
				if volCount == 0 {
					volumeStatus = "empty"
				} else if volCount < boqCount {
					volumeStatus = "partial"
				} else {
					volumeStatus = "complete"
				}
			}

			rabStatus := "empty"
			if rabCount > 0 {
				if boqCount > 0 && rabCount < boqCount {
					rabStatus = "partial"
				} else {
					rabStatus = "complete"
				}
			}

			results[i] = BoqProjectSummaryDTO{
				ID:           p.ID.String(),
				Name:         p.Name,
				Section:      p.Section,
				Plant:        p.Plant,
				VolumeStatus: volumeStatus,
				RabStatus:    rabStatus,
			}
		}

		return c.JSON(results)
	}
}

func GetProjectBoqHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		projectID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid project ID", nil)
		}

		var items []domain.BoqItem
		if err := db.Where("project_id = ?", projectID).Order("created_at ASC").Find(&items).Error; err != nil {
			return httperr.Internal(c, "Failed to fetch BOQ items")
		}

		dto := make([]BoqItemDTO, len(items))
		for i, item := range items {
			d := BoqItemDTO{
				ID:            item.ID.String(),
				ProjectID:     item.ProjectID.String(),
				WorkItem:      item.WorkItem,
				Unit:          item.Unit,
				PlannedVolume: fmt.Sprintf("%.4f", item.PlannedVolume),
			}
			if item.AhspID != nil {
				s := item.AhspID.String()
				d.AhspID = &s
			}
			dto[i] = d
		}

		return c.JSON(dto)
	}
}

