package model

import (
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type BoqItem struct {
	gorm.Model
	ProjectID     uint    `gorm:"not null;index"`
	AhspID        *uint   `gorm:"index"`
	WorkItem      string  `gorm:"not null"`
	Unit          string  `gorm:"not null"`
	PlannedVolume float64 `gorm:"type:numeric(14,4);not null;default:0"`

	Volumes []Volume `gorm:"foreignKey:BoqItemID"`
}

func (BoqItem) TableName() string {
	return "boq_items"
}

func (b *BoqItem) ToDomain() *domain.BoqItem {
	if b == nil {
		return nil
	}
	d := &domain.BoqItem{
		ID:            b.ID,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
		ProjectID:     b.ProjectID,
		AhspID:        b.AhspID,
		WorkItem:      b.WorkItem,
		Unit:          b.Unit,
		PlannedVolume: b.PlannedVolume,
	}

	if len(b.Volumes) > 0 {
		d.Volumes = make([]domain.Volume, len(b.Volumes))
		for i := range b.Volumes {
			d.Volumes[i] = *b.Volumes[i].ToDomain()
		}
	}

	return d
}

func FromBoqItemDomain(d *domain.BoqItem) *BoqItem {
	if d == nil {
		return nil
	}
	b := &BoqItem{
		ProjectID:     d.ProjectID,
		AhspID:        d.AhspID,
		WorkItem:      d.WorkItem,
		Unit:          d.Unit,
		PlannedVolume: d.PlannedVolume,
	}
	b.ID = d.ID
	b.CreatedAt = d.CreatedAt
	b.UpdatedAt = d.UpdatedAt

	if len(d.Volumes) > 0 {
		b.Volumes = make([]Volume, len(d.Volumes))
		for i := range d.Volumes {
			b.Volumes[i] = *FromVolumeDomain(&d.Volumes[i])
		}
	}

	return b
}

