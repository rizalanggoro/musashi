package ahsp

import (
	"fmt"
	"strconv"

	"musashi/backend/internal/domain"
	"musashi/backend/utils"
)

type AhspDTO struct {
	ID                string `json:"id"`
	Code              string `json:"code"`
	WorkItem          string `json:"work_item"`
	Unit              string `json:"unit"`
	OverheadProfitPct string `json:"overhead_profit_pct"`
	TotalLabor        string `json:"total_labor"`
	TotalMaterial     string `json:"total_material"`
	TotalEquipment    string `json:"total_equipment"`
	TotalCost         string `json:"total_cost"`
	UnitPriceOfWork   string `json:"unit_price_of_work"`
	Status            string `json:"status"`
}

type AhspItemDTO struct {
	ID          string  `json:"id"`
	UnitPriceID *string `json:"unit_price_id"`
	Category    string  `json:"category"`
	Sort        int     `json:"sort"`
	Description string  `json:"description"`
	Code        string  `json:"code"`
	Unit        string  `json:"unit"`
	Coefficient string  `json:"coefficient"`
	UnitPrice   string  `json:"unit_price"`
	Cost        string  `json:"cost"`
}

type AhspDetailDTO struct {
	AhspDTO
	Items []AhspItemDTO `json:"items"`
}

type CreateAhspItemInput struct {
	UnitPriceID *string `json:"unit_price_id,omitempty"`
	Category    string  `json:"category"`
	Sort        int     `json:"sort"`
	Description string  `json:"description"`
	Code        string  `json:"code"`
	Unit        string  `json:"unit"`
	Coefficient string  `json:"coefficient"`
	UnitPrice   string  `json:"unit_price"`
}

type CreateAhspInput struct {
	Code              string                `json:"code"`
	WorkItem          string                `json:"work_item"`
	Unit              string                `json:"unit"`
	OverheadProfitPct string                `json:"overhead_profit_pct"`
	Status            string                `json:"status"`
	Items             []CreateAhspItemInput `json:"items"`
}

func toAhspDTO(a *domain.Ahsp) AhspDTO {
	return AhspDTO{
		ID:                utils.FormatID(a.ID),
		Code:              a.Code,
		WorkItem:          a.WorkItem,
		Unit:              a.Unit,
		OverheadProfitPct: fmt.Sprintf("%.2f", a.OverheadProfitPct),
		TotalLabor:        strconv.FormatInt(a.TotalLabor, 10),
		TotalMaterial:     strconv.FormatInt(a.TotalMaterial, 10),
		TotalEquipment:    strconv.FormatInt(a.TotalEquipment, 10),
		TotalCost:         strconv.FormatInt(a.TotalCost, 10),
		UnitPriceOfWork:   strconv.FormatInt(a.UnitPriceOfWork, 10),
		Status:            string(a.Status),
	}
}

func toItemDTO(item *domain.AhspItem) AhspItemDTO {
	dto := AhspItemDTO{
		ID:          utils.FormatID(item.ID),
		Category:    string(item.Category),
		Sort:        item.Sort,
		Description: item.Description,
		Code:        item.Code,
		Unit:        item.Unit,
		Coefficient: fmt.Sprintf("%.4f", item.Coefficient),
		UnitPrice:   strconv.FormatInt(item.UnitPrice, 10),
		Cost:        strconv.FormatInt(item.Cost, 10),
	}
	if item.UnitPriceID != nil {
		s := utils.FormatID(*item.UnitPriceID)
		dto.UnitPriceID = &s
	}
	return dto
}
