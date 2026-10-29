package domain

import "time"

type ProjectStatus string

const (
	ProjectStatusPlanning  ProjectStatus = "planning"
	ProjectStatusOnGoing   ProjectStatus = "on_going"
	ProjectStatusCompleted ProjectStatus = "completed"
	ProjectStatusCancelled ProjectStatus = "cancelled"
)

type ProgressStatus string

const (
	ProgressPending ProgressStatus = "pending"
	ProgressDone    ProgressStatus = "done"
)

type Project struct {
	ID                 uint           `json:"id"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	Name               string         `json:"name"`
	Section            string         `json:"section"`
	Plant              string         `json:"plant"`
	ImplementationPlan *time.Time     `json:"implementation_plan"`
	Budget             int64          `json:"budget"`
	Information        string         `json:"information"`
	Status             ProjectStatus  `json:"status"`
	ProgressSPB        ProgressStatus `json:"progress_spb"`
	ProgressRAB        ProgressStatus `json:"progress_rab"`
	ProgressXXX        ProgressStatus `json:"progress_xxx"`

	Schedules []Schedule `json:"schedules,omitempty"`
	BoqItems  []BoqItem  `json:"boq_items,omitempty"`
	RabItems  []RabItem  `json:"rab_items,omitempty"`
}
