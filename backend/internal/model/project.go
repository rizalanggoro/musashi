package model

import (
	"time"

	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type Project struct {
	gorm.Model
	Name               string                `gorm:"not null"`
	Section            string                `gorm:"not null;default:''"`
	Plant              string                `gorm:"not null;default:''"`
	ImplementationPlan *time.Time            `gorm:"type:date"`
	Budget             int64                 `gorm:"not null;default:0"`
	Information        string                `gorm:"not null;default:''"`
	Status             domain.ProjectStatus  `gorm:"type:text;default:'planning'"`
	ProgressSPB        domain.ProgressStatus `gorm:"column:progress_spb;type:text;default:'pending'"`
	ProgressRAB        domain.ProgressStatus `gorm:"column:progress_rab;type:text;default:'pending'"`
	ProgressXXX        domain.ProgressStatus `gorm:"column:progress_xxx;type:text;default:'pending'"`

	Schedules []Schedule `gorm:"foreignKey:ProjectID"`
	BoqItems  []BoqItem  `gorm:"foreignKey:ProjectID"`
	RabItems  []RabItem  `gorm:"foreignKey:ProjectID"`
}

func (Project) TableName() string {
	return "projects"
}

func (p *Project) ToDomain() *domain.Project {
	if p == nil {
		return nil
	}
	d := &domain.Project{
		ID:                 p.ID,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
		Name:               p.Name,
		Section:            p.Section,
		Plant:              p.Plant,
		ImplementationPlan: p.ImplementationPlan,
		Budget:             p.Budget,
		Information:        p.Information,
		Status:             p.Status,
		ProgressSPB:        p.ProgressSPB,
		ProgressRAB:        p.ProgressRAB,
		ProgressXXX:        p.ProgressXXX,
	}

	if len(p.Schedules) > 0 {
		d.Schedules = make([]domain.Schedule, len(p.Schedules))
		for i := range p.Schedules {
			d.Schedules[i] = *p.Schedules[i].ToDomain()
		}
	}
	if len(p.BoqItems) > 0 {
		d.BoqItems = make([]domain.BoqItem, len(p.BoqItems))
		for i := range p.BoqItems {
			d.BoqItems[i] = *p.BoqItems[i].ToDomain()
		}
	}
	if len(p.RabItems) > 0 {
		d.RabItems = make([]domain.RabItem, len(p.RabItems))
		for i := range p.RabItems {
			d.RabItems[i] = *p.RabItems[i].ToDomain()
		}
	}

	return d
}

func FromProjectDomain(d *domain.Project) *Project {
	if d == nil {
		return nil
	}
	p := &Project{
		Name:               d.Name,
		Section:            d.Section,
		Plant:              d.Plant,
		ImplementationPlan: d.ImplementationPlan,
		Budget:             d.Budget,
		Information:        d.Information,
		Status:             d.Status,
		ProgressSPB:        d.ProgressSPB,
		ProgressRAB:        d.ProgressRAB,
		ProgressXXX:        d.ProgressXXX,
	}
	p.ID = d.ID
	p.CreatedAt = d.CreatedAt
	p.UpdatedAt = d.UpdatedAt

	if len(d.Schedules) > 0 {
		p.Schedules = make([]Schedule, len(d.Schedules))
		for i := range d.Schedules {
			p.Schedules[i] = *FromScheduleDomain(&d.Schedules[i])
		}
	}
	if len(d.BoqItems) > 0 {
		p.BoqItems = make([]BoqItem, len(d.BoqItems))
		for i := range d.BoqItems {
			p.BoqItems[i] = *FromBoqItemDomain(&d.BoqItems[i])
		}
	}
	if len(d.RabItems) > 0 {
		p.RabItems = make([]RabItem, len(d.RabItems))
		for i := range d.RabItems {
			p.RabItems[i] = *FromRabItemDomain(&d.RabItems[i])
		}
	}

	return p
}

