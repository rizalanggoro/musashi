package model

import (
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type RabItem struct {
	gorm.Model
	ProjectID uint    `gorm:"not null;index"`
	BoqItemID *uint   `gorm:"index"`
	WorkItem  string  `gorm:"not null"`
	Unit      string  `gorm:"not null"`
	Volume    float64 `gorm:"type:numeric(14,4);not null;default:0"`
	UnitPrice int64   `gorm:"not null;default:0"`
	TotalCost int64   `gorm:"not null;default:0"`
}

func (RabItem) TableName() string {
	return "rab_items"
}

func (r *RabItem) ToDomain() *domain.RabItem {
	if r == nil {
		return nil
	}
	return &domain.RabItem{
		ID:        r.ID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		ProjectID: r.ProjectID,
		BoqItemID: r.BoqItemID,
		WorkItem:  r.WorkItem,
		Unit:      r.Unit,
		Volume:    r.Volume,
		UnitPrice: r.UnitPrice,
		TotalCost: r.TotalCost,
	}
}

func FromRabItemDomain(d *domain.RabItem) *RabItem {
	if d == nil {
		return nil
	}
	r := &RabItem{
		ProjectID: d.ProjectID,
		BoqItemID: d.BoqItemID,
		WorkItem:  d.WorkItem,
		Unit:      d.Unit,
		Volume:    d.Volume,
		UnitPrice: d.UnitPrice,
		TotalCost: d.TotalCost,
	}
	r.ID = d.ID
	r.CreatedAt = d.CreatedAt
	r.UpdatedAt = d.UpdatedAt
	return r
}

