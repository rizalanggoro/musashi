package domain

import "time"

type RabItem struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ProjectID uint      `json:"project_id"`
	BoqItemID *uint     `json:"boq_item_id,omitempty"`
	WorkItem  string    `json:"work_item"`
	Unit      string    `json:"unit"`
	Volume    float64   `json:"volume"`
	UnitPrice int64     `json:"unit_price"`
	TotalCost int64     `json:"total_cost"`
}
