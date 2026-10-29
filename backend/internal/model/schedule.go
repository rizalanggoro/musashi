package model

import (
	"time"

	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type Schedule struct {
	gorm.Model
	ProjectID uint                  `gorm:"not null;index"`
	Code      string                `gorm:"not null;default:''"`
	WorkItem  string                `gorm:"not null"`
	StartDate time.Time             `gorm:"type:date;not null"`
	EndDate   time.Time             `gorm:"type:date;not null"`
	Duration  int                   `gorm:"not null;default:1"`
	Weight    float64               `gorm:"type:numeric(6,2);not null;default:0"`
	Progress  float64               `gorm:"type:numeric(6,2);not null;default:0"`
	Status    domain.ScheduleStatus `gorm:"type:text;default:'planning'"`
}

func (Schedule) TableName() string {
	return "schedules"
}

func (s *Schedule) ToDomain() *domain.Schedule {
	if s == nil {
		return nil
	}
	return &domain.Schedule{
		ID:        s.ID,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
		ProjectID: s.ProjectID,
		Code:      s.Code,
		WorkItem:  s.WorkItem,
		StartDate: s.StartDate,
		EndDate:   s.EndDate,
		Duration:  s.Duration,
		Weight:    s.Weight,
		Progress:  s.Progress,
		Status:    s.Status,
	}
}

func FromScheduleDomain(d *domain.Schedule) *Schedule {
	if d == nil {
		return nil
	}
	s := &Schedule{
		ProjectID: d.ProjectID,
		Code:      d.Code,
		WorkItem:  d.WorkItem,
		StartDate: d.StartDate,
		EndDate:   d.EndDate,
		Duration:  d.Duration,
		Weight:    d.Weight,
		Progress:  d.Progress,
		Status:    d.Status,
	}
	s.ID = d.ID
	s.CreatedAt = d.CreatedAt
	s.UpdatedAt = d.UpdatedAt
	return s
}

