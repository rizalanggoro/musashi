package unitprice

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

var (
	ErrUnitPriceNotFound    = errors.New("unit price not found")
	ErrRequiredFields       = errors.New("code, work item, and unit are required")
	ErrCodeAlreadyExists    = errors.New("unit price with that code already exists")
	ErrReferencedByAhspItem = errors.New("cannot delete unit price: referenced by one or more AHSP items")
)

type Service interface {
	List(q, category string) ([]UnitPriceDTO, error)
	Create(input CreateUnitPriceInput) (*UnitPriceDTO, error)
	Update(id uint, input UpdateUnitPriceInput) (*UnitPriceDTO, error)
	Delete(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) List(q, category string) ([]UnitPriceDTO, error) {
	items, err := s.repo.List(strings.TrimSpace(q), strings.TrimSpace(category))
	if err != nil {
		return nil, err
	}

	data := make([]UnitPriceDTO, len(items))
	for i := range items {
		data[i] = toDTO(&items[i])
	}
	return data, nil
}

func (s *service) Create(input CreateUnitPriceInput) (*UnitPriceDTO, error) {
	code := strings.TrimSpace(input.Code)
	if code == "" || input.WorkItem == "" || input.Unit == "" {
		return nil, ErrRequiredFields
	}

	existing, err := s.repo.FindByCode(code, nil)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("%w: '%s'", ErrCodeAlreadyExists, code)
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

	if err := s.repo.Create(&item); err != nil {
		return nil, err
	}

	dto := toDTO(&item)
	return &dto, nil
}

func (s *service) Update(id uint, input UpdateUnitPriceInput) (*UnitPriceDTO, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUnitPriceNotFound
		}
		return nil, err
	}

	if input.Code != nil {
		newCode := strings.TrimSpace(*input.Code)
		if newCode != "" && !strings.EqualFold(newCode, item.Code) {
			existing, err := s.repo.FindByCode(newCode, &id)
			if err == nil && existing != nil {
				return nil, fmt.Errorf("%w: '%s'", ErrCodeAlreadyExists, newCode)
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

	if err := s.repo.Update(item); err != nil {
		return nil, err
	}

	dto := toDTO(item)
	return &dto, nil
}

func (s *service) Delete(id uint) error {
	count, err := s.repo.CountAhspUsage(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrReferencedByAhspItem
	}

	deleted, err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrUnitPriceNotFound
	}
	return nil
}
