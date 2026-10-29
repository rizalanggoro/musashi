package project

import (
	"strconv"
	"time"

	"musashi/backend/internal/domain"
	"musashi/backend/utils"
)

type ProjectDTO struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	Section            string  `json:"section"`
	Plant              string  `json:"plant"`
	ImplementationPlan *string `json:"implementation_plan"`
	Budget             string  `json:"budget"`
	Information        string  `json:"information"`
	Status             string  `json:"status"`
	ProgressSPB        string  `json:"progress_spb"`
	ProgressRAB        string  `json:"progress_rab"`
	ProgressXXX        string  `json:"progress_xxx"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

type CreateProjectInput struct {
	Name               string `json:"name"`
	Section            string `json:"section"`
	Plant              string `json:"plant"`
	ImplementationPlan string `json:"implementation_plan"`
	Budget             string `json:"budget"`
	Information        string `json:"information"`
	ProgressSPB        string `json:"progress_spb"`
	ProgressRAB        string `json:"progress_rab"`
	ProgressXXX        string `json:"progress_xxx"`
}

type UpdateProjectInput struct {
	Name               *string `json:"name"`
	Section            *string `json:"section"`
	Plant              *string `json:"plant"`
	ImplementationPlan *string `json:"implementation_plan"`
	Budget             *string `json:"budget"`
	Information        *string `json:"information"`
	Status             *string `json:"status"`
	ProgressSPB        *string `json:"progress_spb"`
	ProgressRAB        *string `json:"progress_rab"`
	ProgressXXX        *string `json:"progress_xxx"`
}

func toDTO(p *domain.Project) ProjectDTO {
	dto := ProjectDTO{
		ID:          utils.FormatID(p.ID),
		Name:        p.Name,
		Section:     p.Section,
		Plant:       p.Plant,
		Budget:      strconv.FormatInt(p.Budget, 10),
		Information: p.Information,
		Status:      string(p.Status),
		ProgressSPB: string(p.ProgressSPB),
		ProgressRAB: string(p.ProgressRAB),
		ProgressXXX: string(p.ProgressXXX),
		CreatedAt:   p.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if p.ImplementationPlan != nil {
		d := p.ImplementationPlan.Format("2006-01-02")
		dto.ImplementationPlan = &d
	}
	return dto
}
