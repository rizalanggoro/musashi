package volume

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

var (
	ErrVolumeNotFound   = errors.New("volume row not found")
	ErrBoqItemNotFound  = errors.New("boq item not found")
	ErrRequiredBoqField = errors.New("work item and unit are required")
)

type Service interface {
	List(projectID uuid.UUID, q string) ([]VolumeRowDTO, error)
	Create(projectID uuid.UUID, input CreateVolumeInput) (*VolumeRowDTO, error)
	Update(id uuid.UUID, input UpdateVolumeInput) (*VolumeRowDTO, error)
	Delete(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
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

func (s *service) List(projectID uuid.UUID, q string) ([]VolumeRowDTO, error) {
	rows, err := s.repo.ListRows(projectID, strings.TrimSpace(q))
	if err != nil {
		return nil, err
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
	return data, nil
}

func (s *service) Create(projectID uuid.UUID, input CreateVolumeInput) (*VolumeRowDTO, error) {
	var boqItem domain.BoqItem
	if input.BoqItemID != nil && *input.BoqItemID != "" {
		bID, err := uuid.Parse(*input.BoqItemID)
		if err != nil {
			return nil, ErrBoqItemNotFound
		}
		item, err := s.repo.FindBoqItemInProject(bID, projectID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrBoqItemNotFound
			}
			return nil, err
		}
		boqItem = *item
	} else {
		if strings.TrimSpace(input.WorkItem) == "" || strings.TrimSpace(input.Unit) == "" {
			return nil, ErrRequiredBoqField
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
		if err := s.repo.CreateBoqItem(&boqItem); err != nil {
			return nil, err
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

	if err := s.repo.CreateVolume(&vol); err != nil {
		return nil, err
	}

	dto := toRowDTO(&vol, &boqItem)
	return &dto, nil
}

func (s *service) Update(id uuid.UUID, input UpdateVolumeInput) (*VolumeRowDTO, error) {
	vol, err := s.repo.FindVolumeByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVolumeNotFound
		}
		return nil, err
	}

	boqItem, err := s.repo.FindBoqItemByID(vol.BoqItemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBoqItemNotFound
		}
		return nil, err
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

	if err := s.repo.UpdateVolume(vol); err != nil {
		return nil, err
	}

	dto := toRowDTO(vol, boqItem)
	return &dto, nil
}

func (s *service) Delete(id uuid.UUID) error {
	deleted, err := s.repo.DeleteVolume(id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrVolumeNotFound
	}
	return nil
}
