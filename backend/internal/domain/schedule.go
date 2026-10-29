package domain

import "time"

type ScheduleStatus string

const (
	SchedulePlanning ScheduleStatus = "planning"
	ScheduleOnGoing  ScheduleStatus = "on_going"
	ScheduleDone     ScheduleStatus = "done"
	ScheduleLate     ScheduleStatus = "late"
)

type Schedule struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	ProjectID uint           `json:"project_id"`
	Code      string         `json:"code"`
	WorkItem  string         `json:"work_item"`
	StartDate time.Time      `json:"start_date"`
	EndDate   time.Time      `json:"end_date"`
	Duration  int            `json:"duration"`
	Weight    float64        `json:"weight"`
	Progress  float64        `json:"progress"`
	Status    ScheduleStatus `json:"status"`
}
