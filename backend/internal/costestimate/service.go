package costestimate

import (
	"errors"
	"strconv"

	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/utils"
)

var (
	ErrProjectNotFound = errors.New("project not found")
)

type Service interface {
	GetSummary(projectID uint) (*CostEstimateDTO, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetSummary(projectID uint) (*CostEstimateDTO, error) {
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
		ProjectID: utils.FormatID(projectID),
		Budget:    strconv.FormatInt(project.Budget, 10),
		Rab:       strconv.FormatInt(totalRab, 10),
		Spb:       strconv.FormatInt(spb, 10),
		Variance:  strconv.FormatInt(variance, 10),
		Sisa:      strconv.FormatInt(sisa, 10),
	}, nil
}
