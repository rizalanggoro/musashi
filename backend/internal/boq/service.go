package boq

import (
	"strings"

	"github.com/google/uuid"

	"musashi/backend/internal/domain"
)

type Service interface {
	ListSummary(q string) ([]BoqProjectSummaryDTO, error)
	GetProjectBoq(projectID uuid.UUID) ([]BoqItemDTO, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func calculateStatuses(p domain.Project, counts ProjectCounts) (string, string) {
	volumeStatus := "empty"
	if counts.BoqCount > 0 {
		if counts.VolCount == 0 {
			volumeStatus = "empty"
		} else if counts.VolCount < counts.BoqCount {
			volumeStatus = "partial"
		} else {
			volumeStatus = "complete"
		}
	}

	rabStatus := "empty"
	if counts.RabCount > 0 {
		if counts.BoqCount > 0 && counts.RabCount < counts.BoqCount {
			rabStatus = "partial"
		} else {
			rabStatus = "complete"
		}
	}

	return volumeStatus, rabStatus
}

func (s *service) ListSummary(q string) ([]BoqProjectSummaryDTO, error) {
	projects, err := s.repo.ListProjects(strings.TrimSpace(q))
	if err != nil {
		return nil, err
	}

	results := make([]BoqProjectSummaryDTO, len(projects))
	for i, p := range projects {
		counts, err := s.repo.GetProjectCounts(p.ID)
		if err != nil {
			return nil, err
		}

		volStatus, rabStatus := calculateStatuses(p, counts)

		results[i] = BoqProjectSummaryDTO{
			ID:           p.ID.String(),
			Name:         p.Name,
			Section:      p.Section,
			Plant:        p.Plant,
			VolumeStatus: volStatus,
			RabStatus:    rabStatus,
		}
	}

	return results, nil
}

func (s *service) GetProjectBoq(projectID uuid.UUID) ([]BoqItemDTO, error) {
	items, err := s.repo.GetBoqItems(projectID)
	if err != nil {
		return nil, err
	}

	dto := make([]BoqItemDTO, len(items))
	for i, item := range items {
		dto[i] = toItemDTO(&item)
	}
	return dto, nil
}
