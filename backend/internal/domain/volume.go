package domain

import "time"

type Volume struct {
	ID            uint       `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	BoqItemID     uint       `json:"boq_item_id"`
	Date          *time.Time `json:"date"`
	Length        *float64   `json:"length"`
	Width         *float64   `json:"width"`
	Height        *float64   `json:"height"`
	Quantity      float64    `json:"quantity"`
	UnitWeight    float64    `json:"unit_weight"`
	TotalQuantity float64    `json:"total_quantity"`
	Description   string     `json:"description"`
}
