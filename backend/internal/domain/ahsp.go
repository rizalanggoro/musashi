package domain

import "time"

type AhspStatus string

const (
	AhspStatusDraft AhspStatus = "draft"
	AhspStatusFinal AhspStatus = "final"
)

type Ahsp struct {
	ID                uint       `json:"id"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	Code              string     `json:"code"`
	WorkItem          string     `json:"work_item"`
	Unit              string     `json:"unit"`
	OverheadProfitPct float64    `json:"overhead_profit_pct"`
	TotalLabor        int64      `json:"total_labor"`
	TotalMaterial     int64      `json:"total_material"`
	TotalEquipment    int64      `json:"total_equipment"`
	TotalCost         int64      `json:"total_cost"`
	UnitPriceOfWork   int64      `json:"unit_price_of_work"`
	Status            AhspStatus `json:"status"`

	Items []AhspItem `json:"items,omitempty"`
}

type AhspItem struct {
	ID          uint              `json:"id"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	AhspID      uint              `json:"ahsp_id"`
	Category    UnitPriceCategory `json:"category"`
	Sort        int               `json:"sort"`
	Description string            `json:"description"`
	Code        string            `json:"code"`
	Unit        string            `json:"unit"`
	Coefficient float64           `json:"coefficient"`
	UnitPrice   int64             `json:"unit_price"`
	Cost        int64             `json:"cost"`
	UnitPriceID *uint             `json:"unit_price_id,omitempty"`
}
