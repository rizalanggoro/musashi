package rab

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

var (
	ErrRabItemNotFound  = errors.New("rab item not found")
	ErrRequiredRabField = errors.New("work item and unit are required")
)

type Service interface {
	GetProjectRab(projectID uuid.UUID, q string) (*RabResponseDTO, error)
	Create(projectID uuid.UUID, input CreateRabItemInput) (*RabItemDTO, error)
	Update(id uuid.UUID, input UpdateRabItemInput) (*RabItemDTO, error)
	Delete(id uuid.UUID) error
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

func (s *service) GetProjectRab(projectID uuid.UUID, q string) (*RabResponseDTO, error) {
	items, err := s.repo.ListByProjectID(projectID, strings.TrimSpace(q))
	if err != nil {
		return nil, err
	}

	var totalRab int64
	itemsDTO := make([]RabItemDTO, len(items))
	for i, item := range items {
		totalRab += item.TotalCost
		itemsDTO[i] = toDTO(&item)
	}

	return &RabResponseDTO{
		ProjectID: projectID.String(),
		Items:     itemsDTO,
		TotalRab:  strconv.FormatInt(totalRab, 10),
	}, nil
}

func (s *service) Create(projectID uuid.UUID, input CreateRabItemInput) (*RabItemDTO, error) {
	if strings.TrimSpace(input.WorkItem) == "" || strings.TrimSpace(input.Unit) == "" {
		return nil, ErrRequiredRabField
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

	if err := s.repo.Create(&item); err != nil {
		return nil, err
	}

	dto := toDTO(&item)
	return &dto, nil
}

func (s *service) Update(id uuid.UUID, input UpdateRabItemInput) (*RabItemDTO, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRabItemNotFound
		}
		return nil, err
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

	if err := s.repo.Update(item); err != nil {
		return nil, err
	}

	dto := toDTO(item)
	return &dto, nil
}

func (s *service) Delete(id uuid.UUID) error {
	deleted, err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrRabItemNotFound
	}
	return nil
}
