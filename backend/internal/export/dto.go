package export

import (
	"musashi/backend/internal/domain"
)

type VolumeExportRow struct {
	domain.Volume
	WorkItem      string  `gorm:"column:work_item"`
	Unit          string  `gorm:"column:unit"`
	PlannedVolume float64 `gorm:"column:planned_volume"`
}
