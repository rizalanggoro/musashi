package unitprice

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/platform/httperr"
)

type UnitPriceDTO struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Category    string `json:"category"`
	WorkItem    string `json:"work_item"`
	Unit        string `json:"unit"`
	UnitPrice   string `json:"unit_price"`
	Information string `json:"information"`
}

type CreateUnitPriceInput struct {
	Code        string `json:"code"`
	Category    string `json:"category"`
	WorkItem    string `json:"work_item"`
	Unit        string `json:"unit"`
	UnitPrice   string `json:"unit_price"`
	Information string `json:"information"`
}

type UpdateUnitPriceInput struct {
	Code        *string `json:"code"`
	Category    *string `json:"category"`
	WorkItem    *string `json:"work_item"`
	Unit        *string `json:"unit"`
	UnitPrice   *string `json:"unit_price"`
	Information *string `json:"information"`
}

func toDTO(u *domain.UnitPrice) UnitPriceDTO {
	return UnitPriceDTO{
		ID:          u.ID.String(),
		Code:        u.Code,
		Category:    string(u.Category),
		WorkItem:    u.WorkItem,
		Unit:        u.Unit,
		UnitPrice:   strconv.FormatInt(u.UnitPrice, 10),
		Information: u.Information,
	}
}

func ListHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		q := strings.TrimSpace(c.Query("q"))
		category := strings.TrimSpace(c.Query("category"))

		query := db.Model(&domain.UnitPrice{})
		if q != "" {
			like := "%" + strings.ToLower(q) + "%"
			query = query.Where("lower(code) LIKE ? OR lower(work_item) LIKE ?", like, like)
		}
		if category != "" {
			query = query.Where("category = ?", category)
		}

		var items []domain.UnitPrice
		if err := query.Order("category ASC, code ASC").Find(&items).Error; err != nil {
			return httperr.Internal(c, "Failed to fetch unit prices")
		}

		data := make([]UnitPriceDTO, len(items))
		for i := range items {
			data[i] = toDTO(&items[i])
		}
		return c.JSON(data)
	}
}

func CreateHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var input CreateUnitPriceInput
		if err := c.Bind().Body(&input); err != nil {
			return httperr.BadRequest(c, "Invalid request body", nil)
		}

		code := strings.TrimSpace(input.Code)
		if code == "" || input.WorkItem == "" || input.Unit == "" {
			return httperr.BadRequest(c, "Code, work item, and unit are required", nil)
		}

		var count int64
		db.Model(&domain.UnitPrice{}).Where("lower(code) = ?", strings.ToLower(code)).Count(&count)
		if count > 0 {
			return httperr.Conflict(c, fmt.Sprintf("Unit price with code '%s' already exists", code))
		}

		price, _ := strconv.ParseInt(input.UnitPrice, 10, 64)

		item := domain.UnitPrice{
			Code:        code,
			Category:    domain.UnitPriceCategory(input.Category),
			WorkItem:    input.WorkItem,
			Unit:        input.Unit,
			UnitPrice:   price,
			Information: input.Information,
		}

		if err := db.Create(&item).Error; err != nil {
			return httperr.Internal(c, "Failed to create unit price")
		}

		return c.Status(fiber.StatusCreated).JSON(toDTO(&item))
	}
}

func UpdateHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid unit price ID", nil)
		}

		var item domain.UnitPrice
		if err := db.First(&item, id).Error; err != nil {
			return httperr.NotFound(c, "Unit price not found")
		}

		var input UpdateUnitPriceInput
		if err := c.Bind().Body(&input); err != nil {
			return httperr.BadRequest(c, "Invalid request body", nil)
		}

		if input.Code != nil {
			newCode := strings.TrimSpace(*input.Code)
			if newCode != "" && !strings.EqualFold(newCode, item.Code) {
				var count int64
				db.Model(&domain.UnitPrice{}).Where("lower(code) = ? AND id != ?", strings.ToLower(newCode), id).Count(&count)
				if count > 0 {
					return httperr.Conflict(c, fmt.Sprintf("Unit price with code '%s' already exists", newCode))
				}
				item.Code = newCode
			}
		}
		if input.Category != nil {
			item.Category = domain.UnitPriceCategory(*input.Category)
		}
		if input.WorkItem != nil {
			item.WorkItem = *input.WorkItem
		}
		if input.Unit != nil {
			item.Unit = *input.Unit
		}
		if input.UnitPrice != nil {
			p, _ := strconv.ParseInt(*input.UnitPrice, 10, 64)
			item.UnitPrice = p
		}
		if input.Information != nil {
			item.Information = *input.Information
		}

		if err := db.Save(&item).Error; err != nil {
			return httperr.Internal(c, "Failed to update unit price")
		}

		return c.JSON(toDTO(&item))
	}
}

func DeleteHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid unit price ID", nil)
		}

		// Conflict check: is this unit price used in AHSP?
		var count int64
		db.Model(&domain.AhspItem{}).Where("unit_price_id = ?", id).Count(&count)
		if count > 0 {
			return httperr.Conflict(c, "Cannot delete unit price: referenced by one or more AHSP items")
		}

		res := db.Delete(&domain.UnitPrice{}, id)
		if res.Error != nil {
			return httperr.Internal(c, "Failed to delete unit price")
		}
		if res.RowsAffected == 0 {
			return httperr.NotFound(c, "Unit price not found")
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

