package volume

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/platform/httperr"
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

func calcQuantities(length, width, height *float64, plannedVolume float64, unitWeight float64) (float64, float64) {
	if unitWeight <= 0 {
		unitWeight = 1
	}

	var qty float64
	if length != nil && width != nil && height != nil && *length > 0 && *width > 0 && *height > 0 {
		qty = (*length) * (*width) * (*height)
	} else {
		qty = plannedVolume
	}

	totalQty := qty * unitWeight
	return qty, totalQty
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

func ListHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		projectID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid project ID", nil)
		}

		q := strings.TrimSpace(c.Query("q"))

		var rows []struct {
			domain.Volume
			WorkItem      string  `gorm:"column:work_item"`
			Unit          string  `gorm:"column:unit"`
			PlannedVolume float64 `gorm:"column:planned_volume"`
		}

		query := db.Table("volumes").
			Select("volumes.*, boq_items.work_item, boq_items.unit, boq_items.planned_volume").
			Joins("JOIN boq_items ON boq_items.id = volumes.boq_item_id").
			Where("boq_items.project_id = ?", projectID)

		if q != "" {
			like := "%" + strings.ToLower(q) + "%"
			query = query.Where("lower(boq_items.work_item) LIKE ? OR lower(volumes.description) LIKE ?", like, like)
		}

		if err := query.Order("boq_items.created_at ASC, volumes.created_at ASC").Scan(&rows).Error; err != nil {
			return httperr.Internal(c, "Failed to fetch volume take-offs")
		}

		data := make([]VolumeRowDTO, len(rows))
		for i, r := range rows {
			b := domain.BoqItem{
				WorkItem:      r.WorkItem,
				Unit:          r.Unit,
				PlannedVolume: r.PlannedVolume,
			}
			data[i] = toRowDTO(&r.Volume, &b)
		}

		return c.JSON(data)
	}
}

func CreateHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		projectID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid project ID", nil)
		}

		var input CreateVolumeInput
		if err := c.Bind().Body(&input); err != nil {
			return httperr.BadRequest(c, "Invalid request body", nil)
		}

		var boqItem domain.BoqItem
		if input.BoqItemID != nil && *input.BoqItemID != "" {
			bID, err := uuid.Parse(*input.BoqItemID)
			if err != nil {
				return httperr.BadRequest(c, "Invalid BOQ item ID", nil)
			}
			if err := db.Where("id = ? AND project_id = ?", bID, projectID).First(&boqItem).Error; err != nil {
				return httperr.NotFound(c, "BOQ item not found")
			}
		} else {
			if strings.TrimSpace(input.WorkItem) == "" || strings.TrimSpace(input.Unit) == "" {
				return httperr.BadRequest(c, "Work item and unit are required", nil)
			}
			plannedVol, _ := strconv.ParseFloat(input.PlannedVolume, 64)

			var ahspID *uuid.UUID
			if input.AhspID != nil && *input.AhspID != "" {
				if parsed, err := uuid.Parse(*input.AhspID); err == nil {
					ahspID = &parsed
				}
			}

			boqItem = domain.BoqItem{
				ProjectID:     projectID,
				AhspID:        ahspID,
				WorkItem:      input.WorkItem,
				Unit:          input.Unit,
				PlannedVolume: plannedVol,
			}
			if err := db.Create(&boqItem).Error; err != nil {
				return httperr.Internal(c, "Failed to create BOQ item")
			}
		}

		var lVal, wVal, hVal *float64
		if input.Length != nil && *input.Length != "" {
			v, _ := strconv.ParseFloat(*input.Length, 64)
			lVal = &v
		}
		if input.Width != nil && *input.Width != "" {
			v, _ := strconv.ParseFloat(*input.Width, 64)
			wVal = &v
		}
		if input.Height != nil && *input.Height != "" {
			v, _ := strconv.ParseFloat(*input.Height, 64)
			hVal = &v
		}

		uWeight := 1.0
		if input.UnitWeight != "" {
			if w, err := strconv.ParseFloat(input.UnitWeight, 64); err == nil && w > 0 {
				uWeight = w
			}
		}

		qty, totalQty := calcQuantities(lVal, wVal, hVal, boqItem.PlannedVolume, uWeight)

		var dateVal *time.Time
		if input.Date != "" {
			if t, err := time.Parse("2006-01-02", input.Date); err == nil {
				dateVal = &t
			}
		}

		vol := domain.Volume{
			BoqItemID:     boqItem.ID,
			Date:          dateVal,
			Length:        lVal,
			Width:         wVal,
			Height:        hVal,
			Quantity:      qty,
			UnitWeight:    uWeight,
			TotalQuantity: totalQty,
			Description:   input.Description,
		}

		if err := db.Create(&vol).Error; err != nil {
			return httperr.Internal(c, "Failed to create volume take-off")
		}

		return c.Status(fiber.StatusCreated).JSON(toRowDTO(&vol, &boqItem))
	}
}

func UpdateHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid volume ID", nil)
		}

		var vol domain.Volume
		if err := db.First(&vol, id).Error; err != nil {
			return httperr.NotFound(c, "Volume row not found")
		}

		var boqItem domain.BoqItem
		if err := db.First(&boqItem, vol.BoqItemID).Error; err != nil {
			return httperr.NotFound(c, "Associated BOQ item not found")
		}

		var input UpdateVolumeInput
		if err := c.Bind().Body(&input); err != nil {
			return httperr.BadRequest(c, "Invalid request body", nil)
		}

		if input.Date != nil {
			if *input.Date != "" {
				if t, err := time.Parse("2006-01-02", *input.Date); err == nil {
					vol.Date = &t
				}
			} else {
				vol.Date = nil
			}
		}
		if input.Length != nil {
			if *input.Length != "" {
				v, _ := strconv.ParseFloat(*input.Length, 64)
				vol.Length = &v
			} else {
				vol.Length = nil
			}
		}
		if input.Width != nil {
			if *input.Width != "" {
				v, _ := strconv.ParseFloat(*input.Width, 64)
				vol.Width = &v
			} else {
				vol.Width = nil
			}
		}
		if input.Height != nil {
			if *input.Height != "" {
				v, _ := strconv.ParseFloat(*input.Height, 64)
				vol.Height = &v
			} else {
				vol.Height = nil
			}
		}
		if input.UnitWeight != nil {
			w, _ := strconv.ParseFloat(*input.UnitWeight, 64)
			if w <= 0 {
				w = 1.0
			}
			vol.UnitWeight = w
		}
		if input.Description != nil {
			vol.Description = *input.Description
		}

		vol.Quantity, vol.TotalQuantity = calcQuantities(vol.Length, vol.Width, vol.Height, boqItem.PlannedVolume, vol.UnitWeight)

		if err := db.Save(&vol).Error; err != nil {
			return httperr.Internal(c, "Failed to update volume row")
		}

		return c.JSON(toRowDTO(&vol, &boqItem))
	}
}

func DeleteHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid volume ID", nil)
		}

		res := db.Delete(&domain.Volume{}, id)
		if res.Error != nil {
			return httperr.Internal(c, "Failed to delete volume row")
		}
		if res.RowsAffected == 0 {
			return httperr.NotFound(c, "Volume row not found")
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

