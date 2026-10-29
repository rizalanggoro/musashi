package model

import (
	"time"

	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type Volume struct {
	gorm.Model
	BoqItemID     uint       `gorm:"not null;index"`
	Date          *time.Time `gorm:"type:date"`
	Length        *float64   `gorm:"type:numeric(14,4)"`
	Width         *float64   `gorm:"type:numeric(14,4)"`
	Height        *float64   `gorm:"type:numeric(14,4)"`
	Quantity      float64    `gorm:"type:numeric(14,4);not null;default:0"`
	UnitWeight    float64    `gorm:"type:numeric(14,4);not null;default:1"`
	TotalQuantity float64    `gorm:"type:numeric(14,4);not null;default:0"`
	Description   string     `gorm:"not null;default:''"`
}

func (Volume) TableName() string {
	return "volumes"
}

func (v *Volume) ToDomain() *domain.Volume {
	if v == nil {
		return nil
	}
	return &domain.Volume{
		ID:            v.ID,
		CreatedAt:     v.CreatedAt,
		UpdatedAt:     v.UpdatedAt,
		BoqItemID:     v.BoqItemID,
		Date:          v.Date,
		Length:        v.Length,
		Width:         v.Width,
		Height:        v.Height,
		Quantity:      v.Quantity,
		UnitWeight:    v.UnitWeight,
		TotalQuantity: v.TotalQuantity,
		Description:   v.Description,
	}
}

func FromVolumeDomain(d *domain.Volume) *Volume {
	if d == nil {
		return nil
	}
	v := &Volume{
		BoqItemID:     d.BoqItemID,
		Date:          d.Date,
		Length:        d.Length,
		Width:         d.Width,
		Height:        d.Height,
		Quantity:      d.Quantity,
		UnitWeight:    d.UnitWeight,
		TotalQuantity: d.TotalQuantity,
		Description:   d.Description,
	}
	v.ID = d.ID
	v.CreatedAt = d.CreatedAt
	v.UpdatedAt = d.UpdatedAt
	return v
}

