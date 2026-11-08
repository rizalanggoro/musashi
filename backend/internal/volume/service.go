package volume

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/utils"
)

var (
	ErrVolumeNotFound   = errors.New("volume row not found")
	ErrBoqItemNotFound  = errors.New("boq item not found")
	ErrRequiredBoqField = errors.New("work item and unit are required")
)

type Service interface {
	List(projectID uint, q string) ([]VolumeRowDTO, error)
	Create(projectID uint, input CreateVolumeInput) (*VolumeRowDTO, error)
	Update(id uint, input UpdateVolumeInput) (*VolumeRowDTO, error)
	Delete(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func calcQuantities(formula, unit string, length, width, height *float64, plannedVolume float64, unitWeight float64, customQty *float64) (float64, float64) {
	if unitWeight <= 0 {
		unitWeight = 1
	}

	if customQty != nil && *customQty > 0 {
		return *customQty, (*customQty) * unitWeight
	}

	u := strings.ToLower(strings.TrimSpace(unit))
	f := strings.ToLower(strings.TrimSpace(formula))

	var l, w, h float64
	if length != nil {
		l = *length
	}
	if width != nil {
		w = *width
	}
	if height != nil {
		h = *height
	}

	var qty float64

	switch f {
	case "cylinder": // 1/4 * pi * D^2 * T (lubang bor / bore pile / tiang silinder)
		depth := h
		if depth <= 0 && w > 0 {
			depth = w
		}
		if depth <= 0 {
			depth = 1
		}
		qty = 0.25 * math.Pi * l * l * depth

	case "area_circle": // 1/4 * pi * D^2 (lingkaran)
		qty = 0.25 * math.Pi * l * l

	case "area_pl": // P x L
		qty = l * w

	case "area_pt": // P x T
		qty = l * h

	case "area_lt": // L x T
		qty = w * h

	case "linear_p": // P (meter lari)
		qty = l

	case "linear_t": // T (meter lari tinggi)
		qty = h

	case "weight_rebar": // P * berat_jenis (kg)
		bj := w
		if bj <= 0 {
			bj = 1
		}
		qty = l * bj

	case "unit_only": // lot, pcs, unit
		qty = 1.0

	case "box": // P x L x T
		qty = l * w * h

	default:
		// Smart auto-detection fallback if formula not explicitly set
		switch u {
		case "lot", "ls", "lumpsum", "unit", "pcs", "titik", "buah", "set", "btg", "bh":
			if plannedVolume > 0 {
				qty = plannedVolume
			} else {
				qty = 1.0
			}
		case "m'", "m", "meter", "ml":
			if l > 0 {
				qty = l
			} else if h > 0 {
				qty = h
			} else if w > 0 {
				qty = w
			} else {
				qty = plannedVolume
			}
		case "m2", "m²":
			if l > 0 && w > 0 {
				qty = l * w
			} else if l > 0 && h > 0 {
				qty = l * h
			} else if w > 0 && h > 0 {
				qty = w * h
			} else if l > 0 {
				qty = 0.25 * math.Pi * l * l
			} else {
				qty = plannedVolume
			}
		case "kg", "ton":
			if l > 0 && w > 0 {
				qty = l * w
			} else if l > 0 {
				qty = l
			} else {
				qty = plannedVolume
			}
		case "m3", "m³":
			if l > 0 && w > 0 && h > 0 {
				qty = l * w * h
			} else if l > 0 && h > 0 { // silinder D x T
				qty = 0.25 * math.Pi * l * l * h
			} else {
				qty = plannedVolume
			}
		default:
			if l > 0 && w > 0 && h > 0 {
				qty = l * w * h
			} else if l > 0 && w > 0 {
				qty = l * w
			} else if l > 0 && h > 0 {
				qty = l * h
			} else if l > 0 {
				qty = l
			} else {
				qty = plannedVolume
			}
		}
	}

	if qty <= 0 && plannedVolume > 0 {
		qty = plannedVolume
	}

	totalQty := qty * unitWeight
	return qty, totalQty
}

func (s *service) List(projectID uint, q string) ([]VolumeRowDTO, error) {
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

func (s *service) Create(projectID uint, input CreateVolumeInput) (*VolumeRowDTO, error) {
	var boqItem domain.BoqItem
	if input.BoqItemID != nil && *input.BoqItemID != "" {
		bID, err := utils.ParseID(*input.BoqItemID)
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

		var ahspID *uint
		if input.AhspID != nil && *input.AhspID != "" {
			if parsed, err := utils.ParseID(*input.AhspID); err == nil {
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

	formulaStr := ""
	if input.Formula != nil {
		formulaStr = *input.Formula
	}

	var customQty *float64
	if input.Quantity != nil && *input.Quantity != "" {
		if q, err := strconv.ParseFloat(*input.Quantity, 64); err == nil && q > 0 {
			customQty = &q
		}
	}

	qty, totalQty := calcQuantities(formulaStr, boqItem.Unit, lVal, wVal, hVal, boqItem.PlannedVolume, uWeight, customQty)

	var dateVal *time.Time
	if input.Date != nil && *input.Date != "" {
		if t, err := time.Parse("2006-01-02", *input.Date); err == nil {
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
		Formula:       formulaStr,
		Description:   input.Description,
	}

	if err := s.repo.CreateVolume(&vol); err != nil {
		return nil, err
	}

	dto := toRowDTO(&vol, &boqItem)
	return &dto, nil
}

func (s *service) Update(id uint, input UpdateVolumeInput) (*VolumeRowDTO, error) {
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
	if input.Formula != nil {
		vol.Formula = *input.Formula
	}
	if input.Description != nil {
		vol.Description = *input.Description
	}

	var customQty *float64
	if input.Quantity != nil && *input.Quantity != "" {
		if q, err := strconv.ParseFloat(*input.Quantity, 64); err == nil && q > 0 {
			customQty = &q
		}
	}

	vol.Quantity, vol.TotalQuantity = calcQuantities(vol.Formula, boqItem.Unit, vol.Length, vol.Width, vol.Height, boqItem.PlannedVolume, vol.UnitWeight, customQty)

	if err := s.repo.UpdateVolume(vol); err != nil {
		return nil, err
	}

	dto := toRowDTO(vol, boqItem)
	return &dto, nil
}

func (s *service) Delete(id uint) error {
	deleted, err := s.repo.DeleteVolume(id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrVolumeNotFound
	}
	return nil
}
