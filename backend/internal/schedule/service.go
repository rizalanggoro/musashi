package schedule

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

var (
	ErrScheduleNotFound = errors.New("schedule not found")
	ErrRequiredFields   = errors.New("code, work item, start date, and end date are required")
	ErrInvalidDate      = errors.New("start date and end date must be YYYY-MM-DD")
)

type Service interface {
	ListByProject(projectID uint) ([]ScheduleDTO, error)
	Create(projectID uint, input CreateScheduleInput) (*ScheduleDTO, error)
	Update(id uint, input UpdateScheduleInput) (*ScheduleDTO, error)
	Delete(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func calcDuration(start, end time.Time) int {
	days := int(end.Sub(start).Hours()/24) + 1
	if days < 1 {
		return 1
	}
	return days
}

func (s *service) ListByProject(projectID uint) ([]ScheduleDTO, error) {
	schedules, err := s.repo.FindByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	data := make([]ScheduleDTO, len(schedules))
	for i := range schedules {
		data[i] = toDTO(&schedules[i])
	}
	return data, nil
}

func (s *service) Create(projectID uint, input CreateScheduleInput) (*ScheduleDTO, error) {
	if strings.TrimSpace(input.Code) == "" || strings.TrimSpace(input.WorkItem) == "" || input.StartDate == "" || input.EndDate == "" {
		return nil, ErrRequiredFields
	}

	start, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return nil, ErrInvalidDate
	}
	end, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return nil, ErrInvalidDate
	}

	duration := calcDuration(start, end)
	weight, _ := strconv.ParseFloat(input.Weight, 64)
	progress, _ := strconv.ParseFloat(input.Progress, 64)

	status := domain.SchedulePlanning
	if input.Status != "" {
		status = domain.ScheduleStatus(input.Status)
	}

	sched := domain.Schedule{
		ProjectID: projectID,
		Code:      strings.TrimSpace(input.Code),
		WorkItem:  strings.TrimSpace(input.WorkItem),
		StartDate: start,
		EndDate:   end,
		Duration:  duration,
		Weight:    weight,
		Progress:  progress,
		Status:    status,
	}

	if err := s.repo.Create(&sched); err != nil {
		return nil, err
	}

	dto := toDTO(&sched)
	return &dto, nil
}

func (s *service) Update(id uint, input UpdateScheduleInput) (*ScheduleDTO, error) {
	sched, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScheduleNotFound
		}
		return nil, err
	}

	if input.Code != nil {
		sched.Code = strings.TrimSpace(*input.Code)
	}
	if input.WorkItem != nil {
		sched.WorkItem = strings.TrimSpace(*input.WorkItem)
	}
	if input.StartDate != nil {
		t, err := time.Parse("2006-01-02", *input.StartDate)
		if err == nil {
			sched.StartDate = t
		}
	}
	if input.EndDate != nil {
		t, err := time.Parse("2006-01-02", *input.EndDate)
		if err == nil {
			sched.EndDate = t
		}
	}
	sched.Duration = calcDuration(sched.StartDate, sched.EndDate)

	if input.Weight != nil {
		w, _ := strconv.ParseFloat(*input.Weight, 64)
		sched.Weight = w
	}
	if input.Progress != nil {
		p, _ := strconv.ParseFloat(*input.Progress, 64)
		sched.Progress = p
	}
	if input.Status != nil {
		sched.Status = domain.ScheduleStatus(*input.Status)
	}

	if err := s.repo.Update(sched); err != nil {
		return nil, err
	}

	dto := toDTO(sched)
	return &dto, nil
}

func (s *service) Delete(id uint) error {
	deleted, err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrScheduleNotFound
	}
	return nil
}
