package model

import (
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type Ahsp struct {
	gorm.Model
	Code              string            `gorm:"uniqueIndex;not null"`
	WorkItem          string            `gorm:"not null"`
	Unit              string            `gorm:"not null"`
	OverheadProfitPct float64           `gorm:"type:numeric(6,2);not null;default:10.00"`
	TotalLabor        int64             `gorm:"not null;default:0"`
	TotalMaterial     int64             `gorm:"not null;default:0"`
	TotalEquipment    int64             `gorm:"not null;default:0"`
	TotalCost         int64             `gorm:"not null;default:0"`
	UnitPriceOfWork   int64             `gorm:"not null;default:0"`
	Status            domain.AhspStatus `gorm:"type:text;default:'draft'"`

	Items []AhspItem `gorm:"foreignKey:AhspID"`
}

func (Ahsp) TableName() string {
	return "ahsps"
}

func (a *Ahsp) ToDomain() *domain.Ahsp {
	if a == nil {
		return nil
	}
	d := &domain.Ahsp{
		ID:                a.ID,
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
		Code:              a.Code,
		WorkItem:          a.WorkItem,
		Unit:              a.Unit,
		OverheadProfitPct: a.OverheadProfitPct,
		TotalLabor:        a.TotalLabor,
		TotalMaterial:     a.TotalMaterial,
		TotalEquipment:    a.TotalEquipment,
		TotalCost:         a.TotalCost,
		UnitPriceOfWork:   a.UnitPriceOfWork,
		Status:            a.Status,
	}

	if len(a.Items) > 0 {
		d.Items = make([]domain.AhspItem, len(a.Items))
		for i := range a.Items {
			d.Items[i] = *a.Items[i].ToDomain()
		}
	}

	return d
}

func FromAhspDomain(d *domain.Ahsp) *Ahsp {
	if d == nil {
		return nil
	}
	a := &Ahsp{
		Code:              d.Code,
		WorkItem:          d.WorkItem,
		Unit:              d.Unit,
		OverheadProfitPct: d.OverheadProfitPct,
		TotalLabor:        d.TotalLabor,
		TotalMaterial:     d.TotalMaterial,
		TotalEquipment:    d.TotalEquipment,
		TotalCost:         d.TotalCost,
		UnitPriceOfWork:   d.UnitPriceOfWork,
		Status:            d.Status,
	}
	a.ID = d.ID
	a.CreatedAt = d.CreatedAt
	a.UpdatedAt = d.UpdatedAt

	if len(d.Items) > 0 {
		a.Items = make([]AhspItem, len(d.Items))
		for i := range d.Items {
			a.Items[i] = *FromAhspItemDomain(&d.Items[i])
		}
	}

	return a
}

type AhspItem struct {
	gorm.Model
	AhspID      uint                     `gorm:"not null;index"`
	Category    domain.UnitPriceCategory `gorm:"type:text;not null"`
	Sort        int                      `gorm:"not null;default:0"`
	Description string                   `gorm:"not null;default:''"`
	Code        string                   `gorm:"not null;default:''"`
	Unit        string                   `gorm:"not null;default:''"`
	Coefficient float64                  `gorm:"type:numeric(14,4);not null;default:0"`
	UnitPrice   int64                    `gorm:"not null;default:0"`
	Cost        int64                    `gorm:"not null;default:0"`
	UnitPriceID *uint                    `gorm:"index"`
}

func (AhspItem) TableName() string {
	return "ahsp_items"
}

func (item *AhspItem) ToDomain() *domain.AhspItem {
	if item == nil {
		return nil
	}
	return &domain.AhspItem{
		ID:          item.ID,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
		AhspID:      item.AhspID,
		Category:    item.Category,
		Sort:        item.Sort,
		Description: item.Description,
		Code:        item.Code,
		Unit:        item.Unit,
		Coefficient: item.Coefficient,
		UnitPrice:   item.UnitPrice,
		Cost:        item.Cost,
		UnitPriceID: item.UnitPriceID,
	}
}

func FromAhspItemDomain(d *domain.AhspItem) *AhspItem {
	if d == nil {
		return nil
	}
	item := &AhspItem{
		AhspID:      d.AhspID,
		Category:    d.Category,
		Sort:        d.Sort,
		Description: d.Description,
		Code:        d.Code,
		Unit:        d.Unit,
		Coefficient: d.Coefficient,
		UnitPrice:   d.UnitPrice,
		Cost:        d.Cost,
		UnitPriceID: d.UnitPriceID,
	}
	item.ID = d.ID
	item.CreatedAt = d.CreatedAt
	item.UpdatedAt = d.UpdatedAt
	return item
}

