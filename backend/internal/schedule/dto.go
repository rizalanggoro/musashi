package schedule

import (
	"fmt"

	"musashi/backend/internal/domain"
	"musashi/backend/utils"
)

type ScheduleDTO struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Code      string `json:"code"`
	WorkItem  string `json:"work_item"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Duration  int    `json:"duration"`
	Weight    string `json:"weight"`
	Progress  string `json:"progress"`
	Status    string `json:"status"`
}

type CreateScheduleInput struct {
	Code      string `json:"code"`
	WorkItem  string `json:"work_item"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Weight    string `json:"weight"`
	Progress  string `json:"progress"`
	Status    string `json:"status"`
}

type UpdateScheduleInput struct {
	Code      *string `json:"code,omitempty"`
	WorkItem  *string `json:"work_item,omitempty"`
	StartDate *string `json:"start_date,omitempty"`
	EndDate   *string `json:"end_date,omitempty"`
	Weight    *string `json:"weight,omitempty"`
	Progress  *string `json:"progress,omitempty"`
	Status    *string `json:"status,omitempty"`
}

func toDTO(s *domain.Schedule) ScheduleDTO {
	return ScheduleDTO{
		ID:        utils.FormatID(s.ID),
		ProjectID: utils.FormatID(s.ProjectID),
		Code:      s.Code,
		WorkItem:  s.WorkItem,
		StartDate: s.StartDate.Format("2006-01-02"),
		EndDate:   s.EndDate.Format("2006-01-02"),
		Duration:  s.Duration,
		Weight:    fmt.Sprintf("%.2f", s.Weight),
		Progress:  fmt.Sprintf("%.2f", s.Progress),
		Status:    string(s.Status),
	}
}
