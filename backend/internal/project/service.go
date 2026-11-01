package project

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

var (
	ErrProjectNotFound      = errors.New("project not found")
	ErrProjectNameRequired  = errors.New("project name is required")
	ErrProjectAlreadyExists = errors.New("project with that name already exists")
	ErrInvalidSection       = errors.New("section must be one of: GA Building, EHS, Wellfare")
	ErrInvalidPlant         = errors.New("plant must be one of: Cikarang, Karawang")
)

type Service interface {
	List(q string, page, perPage int) ([]ProjectDTO, int64, error)
	Get(id uint) (*ProjectDTO, error)
	Create(input CreateProjectInput) (*ProjectDTO, error)
	Update(id uint, input UpdateProjectInput) (*ProjectDTO, error)
	Delete(id uint) error
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

func (s *service) Get(id uint) (*ProjectDTO, error) {
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
	if input.ImplementationPlan != nil && *input.ImplementationPlan != "" {
		t, err := time.Parse("2006-01-02", *input.ImplementationPlan)
		if err == nil {
			planDate = &t
		}
	}

	sec := strings.TrimSpace(input.Section)
	if sec == "" {
		sec = string(domain.ProjectSectionGABuilding)
	}
	if !domain.IsValidSection(sec) {
		return nil, ErrInvalidSection
	}

	pl := strings.TrimSpace(input.Plant)
	if pl == "" {
		pl = string(domain.ProjectPlantCikarang)
	}
	if !domain.IsValidPlant(pl) {
		return nil, ErrInvalidPlant
	}

	p := domain.Project{
		Name:                 name,
		Section:              sec,
		Plant:                pl,
		ImplementationPlan:   planDate,
		Budget:               budgetVal,
		Note:                 input.Note,
		Status:               domain.ProjectStatusPlanning,
		ProgressPPT:          parseProgress(input.ProgressPPT),
		ProgressDrawing:      parseProgress(input.ProgressDrawing),
		ProgressRAB:          parseProgress(input.ProgressRAB),
		ProgressVendor:       parseProgress(input.ProgressVendor),
		ProgressSPB:          parseProgress(input.ProgressSPB),
		ProgressPR:           parseProgress(input.ProgressPR),
		ProgressRINGI:        parseProgress(input.ProgressRINGI),
		ProgressPO:           parseProgress(input.ProgressPO),
		ProgressConstruction: parseProgress(input.ProgressConstruction),
		ProgressBAST:         parseProgress(input.ProgressBAST),
	}

	if err := s.repo.Create(&p); err != nil {
		return nil, err
	}

	dto := toDTO(&p)
	return &dto, nil
}

func (s *service) Update(id uint, input UpdateProjectInput) (*ProjectDTO, error) {
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
		sec := strings.TrimSpace(*input.Section)
		if !domain.IsValidSection(sec) {
			return nil, ErrInvalidSection
		}
		p.Section = sec
	}
	if input.Plant != nil {
		pl := strings.TrimSpace(*input.Plant)
		if !domain.IsValidPlant(pl) {
			return nil, ErrInvalidPlant
		}
		p.Plant = pl
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
	if input.Note != nil {
		p.Note = *input.Note
	}
	if input.Status != nil {
		p.Status = domain.ProjectStatus(*input.Status)
	}
	if input.ProgressPPT != nil {
		p.ProgressPPT = domain.ProgressStatus(*input.ProgressPPT)
	}
	if input.ProgressDrawing != nil {
		p.ProgressDrawing = domain.ProgressStatus(*input.ProgressDrawing)
	}
	if input.ProgressRAB != nil {
		p.ProgressRAB = domain.ProgressStatus(*input.ProgressRAB)
	}
	if input.ProgressVendor != nil {
		p.ProgressVendor = domain.ProgressStatus(*input.ProgressVendor)
	}
	if input.ProgressSPB != nil {
		p.ProgressSPB = domain.ProgressStatus(*input.ProgressSPB)
	}
	if input.ProgressPR != nil {
		p.ProgressPR = domain.ProgressStatus(*input.ProgressPR)
	}
	if input.ProgressRINGI != nil {
		p.ProgressRINGI = domain.ProgressStatus(*input.ProgressRINGI)
	}
	if input.ProgressPO != nil {
		p.ProgressPO = domain.ProgressStatus(*input.ProgressPO)
	}
	if input.ProgressConstruction != nil {
		p.ProgressConstruction = domain.ProgressStatus(*input.ProgressConstruction)
	}
	if input.ProgressBAST != nil {
		p.ProgressBAST = domain.ProgressStatus(*input.ProgressBAST)
	}

	if err := s.repo.Update(p); err != nil {
		return nil, err
	}

	dto := toDTO(p)
	return &dto, nil
}

func (s *service) Delete(id uint) error {
	deleted, err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrProjectNotFound
	}
	return nil
}

func parseProgress(val string) domain.ProgressStatus {
	if strings.EqualFold(val, "done") {
		return domain.ProgressDone
	}
	return domain.ProgressPending
}
