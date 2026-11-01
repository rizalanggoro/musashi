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

type ProjectSection string

const (
	ProjectSectionGABuilding ProjectSection = "GA Building"
	ProjectSectionEHS        ProjectSection = "EHS"
	ProjectSectionWellfare   ProjectSection = "Wellfare"
)

func IsValidSection(s string) bool {
	switch s {
	case string(ProjectSectionGABuilding), string(ProjectSectionEHS), string(ProjectSectionWellfare):
		return true
	default:
		return false
	}
}

type ProjectPlant string

const (
	ProjectPlantCikarang ProjectPlant = "Cikarang"
	ProjectPlantKarawang ProjectPlant = "Karawang"
)

func IsValidPlant(p string) bool {
	switch p {
	case string(ProjectPlantCikarang), string(ProjectPlantKarawang):
		return true
	default:
		return false
	}
}

type Project struct {
	ID                 uint           `json:"id"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	Name               string         `json:"name"`
	Section            string         `json:"section"`
	Plant              string         `json:"plant"`
	ImplementationPlan *time.Time     `json:"implementation_plan"`
	Budget             int64          `json:"budget"`
	Note               string         `json:"note"`
	Status               ProjectStatus  `json:"status"`
	ProgressPPT          ProgressStatus `json:"progress_ppt"`
	ProgressDrawing      ProgressStatus `json:"progress_drawing"`
	ProgressRAB          ProgressStatus `json:"progress_rab"`
	ProgressVendor       ProgressStatus `json:"progress_vendor"`
	ProgressSPB          ProgressStatus `json:"progress_spb"`
	ProgressPR           ProgressStatus `json:"progress_pr"`
	ProgressRINGI        ProgressStatus `json:"progress_ringi"`
	ProgressPO           ProgressStatus `json:"progress_po"`
	ProgressConstruction ProgressStatus `json:"progress_construction"`
	ProgressBAST         ProgressStatus `json:"progress_bast"`

	Schedules []Schedule `json:"schedules,omitempty"`
	BoqItems  []BoqItem  `json:"boq_items,omitempty"`
	RabItems  []RabItem  `json:"rab_items,omitempty"`
}
