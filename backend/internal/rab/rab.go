package rab

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/platform/httperr"
)

type RabItemDTO struct {
	ID        string  `json:"id"`
	ProjectID string  `json:"project_id"`
	BoqItemID *string `json:"boq_item_id"`
	WorkItem  string  `json:"work_item"`
	Unit      string  `json:"unit"`
	Volume    string  `json:"volume"`
	UnitPrice string  `json:"unit_price"`
	TotalCost string  `json:"total_cost"`
}

type RabResponseDTO struct {
	ProjectID string       `json:"project_id"`
	Items     []RabItemDTO `json:"items"`
	TotalRab  string       `json:"total_rab"`
}

type CreateRabItemInput struct {
	BoqItemID *string `json:"boq_item_id"`
	WorkItem  string  `json:"work_item"`
	Unit      string  `json:"unit"`
	Volume    string  `json:"volume"`
	UnitPrice string  `json:"unit_price"`
}

type UpdateRabItemInput struct {
	WorkItem  *string `json:"work_item"`
	Unit      *string `json:"unit"`
	Volume    *string `json:"volume"`
	UnitPrice *string `json:"unit_price"`
}

func roundHalfUp(val float64) int64 {
	if val < 0 {
		return int64(math.Ceil(val - 0.5))
	}
	return int64(math.Floor(val + 0.5))
}

func toDTO(item *domain.RabItem) RabItemDTO {
	dto := RabItemDTO{
		ID:        item.ID.String(),
		ProjectID: item.ProjectID.String(),
		WorkItem:  item.WorkItem,
		Unit:      item.Unit,
		Volume:    fmt.Sprintf("%.4f", item.Volume),
		UnitPrice: strconv.FormatInt(item.UnitPrice, 10),
		TotalCost: strconv.FormatInt(item.TotalCost, 10),
	}
	if item.BoqItemID != nil {
		s := item.BoqItemID.String()
		dto.BoqItemID = &s
	}
	return dto
}

func GetProjectRabHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		projectID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid project ID", nil)
		}

		q := strings.TrimSpace(c.Query("q"))
		query := db.Where("project_id = ?", projectID)
		if q != "" {
			like := "%" + strings.ToLower(q) + "%"
			query = query.Where("lower(work_item) LIKE ?", like)
		}

		var items []domain.RabItem
		if err := query.Order("created_at ASC").Find(&items).Error; err != nil {
			return httperr.Internal(c, "Failed to fetch RAB items")
		}

		var totalRab int64
		itemsDTO := make([]RabItemDTO, len(items))
		for i, item := range items {
			totalRab += item.TotalCost
			itemsDTO[i] = toDTO(&item)
		}

		return c.JSON(RabResponseDTO{
			ProjectID: projectID.String(),
			Items:     itemsDTO,
			TotalRab:  strconv.FormatInt(totalRab, 10),
		})
	}
}

func CreateHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		projectID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid project ID", nil)
		}

		var input CreateRabItemInput
		if err := c.Bind().Body(&input); err != nil {
			return httperr.BadRequest(c, "Invalid request body", nil)
		}
		if strings.TrimSpace(input.WorkItem) == "" || strings.TrimSpace(input.Unit) == "" {
			return httperr.BadRequest(c, "Work item and unit are required", nil)
		}

		vol, _ := strconv.ParseFloat(input.Volume, 64)
		price, _ := strconv.ParseInt(input.UnitPrice, 10, 64)
		totalCost := roundHalfUp(vol * float64(price))

		var boqID *uuid.UUID
		if input.BoqItemID != nil && *input.BoqItemID != "" {
			if parsed, err := uuid.Parse(*input.BoqItemID); err == nil {
				boqID = &parsed
			}
		}

		item := domain.RabItem{
			ProjectID: projectID,
			BoqItemID: boqID,
			WorkItem:  input.WorkItem,
			Unit:      input.Unit,
			Volume:    vol,
			UnitPrice: price,
			TotalCost: totalCost,
		}

		if err := db.Create(&item).Error; err != nil {
			return httperr.Internal(c, "Failed to create RAB item")
		}

		return c.Status(fiber.StatusCreated).JSON(toDTO(&item))
	}
}

func UpdateHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid RAB item ID", nil)
		}

		var item domain.RabItem
		if err := db.First(&item, id).Error; err != nil {
			return httperr.NotFound(c, "RAB item not found")
		}

		var input UpdateRabItemInput
		if err := c.Bind().Body(&input); err != nil {
			return httperr.BadRequest(c, "Invalid request body", nil)
		}

		if input.WorkItem != nil {
			item.WorkItem = *input.WorkItem
		}
		if input.Unit != nil {
			item.Unit = *input.Unit
		}
		if input.Volume != nil {
			v, _ := strconv.ParseFloat(*input.Volume, 64)
			item.Volume = v
		}
		if input.UnitPrice != nil {
			p, _ := strconv.ParseInt(*input.UnitPrice, 10, 64)
			item.UnitPrice = p
		}
		item.TotalCost = roundHalfUp(item.Volume * float64(item.UnitPrice))

		if err := db.Save(&item).Error; err != nil {
			return httperr.Internal(c, "Failed to update RAB item")
		}

		return c.JSON(toDTO(&item))
	}
}

func DeleteHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid RAB item ID", nil)
		}

		res := db.Delete(&domain.RabItem{}, id)
		if res.Error != nil {
			return httperr.Internal(c, "Failed to delete RAB item")
		}
		if res.RowsAffected == 0 {
			return httperr.NotFound(c, "RAB item not found")
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

