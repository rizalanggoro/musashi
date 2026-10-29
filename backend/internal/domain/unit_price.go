package domain

import "time"

type UnitPriceCategory string

const (
	CategoryLabor     UnitPriceCategory = "LABOR"
	CategoryMaterials UnitPriceCategory = "MATERIALS"
	CategoryEquipment UnitPriceCategory = "EQUIPMENT"
	CategoryOther     UnitPriceCategory = "OTHER"
)

type UnitPrice struct {
	ID          uint              `json:"id"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Code        string            `json:"code"`
	Category    UnitPriceCategory `json:"category"`
	WorkItem    string            `json:"work_item"`
	Unit        string            `json:"unit"`
	UnitPrice   int64             `json:"unit_price"`
	Information string            `json:"information"`
}
