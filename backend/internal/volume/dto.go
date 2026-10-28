package volume

import (
	"fmt"

	"musashi/backend/internal/domain"
)

type VolumeRowDTO struct {
	ID            string  `json:"id"`
	BoqItemID     string  `json:"boq_item_id"`
	WorkItem      string  `json:"work_item"`
	Unit          string  `json:"unit"`
	PlannedVolume string  `json:"planned_volume"`
	Date          *string `json:"date"`
	Length        *string `json:"length"`
	Width         *string `json:"width"`
	Height        *string `json:"height"`
	Quantity      string  `json:"quantity"`
	UnitWeight    string  `json:"unit_weight"`
	TotalQuantity string  `json:"total_quantity"`
	Description   string  `json:"description"`
}

type CreateVolumeInput struct {
	BoqItemID     *string `json:"boq_item_id"`
	AhspID        *string `json:"ahsp_id"`
	WorkItem      string  `json:"work_item"`
	Unit          string  `json:"unit"`
	PlannedVolume string  `json:"planned_volume"`
	Date          string  `json:"date"`
	Length        *string `json:"length"`
	Width         *string `json:"width"`
	Height        *string `json:"height"`
	UnitWeight    string  `json:"unit_weight"`
	Description   string  `json:"description"`
}

type UpdateVolumeInput struct {
	Date        *string `json:"date"`
	Length      *string `json:"length"`
	Width       *string `json:"width"`
	Height      *string `json:"height"`
	UnitWeight  *string `json:"unit_weight"`
	Description *string `json:"description"`
}

type VolumeRowQueryResult struct {
	domain.Volume
	WorkItem      string  `gorm:"column:work_item"`
	Unit          string  `gorm:"column:unit"`
	PlannedVolume float64 `gorm:"column:planned_volume"`
}

func toRowDTO(v *domain.Volume, b *domain.BoqItem) VolumeRowDTO {
	dto := VolumeRowDTO{
		ID:            v.ID.String(),
		BoqItemID:     v.BoqItemID.String(),
		WorkItem:      b.WorkItem,
		Unit:          b.Unit,
		PlannedVolume: fmt.Sprintf("%.4f", b.PlannedVolume),
		Quantity:      fmt.Sprintf("%.4f", v.Quantity),
		UnitWeight:    fmt.Sprintf("%.4f", v.UnitWeight),
		TotalQuantity: fmt.Sprintf("%.4f", v.TotalQuantity),
		Description:   v.Description,
	}
	if v.Date != nil {
		d := v.Date.Format("2006-01-02")
		dto.Date = &d
	}
	if v.Length != nil {
		l := fmt.Sprintf("%.4f", *v.Length)
		dto.Length = &l
	}
	if v.Width != nil {
		w := fmt.Sprintf("%.4f", *v.Width)
		dto.Width = &w
	}
	if v.Height != nil {
		h := fmt.Sprintf("%.4f", *v.Height)
		dto.Height = &h
	}
	return dto
}
