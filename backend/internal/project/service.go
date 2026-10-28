package project

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

var (
	ErrProjectNotFound      = errors.New("project not found")
	ErrProjectNameRequired  = errors.New("project name is required")
	ErrProjectAlreadyExists = errors.New("project with that name already exists")
)

type Service interface {
	List(q string, page, perPage int) ([]ProjectDTO, int64, error)
	Get(id uuid.UUID) (*ProjectDTO, error)
	Create(input CreateProjectInput) (*ProjectDTO, error)
	Update(id uuid.UUID, input UpdateProjectInput) (*ProjectDTO, error)
	Delete(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) List(q string, page, perPage int) ([]ProjectDTO, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 25
	}

	projects, total, err := s.repo.List(q, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	data := make([]ProjectDTO, len(projects))
	for i := range projects {
		data[i] = toDTO(&projects[i])
	}
	return data, total, nil
}

func (s *service) Get(id uuid.UUID) (*ProjectDTO, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	dto := toDTO(p)
	return &dto, nil
}

func (s *service) Create(input CreateProjectInput) (*ProjectDTO, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrProjectNameRequired
	}

	existing, err := s.repo.FindByName(name, nil)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("%w: '%s'", ErrProjectAlreadyExists, name)
	}

	var budgetVal int64
	if input.Budget != "" {
		budgetVal, _ = strconv.ParseInt(input.Budget, 10, 64)
	}

	var planDate *time.Time
	if input.ImplementationPlan != "" {
		t, err := time.Parse("2006-01-02", input.ImplementationPlan)
		if err == nil {
			planDate = &t
		}
	}

	spb := domain.ProgressPending
	if input.ProgressSPB == "done" {
		spb = domain.ProgressDone
	}
	rab := domain.ProgressPending
	if input.ProgressRAB == "done" {
		rab = domain.ProgressDone
	}
	xxx := domain.ProgressPending
	if input.ProgressXXX == "done" {
		xxx = domain.ProgressDone
	}

	p := domain.Project{
		Name:               name,
		Section:            input.Section,
		Plant:              input.Plant,
		ImplementationPlan: planDate,
		Budget:             budgetVal,
		Information:        input.Information,
		Status:             domain.ProjectStatusPlanning,
		ProgressSPB:        spb,
		ProgressRAB:        rab,
		ProgressXXX:        xxx,
	}

	if err := s.repo.Create(&p); err != nil {
		return nil, err
	}

	dto := toDTO(&p)
	return &dto, nil
}

func (s *service) Update(id uuid.UUID, input UpdateProjectInput) (*ProjectDTO, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	if input.Name != nil {
		newName := strings.TrimSpace(*input.Name)
		if newName != "" && !strings.EqualFold(newName, p.Name) {
			existing, err := s.repo.FindByName(newName, &id)
			if err == nil && existing != nil {
				return nil, fmt.Errorf("%w: '%s'", ErrProjectAlreadyExists, newName)
			}
			p.Name = newName
		}
	}
	if input.Section != nil {
		p.Section = *input.Section
	}
	if input.Plant != nil {
		p.Plant = *input.Plant
	}
	if input.ImplementationPlan != nil {
		if *input.ImplementationPlan != "" {
			t, err := time.Parse("2006-01-02", *input.ImplementationPlan)
			if err == nil {
				p.ImplementationPlan = &t
			}
		} else {
			p.ImplementationPlan = nil
		}
	}
	if input.Budget != nil {
		b, _ := strconv.ParseInt(*input.Budget, 10, 64)
		p.Budget = b
	}
	if input.Information != nil {
		p.Information = *input.Information
	}
	if input.Status != nil {
		p.Status = domain.ProjectStatus(*input.Status)
	}
	if input.ProgressSPB != nil {
		p.ProgressSPB = domain.ProgressStatus(*input.ProgressSPB)
	}
	if input.ProgressRAB != nil {
		p.ProgressRAB = domain.ProgressStatus(*input.ProgressRAB)
	}
	if input.ProgressXXX != nil {
		p.ProgressXXX = domain.ProgressStatus(*input.ProgressXXX)
	}

	if err := s.repo.Update(p); err != nil {
		return nil, err
	}

	dto := toDTO(p)
	return &dto, nil
}

func (s *service) Delete(id uuid.UUID) error {
	deleted, err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrProjectNotFound
	}
	return nil
}
