package domain

import "time"

type BoqItem struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ProjectID     uint      `json:"project_id"`
	AhspID        *uint     `json:"ahsp_id,omitempty"`
	WorkItem      string    `json:"work_item"`
	Unit          string    `json:"unit"`
	PlannedVolume float64   `json:"planned_volume"`

	Volumes []Volume `json:"volumes,omitempty"`
}
