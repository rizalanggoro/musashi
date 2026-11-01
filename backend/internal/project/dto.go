package project

import (
	"strconv"
	"time"

	"musashi/backend/internal/domain"
	"musashi/backend/utils"
)

type PaginationMeta struct {
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	Total   int64 `json:"total"`
}

type ProjectListResponse struct {
	Data []Project      `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type Project struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	Section            string  `json:"section"`
	Plant              string  `json:"plant"`
	ImplementationPlan *string `json:"implementation_plan"`
	Budget             string  `json:"budget"`
	Note               string  `json:"note"`
	Status               string  `json:"status"`
	ProgressPPT          string  `json:"progress_ppt"`
	ProgressDrawing      string  `json:"progress_drawing"`
	ProgressRAB          string  `json:"progress_rab"`
	ProgressVendor       string  `json:"progress_vendor"`
	ProgressSPB          string  `json:"progress_spb"`
	ProgressPR           string  `json:"progress_pr"`
	ProgressRINGI        string  `json:"progress_ringi"`
	ProgressPO           string  `json:"progress_po"`
	ProgressConstruction string  `json:"progress_construction"`
	ProgressBAST         string  `json:"progress_bast"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

type ProjectDTO = Project

type CreateProjectInput struct {
	Name                 string  `json:"name"`
	Section              string  `json:"section"`
	Plant                string  `json:"plant"`
	ImplementationPlan   *string `json:"implementation_plan,omitempty"`
	Budget               string  `json:"budget"`
	Note                 string  `json:"note"`
	ProgressPPT          string  `json:"progress_ppt,omitempty"`
	ProgressDrawing      string  `json:"progress_drawing,omitempty"`
	ProgressRAB          string  `json:"progress_rab,omitempty"`
	ProgressVendor       string  `json:"progress_vendor,omitempty"`
	ProgressSPB          string  `json:"progress_spb,omitempty"`
	ProgressPR           string  `json:"progress_pr,omitempty"`
	ProgressRINGI        string  `json:"progress_ringi,omitempty"`
	ProgressPO           string  `json:"progress_po,omitempty"`
	ProgressConstruction string  `json:"progress_construction,omitempty"`
	ProgressBAST         string  `json:"progress_bast,omitempty"`
}

type UpdateProjectInput struct {
	Name                 *string `json:"name,omitempty"`
	Section              *string `json:"section,omitempty"`
	Plant                *string `json:"plant,omitempty"`
	ImplementationPlan   *string `json:"implementation_plan,omitempty"`
	Budget               *string `json:"budget,omitempty"`
	Note                 *string `json:"note,omitempty"`
	Status               *string `json:"status,omitempty"`
	ProgressPPT          *string `json:"progress_ppt,omitempty"`
	ProgressDrawing      *string `json:"progress_drawing,omitempty"`
	ProgressRAB          *string `json:"progress_rab,omitempty"`
	ProgressVendor       *string `json:"progress_vendor,omitempty"`
	ProgressSPB          *string `json:"progress_spb,omitempty"`
	ProgressPR           *string `json:"progress_pr,omitempty"`
	ProgressRINGI        *string `json:"progress_ringi,omitempty"`
	ProgressPO           *string `json:"progress_po,omitempty"`
	ProgressConstruction *string `json:"progress_construction,omitempty"`
	ProgressBAST         *string `json:"progress_bast,omitempty"`
}

func toDTO(p *domain.Project) ProjectDTO {
	dto := ProjectDTO{
		ID:                   utils.FormatID(p.ID),
		Name:                 p.Name,
		Section:              p.Section,
		Plant:                p.Plant,
		Budget:               strconv.FormatInt(p.Budget, 10),
		Note:                 p.Note,
		Status:               string(p.Status),
		ProgressPPT:          string(p.ProgressPPT),
		ProgressDrawing:      string(p.ProgressDrawing),
		ProgressRAB:          string(p.ProgressRAB),
		ProgressVendor:       string(p.ProgressVendor),
		ProgressSPB:          string(p.ProgressSPB),
		ProgressPR:           string(p.ProgressPR),
		ProgressRINGI:        string(p.ProgressRINGI),
		ProgressPO:           string(p.ProgressPO),
		ProgressConstruction: string(p.ProgressConstruction),
		ProgressBAST:         string(p.ProgressBAST),
		CreatedAt:            p.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:            p.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if p.ImplementationPlan != nil {
		d := p.ImplementationPlan.Format("2006-01-02")
		dto.ImplementationPlan = &d
	}
	return dto
}
