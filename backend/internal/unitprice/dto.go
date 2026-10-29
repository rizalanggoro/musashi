package unitprice

import (
	"strconv"

	"musashi/backend/internal/domain"
	"musashi/backend/utils"
)

type UnitPriceDTO struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Category    string `json:"category"`
	WorkItem    string `json:"work_item"`
	Unit        string `json:"unit"`
	UnitPrice   string `json:"unit_price"`
	Information string `json:"information"`
}

type CreateUnitPriceInput struct {
	Code        string `json:"code"`
	Category    string `json:"category"`
	WorkItem    string `json:"work_item"`
	Unit        string `json:"unit"`
	UnitPrice   string `json:"unit_price"`
	Information string `json:"information"`
}

type UpdateUnitPriceInput struct {
	Code        *string `json:"code"`
	Category    *string `json:"category"`
	WorkItem    *string `json:"work_item"`
	Unit        *string `json:"unit"`
	UnitPrice   *string `json:"unit_price"`
	Information *string `json:"information"`
}

func toDTO(u *domain.UnitPrice) UnitPriceDTO {
	return UnitPriceDTO{
		ID:          utils.FormatID(u.ID),
		Code:        u.Code,
		Category:    string(u.Category),
		WorkItem:    u.WorkItem,
		Unit:        u.Unit,
		UnitPrice:   strconv.FormatInt(u.UnitPrice, 10),
		Information: u.Information,
	}
}
