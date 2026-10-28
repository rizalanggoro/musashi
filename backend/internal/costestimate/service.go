package costestimate

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

var (
	ErrProjectNotFound = errors.New("project not found")
)

type Service interface {
	GetSummary(projectID uuid.UUID) (*CostEstimateDTO, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetSummary(projectID uuid.UUID) (*CostEstimateDTO, error) {
	project, err := s.repo.GetProject(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	totalRab, err := s.repo.GetTotalRab(projectID)
	if err != nil {
		return nil, err
	}

	var spb int64
	if project.ProgressSPB == domain.ProgressDone {
		spb = totalRab
		if spb == 0 {
			spb = project.Budget
		}
	}

	variance := project.Budget - totalRab
	sisa := project.Budget - spb

	return &CostEstimateDTO{
		ProjectID: projectID.String(),
		Budget:    strconv.FormatInt(project.Budget, 10),
		Rab:       strconv.FormatInt(totalRab, 10),
		Spb:       strconv.FormatInt(spb, 10),
		Variance:  strconv.FormatInt(variance, 10),
		Sisa:      strconv.FormatInt(sisa, 10),
	}, nil
}
