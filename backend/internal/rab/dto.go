package rab

import (
	"fmt"
	"strconv"

	"musashi/backend/internal/domain"
	"musashi/backend/utils"
)

type RabItemDTO struct {
	ID        string  `json:"id"`
	ProjectID string  `json:"project_id"`
	BoqItemID *string `json:"boq_item_id"`
	WorkItem  string  `json:"work_item"`
	Unit      string  `json:"unit"`
	Volume    string  `json:"volume"`
	UnitPrice string  `json:"unit_price"`
	TotalCost string  `json:"total_cost"`
}

type RabResponseDTO struct {
	ProjectID string       `json:"project_id"`
	Items     []RabItemDTO `json:"items"`
	TotalRab  string       `json:"total_rab"`
}

type CreateRabItemInput struct {
	BoqItemID *string `json:"boq_item_id,omitempty"`
	WorkItem  string  `json:"work_item"`
	Unit      string  `json:"unit"`
	Volume    string  `json:"volume"`
	UnitPrice string  `json:"unit_price"`
}

type UpdateRabItemInput struct {
	WorkItem  *string `json:"work_item,omitempty"`
	Unit      *string `json:"unit,omitempty"`
	Volume    *string `json:"volume,omitempty"`
	UnitPrice *string `json:"unit_price,omitempty"`
}

func toDTO(item *domain.RabItem) RabItemDTO {
	dto := RabItemDTO{
		ID:        utils.FormatID(item.ID),
		ProjectID: utils.FormatID(item.ProjectID),
		WorkItem:  item.WorkItem,
		Unit:      item.Unit,
		Volume:    fmt.Sprintf("%.4f", item.Volume),
		UnitPrice: strconv.FormatInt(item.UnitPrice, 10),
		TotalCost: strconv.FormatInt(item.TotalCost, 10),
	}
	if item.BoqItemID != nil {
		s := utils.FormatID(*item.BoqItemID)
		dto.BoqItemID = &s
	}
	return dto
}
