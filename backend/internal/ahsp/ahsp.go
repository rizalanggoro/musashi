package ahsp

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

type AhspDTO struct {
	ID                string `json:"id"`
	Code              string `json:"code"`
	WorkItem          string `json:"work_item"`
	Unit              string `json:"unit"`
	OverheadProfitPct string `json:"overhead_profit_pct"`
	TotalLabor        string `json:"total_labor"`
	TotalMaterial     string `json:"total_material"`
	TotalEquipment    string `json:"total_equipment"`
	TotalCost         string `json:"total_cost"`
	UnitPriceOfWork   string `json:"unit_price_of_work"`
	Status            string `json:"status"`
}

type AhspItemDTO struct {
	ID          string  `json:"id"`
	UnitPriceID *string `json:"unit_price_id"`
	Category    string  `json:"category"`
	Sort        int     `json:"sort"`
	Description string  `json:"description"`
	Code        string  `json:"code"`
	Unit        string  `json:"unit"`
	Coefficient string  `json:"coefficient"`
	UnitPrice   string  `json:"unit_price"`
	Cost        string  `json:"cost"`
}

type AhspDetailDTO struct {
	AhspDTO
	Items []AhspItemDTO `json:"items"`
}

type CreateAhspItemInput struct {
	UnitPriceID *string `json:"unit_price_id"`
	Category    string  `json:"category"`
	Sort        int     `json:"sort"`
	Description string  `json:"description"`
	Code        string  `json:"code"`
	Unit        string  `json:"unit"`
	Coefficient string  `json:"coefficient"`
	UnitPrice   string  `json:"unit_price"`
}

type CreateAhspInput struct {
	Code              string                `json:"code"`
	WorkItem          string                `json:"work_item"`
	Unit              string                `json:"unit"`
	OverheadProfitPct string                `json:"overhead_profit_pct"`
	Status            string                `json:"status"`
	Items             []CreateAhspItemInput `json:"items"`
}

func toAhspDTO(a *domain.Ahsp) AhspDTO {
	return AhspDTO{
		ID:                a.ID.String(),
		Code:              a.Code,
		WorkItem:          a.WorkItem,
		Unit:              a.Unit,
		OverheadProfitPct: fmt.Sprintf("%.2f", a.OverheadProfitPct),
		TotalLabor:        strconv.FormatInt(a.TotalLabor, 10),
		TotalMaterial:     strconv.FormatInt(a.TotalMaterial, 10),
		TotalEquipment:    strconv.FormatInt(a.TotalEquipment, 10),
		TotalCost:         strconv.FormatInt(a.TotalCost, 10),
		UnitPriceOfWork:   strconv.FormatInt(a.UnitPriceOfWork, 10),
		Status:            string(a.Status),
	}
}

func toItemDTO(item *domain.AhspItem) AhspItemDTO {
	dto := AhspItemDTO{
		ID:          item.ID.String(),
		Category:    string(item.Category),
		Sort:        item.Sort,
		Description: item.Description,
		Code:        item.Code,
		Unit:        item.Unit,
		Coefficient: fmt.Sprintf("%.4f", item.Coefficient),
		UnitPrice:   strconv.FormatInt(item.UnitPrice, 10),
		Cost:        strconv.FormatInt(item.Cost, 10),
	}
	if item.UnitPriceID != nil {
		s := item.UnitPriceID.String()
		dto.UnitPriceID = &s
	}
	return dto
}

func roundHalfUp(val float64) int64 {
	if val < 0 {
		return int64(math.Ceil(val - 0.5))
	}
	return int64(math.Floor(val + 0.5))
}

func ListHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		q := strings.TrimSpace(c.Query("q"))
		query := db.Model(&domain.Ahsp{})
		if q != "" {
			like := "%" + strings.ToLower(q) + "%"
			query = query.Where("lower(code) LIKE ? OR lower(work_item) LIKE ?", like, like)
		}

		var ahsps []domain.Ahsp
		if err := query.Order("code ASC").Find(&ahsps).Error; err != nil {
			return httperr.Internal(c, "Failed to fetch AHSPs")
		}

		data := make([]AhspDTO, len(ahsps))
		for i := range ahsps {
			data[i] = toAhspDTO(&ahsps[i])
		}
		return c.JSON(data)
	}
}

func GetHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid AHSP ID", nil)
		}

		var a domain.Ahsp
		if err := db.Preload("Items", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("sort ASC, created_at ASC")
		}).First(&a, id).Error; err != nil {
			return httperr.NotFound(c, "AHSP not found")
		}

		itemsDTO := make([]AhspItemDTO, len(a.Items))
		for i := range a.Items {
			itemsDTO[i] = toItemDTO(&a.Items[i])
		}

		return c.JSON(AhspDetailDTO{
			AhspDTO: toAhspDTO(&a),
			Items:   itemsDTO,
		})
	}
}

