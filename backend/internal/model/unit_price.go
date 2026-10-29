package model

import (
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type UnitPrice struct {
	gorm.Model
	Code        string                   `gorm:"uniqueIndex;not null"`
	Category    domain.UnitPriceCategory `gorm:"type:text;not null;index"`
	WorkItem    string                   `gorm:"not null"`
	Unit        string                   `gorm:"not null"`
	UnitPrice   int64                    `gorm:"not null;default:0"`
	Information string                   `gorm:"not null;default:''"`
}

func (UnitPrice) TableName() string {
	return "unit_prices"
}

func (u *UnitPrice) ToDomain() *domain.UnitPrice {
	if u == nil {
		return nil
	}
	return &domain.UnitPrice{
		ID:          u.ID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Code:        u.Code,
		Category:    u.Category,
		WorkItem:    u.WorkItem,
		Unit:        u.Unit,
		UnitPrice:   u.UnitPrice,
		Information: u.Information,
	}
}

func FromUnitPriceDomain(d *domain.UnitPrice) *UnitPrice {
	if d == nil {
		return nil
	}
	u := &UnitPrice{
		Code:        d.Code,
		Category:    d.Category,
		WorkItem:    d.WorkItem,
		Unit:        d.Unit,
		UnitPrice:   d.UnitPrice,
		Information: d.Information,
	}
	u.ID = d.ID
	u.CreatedAt = d.CreatedAt
	u.UpdatedAt = d.UpdatedAt
	return u
}

