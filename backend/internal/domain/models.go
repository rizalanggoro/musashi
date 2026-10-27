package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleStaff UserRole = "staff"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	Name         string    `gorm:"not null" json:"name"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         UserRole  `gorm:"type:text;default:'staff'" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

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
	ID                 uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name               string         `gorm:"not null" json:"name"`
	Section            string         `gorm:"not null;default:''" json:"section"`
	Plant              string         `gorm:"not null;default:''" json:"plant"`
	ImplementationPlan *time.Time     `gorm:"type:date" json:"implementation_plan"`
	Budget             int64          `gorm:"not null;default:0" json:"budget"`
	Information        string         `gorm:"not null;default:''" json:"information"`
	Status             ProjectStatus  `gorm:"type:text;default:'planning'" json:"status"`
	ProgressSPB        ProgressStatus `gorm:"column:progress_spb;type:text;default:'pending'" json:"progress_spb"`
	ProgressRAB        ProgressStatus `gorm:"column:progress_rab;type:text;default:'pending'" json:"progress_rab"`
	ProgressXXX        ProgressStatus `gorm:"column:progress_xxx;type:text;default:'pending'" json:"progress_xxx"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`

	Schedules []Schedule `gorm:"foreignKey:ProjectID" json:"schedules,omitempty"`
	BoqItems  []BoqItem  `gorm:"foreignKey:ProjectID" json:"boq_items,omitempty"`
	RabItems  []RabItem  `gorm:"foreignKey:ProjectID" json:"rab_items,omitempty"`
}

type ScheduleStatus string

const (
	SchedulePlanning ScheduleStatus = "planning"
	ScheduleOnGoing  ScheduleStatus = "on_going"
	ScheduleDone     ScheduleStatus = "done"
	ScheduleLate     ScheduleStatus = "late"
)

type Schedule struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProjectID uuid.UUID      `gorm:"type:uuid;not null;index" json:"project_id"`
	Code      string         `gorm:"not null;default:''" json:"code"`
	WorkItem  string         `gorm:"not null" json:"work_item"`
	StartDate time.Time      `gorm:"type:date;not null" json:"start_date"`
	EndDate   time.Time      `gorm:"type:date;not null" json:"end_date"`
	Duration  int            `gorm:"not null;default:1" json:"duration"`
	Weight    float64        `gorm:"type:numeric(6,2);not null;default:0" json:"weight"`
	Progress  float64        `gorm:"type:numeric(6,2);not null;default:0" json:"progress"`
	Status    ScheduleStatus `gorm:"type:text;default:'planning'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type UnitPriceCategory string

const (
	CategoryLabor     UnitPriceCategory = "LABOR"
	CategoryMaterials UnitPriceCategory = "MATERIALS"
	CategoryEquipment UnitPriceCategory = "EQUIPMENT"
	CategoryOther     UnitPriceCategory = "OTHER"
)

type UnitPrice struct {
	ID          uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Code        string            `gorm:"uniqueIndex;not null" json:"code"`
	Category    UnitPriceCategory `gorm:"type:text;not null;index" json:"category"`
	WorkItem    string            `gorm:"not null" json:"work_item"`
	Unit        string            `gorm:"not null" json:"unit"`
	UnitPrice   int64             `gorm:"not null;default:0" json:"unit_price"`
	Information string            `gorm:"not null;default:''" json:"information"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type AhspStatus string

const (
	AhspStatusDraft AhspStatus = "draft"
	AhspStatusFinal AhspStatus = "final"
)

type Ahsp struct {
	ID                uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Code              string     `gorm:"uniqueIndex;not null" json:"code"`
	WorkItem          string     `gorm:"not null" json:"work_item"`
	Unit              string     `gorm:"not null" json:"unit"`
	OverheadProfitPct float64    `gorm:"type:numeric(6,2);not null;default:10.00" json:"overhead_profit_pct"`
	TotalLabor        int64      `gorm:"not null;default:0" json:"total_labor"`
	TotalMaterial     int64      `gorm:"not null;default:0" json:"total_material"`
	TotalEquipment    int64      `gorm:"not null;default:0" json:"total_equipment"`
	TotalCost         int64      `gorm:"not null;default:0" json:"total_cost"`
	UnitPriceOfWork   int64      `gorm:"not null;default:0" json:"unit_price_of_work"`
	Status            AhspStatus `gorm:"type:text;default:'draft'" json:"status"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	Items []AhspItem `gorm:"foreignKey:AhspID" json:"items,omitempty"`
}

type AhspItem struct {
	ID          uuid.UUID          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AhspID      uuid.UUID          `gorm:"type:uuid;not null;index" json:"ahsp_id"`
	Category    UnitPriceCategory  `gorm:"type:text;not null" json:"category"`
	Sort        int                `gorm:"not null;default:0" json:"sort"`
	Description string             `gorm:"not null;default:''" json:"description"`
	Code        string             `gorm:"not null;default:''" json:"code"`
	Unit        string             `gorm:"not null;default:''" json:"unit"`
	Coefficient float64            `gorm:"type:numeric(14,4);not null;default:0" json:"coefficient"`
	UnitPrice   int64              `gorm:"not null;default:0" json:"unit_price"`
	Cost        int64              `gorm:"not null;default:0" json:"cost"`
	UnitPriceID *uuid.UUID         `gorm:"type:uuid;index" json:"unit_price_id,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type BoqItem struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProjectID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"project_id"`
	AhspID        *uuid.UUID `gorm:"type:uuid" json:"ahsp_id,omitempty"`
	WorkItem      string     `gorm:"not null" json:"work_item"`
	Unit          string     `gorm:"not null" json:"unit"`
	PlannedVolume float64    `gorm:"type:numeric(14,4);not null;default:0" json:"planned_volume"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	Volumes []Volume `gorm:"foreignKey:BoqItemID" json:"volumes,omitempty"`
}

type Volume struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BoqItemID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"boq_item_id"`
	Date          *time.Time `gorm:"type:date" json:"date"`
	Length        *float64   `gorm:"type:numeric(14,4)" json:"length"`
	Width         *float64   `gorm:"type:numeric(14,4)" json:"width"`
	Height        *float64   `gorm:"type:numeric(14,4)" json:"height"`
	Quantity      float64    `gorm:"type:numeric(14,4);not null;default:0" json:"quantity"`
	UnitWeight    float64    `gorm:"type:numeric(14,4);not null;default:1" json:"unit_weight"`
	TotalQuantity float64    `gorm:"type:numeric(14,4);not null;default:0" json:"total_quantity"`
	Description   string     `gorm:"not null;default:''" json:"description"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type RabItem struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProjectID uuid.UUID  `gorm:"type:uuid;not null;index" json:"project_id"`
	BoqItemID *uuid.UUID `gorm:"type:uuid" json:"boq_item_id,omitempty"`
	WorkItem  string     `gorm:"not null" json:"work_item"`
	Unit      string     `gorm:"not null" json:"unit"`
	Volume    float64    `gorm:"type:numeric(14,4);not null;default:0" json:"volume"`
	UnitPrice int64      `gorm:"not null;default:0" json:"unit_price"`
	TotalCost int64      `gorm:"not null;default:0" json:"total_cost"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