func calculateAndBuild(input CreateAhspInput) (*domain.Ahsp, []domain.AhspItem, error) {
	code := strings.TrimSpace(input.Code)
	if code == "" || input.WorkItem == "" || input.Unit == "" {
		return nil, nil, fmt.Errorf("code, work item, and unit are required")
	}

	overheadPct := 10.0
	if input.OverheadProfitPct != "" {
		if p, err := strconv.ParseFloat(input.OverheadProfitPct, 64); err == nil {
			overheadPct = p
		}
	}

	var totalLabor int64
	var totalMaterial int64
	var totalEquipment int64

	items := make([]domain.AhspItem, len(input.Items))
	for i, inItem := range input.Items {
		coef, _ := strconv.ParseFloat(inItem.Coefficient, 64)
		price, _ := strconv.ParseInt(inItem.UnitPrice, 10, 64)
		cost := roundHalfUp(coef * float64(price))

		cat := domain.UnitPriceCategory(inItem.Category)
		switch cat {
		case domain.CategoryLabor:
			totalLabor += cost
		case domain.CategoryMaterials:
			totalMaterial += cost
		case domain.CategoryEquipment:
			totalEquipment += cost
		}

		var upID *uuid.UUID
		if inItem.UnitPriceID != nil && *inItem.UnitPriceID != "" {
			if parsed, err := uuid.Parse(*inItem.UnitPriceID); err == nil {
				upID = &parsed
			}
		}

		items[i] = domain.AhspItem{
			Category:    cat,
			Sort:        inItem.Sort,
			Description: inItem.Description,
			Code:        inItem.Code,
			Unit:        inItem.Unit,
			Coefficient: coef,
			UnitPrice:   price,
			Cost:        cost,
			UnitPriceID: upID,
		}
	}

	totalCostD := totalLabor + totalMaterial + totalEquipment
	overheadCostE := roundHalfUp((overheadPct / 100.0) * float64(totalCostD))
	unitPriceOfWorkF := totalCostD + overheadCostE

	status := domain.AhspStatusDraft
	if input.Status == "final" {
		status = domain.AhspStatusFinal
	}

	ahsp := &domain.Ahsp{
		Code:              code,
		WorkItem:          input.WorkItem,
		Unit:              input.Unit,
		OverheadProfitPct: overheadPct,
		TotalLabor:        totalLabor,
		TotalMaterial:     totalMaterial,
		TotalEquipment:    totalEquipment,
		TotalCost:         totalCostD,
		UnitPriceOfWork:   unitPriceOfWorkF,
		Status:            status,
	}

	return ahsp, items, nil
}

func CreateHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var input CreateAhspInput
		if err := c.Bind().Body(&input); err != nil {
			return httperr.BadRequest(c, "Invalid request body", nil)
		}

		ahsp, items, err := calculateAndBuild(input)
		if err != nil {
			return httperr.BadRequest(c, err.Error(), nil)
		}

		var count int64
		db.Model(&domain.Ahsp{}).Where("lower(code) = ?", strings.ToLower(ahsp.Code)).Count(&count)
		if count > 0 {
			return httperr.Conflict(c, fmt.Sprintf("AHSP with code '%s' already exists", ahsp.Code))
		}

		err = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(ahsp).Error; err != nil {
				return err
			}
			for i := range items {
				items[i].AhspID = ahsp.ID
				if err := tx.Create(&items[i]).Error; err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			return httperr.Internal(c, "Failed to create AHSP")
		}

		itemsDTO := make([]AhspItemDTO, len(items))
		for i := range items {
			itemsDTO[i] = toItemDTO(&items[i])
		}

		return c.Status(fiber.StatusCreated).JSON(AhspDetailDTO{
			AhspDTO: toAhspDTO(ahsp),
			Items:   itemsDTO,
		})
	}
}

func UpdateHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid AHSP ID", nil)
		}

		var existing domain.Ahsp
		if err := db.First(&existing, id).Error; err != nil {
			return httperr.NotFound(c, "AHSP not found")
		}

		var input CreateAhspInput
		if err := c.Bind().Body(&input); err != nil {
			return httperr.BadRequest(c, "Invalid request body", nil)
		}

		ahsp, items, err := calculateAndBuild(input)
		if err != nil {
			return httperr.BadRequest(c, err.Error(), nil)
		}

		if !strings.EqualFold(ahsp.Code, existing.Code) {
			var count int64
			db.Model(&domain.Ahsp{}).Where("lower(code) = ? AND id != ?", strings.ToLower(ahsp.Code), id).Count(&count)
			if count > 0 {
				return httperr.Conflict(c, fmt.Sprintf("AHSP with code '%s' already exists", ahsp.Code))
			}
		}

		err = db.Transaction(func(tx *gorm.DB) error {
			ahsp.ID = id
			if err := tx.Model(&existing).Updates(ahsp).Error; err != nil {
				return err
			}
			// Replace items
			if err := tx.Where("ahsp_id = ?", id).Delete(&domain.AhspItem{}).Error; err != nil {
				return err
			}
			for i := range items {
				items[i].AhspID = id
				if err := tx.Create(&items[i]).Error; err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			return httperr.Internal(c, "Failed to update AHSP")
		}

		itemsDTO := make([]AhspItemDTO, len(items))
		for i := range items {
			itemsDTO[i] = toItemDTO(&items[i])
		}

		return c.JSON(AhspDetailDTO{
			AhspDTO: toAhspDTO(ahsp),
			Items:   itemsDTO,
		})
	}
}

func DeleteHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid AHSP ID", nil)
		}

		res := db.Delete(&domain.Ahsp{}, id)
		if res.Error != nil {
			return httperr.Internal(c, "Failed to delete AHSP")
		}
		if res.RowsAffected == 0 {
			return httperr.NotFound(c, "AHSP not found")
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

