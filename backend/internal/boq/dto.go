package boq

import (
	"fmt"

	"musashi/backend/internal/domain"
	"musashi/backend/utils"
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

func toItemDTO(item *domain.BoqItem) BoqItemDTO {
	d := BoqItemDTO{
		ID:            utils.FormatID(item.ID),
		ProjectID:     utils.FormatID(item.ProjectID),
		WorkItem:      item.WorkItem,
		Unit:          item.Unit,
		PlannedVolume: fmt.Sprintf("%.4f", item.PlannedVolume),
	}
	if item.AhspID != nil {
		s := utils.FormatID(*item.AhspID)
		d.AhspID = &s
	}
	return d
}
