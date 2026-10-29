package ahsp

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/utils"
)

var (
	ErrAhspNotFound      = errors.New("ahsp not found")
	ErrCodeAlreadyExists = errors.New("ahsp with that code already exists")
	ErrRequiredFields    = errors.New("code, work item, and unit are required")
)

type Service interface {
	List(q string) ([]AhspDTO, error)
	Get(id uint) (*AhspDetailDTO, error)
	Create(input CreateAhspInput) (*AhspDetailDTO, error)
	Update(id uint, input CreateAhspInput) (*AhspDetailDTO, error)
	Delete(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func roundHalfUp(val float64) int64 {
	if val < 0 {
		return int64(math.Ceil(val - 0.5))
	}
	return int64(math.Floor(val + 0.5))
}

func calculateAndBuild(input CreateAhspInput) (*domain.Ahsp, []domain.AhspItem, error) {
	code := strings.TrimSpace(input.Code)
	if code == "" || input.WorkItem == "" || input.Unit == "" {
		return nil, nil, ErrRequiredFields
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

		var upID *uint
		if inItem.UnitPriceID != nil && *inItem.UnitPriceID != "" {
			if parsed, err := utils.ParseID(*inItem.UnitPriceID); err == nil {
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

func (s *service) List(q string) ([]AhspDTO, error) {
	ahsps, err := s.repo.List(q)
	if err != nil {
		return nil, err
	}

	data := make([]AhspDTO, len(ahsps))
	for i := range ahsps {
		data[i] = toAhspDTO(&ahsps[i])
	}
	return data, nil
}

func (s *service) Get(id uint) (*AhspDetailDTO, error) {
	a, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAhspNotFound
		}
		return nil, err
	}

	itemsDTO := make([]AhspItemDTO, len(a.Items))
	for i := range a.Items {
		itemsDTO[i] = toItemDTO(&a.Items[i])
	}

	return &AhspDetailDTO{
		AhspDTO: toAhspDTO(a),
		Items:   itemsDTO,
	}, nil
}

func (s *service) Create(input CreateAhspInput) (*AhspDetailDTO, error) {
	ahsp, items, err := calculateAndBuild(input)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByCode(ahsp.Code, nil)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("%w: '%s'", ErrCodeAlreadyExists, ahsp.Code)
	}

	if err := s.repo.Create(ahsp, items); err != nil {
		return nil, err
	}

	itemsDTO := make([]AhspItemDTO, len(items))
	for i := range items {
		itemsDTO[i] = toItemDTO(&items[i])
	}

	return &AhspDetailDTO{
		AhspDTO: toAhspDTO(ahsp),
		Items:   itemsDTO,
	}, nil
}

func (s *service) Update(id uint, input CreateAhspInput) (*AhspDetailDTO, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAhspNotFound
		}
		return nil, err
	}

	ahsp, items, err := calculateAndBuild(input)
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(ahsp.Code, existing.Code) {
		dup, err := s.repo.FindByCode(ahsp.Code, &id)
		if err == nil && dup != nil {
			return nil, fmt.Errorf("%w: '%s'", ErrCodeAlreadyExists, ahsp.Code)
		}
	}

	if err := s.repo.Update(id, ahsp, items); err != nil {
		return nil, err
	}

	itemsDTO := make([]AhspItemDTO, len(items))
	for i := range items {
		itemsDTO[i] = toItemDTO(&items[i])
	}

	return &AhspDetailDTO{
		AhspDTO: toAhspDTO(ahsp),
		Items:   itemsDTO,
	}, nil
}

func (s *service) Delete(id uint) error {
	deleted, err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrAhspNotFound
	}
	return nil
}
