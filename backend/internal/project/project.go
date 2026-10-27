package project

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/platform/httperr"
)

type ProjectDTO struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Section            string `json:"section"`
	Plant              string `json:"plant"`
	ImplementationPlan *string `json:"implementation_plan"`
	Budget             string `json:"budget"`
	Information        string `json:"information"`
	Status             string `json:"status"`
	ProgressSPB        string `json:"progress_spb"`
	ProgressRAB        string `json:"progress_rab"`
	ProgressXXX        string `json:"progress_xxx"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
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
		ID:          p.ID.String(),
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

func ListHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		q := strings.TrimSpace(c.Query("q"))
		page, _ := strconv.Atoi(c.Query("page", "1"))
		if page < 1 {
			page = 1
		}
		perPage, _ := strconv.Atoi(c.Query("per_page", "25"))
		if perPage < 1 || perPage > 100 {
			perPage = 25
		}

		query := db.Model(&domain.Project{}).Where("deleted_at IS NULL")
		if q != "" {
			like := "%" + strings.ToLower(q) + "%"
			query = query.Where("lower(name) LIKE ? OR lower(section) LIKE ? OR lower(plant) LIKE ?", like, like, like)
		}

		var total int64
		if err := query.Count(&total).Error; err != nil {
			return httperr.Internal(c, "Failed to count projects")
		}

		var projects []domain.Project
		offset := (page - 1) * perPage
		if err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&projects).Error; err != nil {
			return httperr.Internal(c, "Failed to fetch projects")
		}

		data := make([]ProjectDTO, len(projects))
		for i := range projects {
			data[i] = toDTO(&projects[i])
		}

		return c.JSON(fiber.Map{
			"data": data,
			"meta": fiber.Map{
				"page":     page,
				"per_page": perPage,
				"total":    total,
			},
		})
	}
}

func CreateHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var input CreateProjectInput
		if err := c.Bind().Body(&input); err != nil {
			return httperr.BadRequest(c, "Invalid request body", nil)
		}
		name := strings.TrimSpace(input.Name)
		if name == "" {
			return httperr.BadRequest(c, "Project name is required", map[string]string{"name": "Project name is required"})
		}

		// Check uniqueness
		var count int64
		db.Model(&domain.Project{}).Where("lower(name) = ? AND deleted_at IS NULL", strings.ToLower(name)).Count(&count)
		if count > 0 {
			return httperr.Conflict(c, fmt.Sprintf("Project with name '%s' already exists", name))
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

		if err := db.Create(&p).Error; err != nil {
			return httperr.Internal(c, "Failed to create project")
		}

		return c.Status(fiber.StatusCreated).JSON(toDTO(&p))
	}
}

func GetHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid project ID", nil)
		}

		var p domain.Project
		if err := db.Where("id = ? AND deleted_at IS NULL", id).First(&p).Error; err != nil {
			return httperr.NotFound(c, "Project not found")
		}

		return c.JSON(toDTO(&p))
	}
}

func UpdateHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid project ID", nil)
		}

		var p domain.Project
		if err := db.Where("id = ? AND deleted_at IS NULL", id).First(&p).Error; err != nil {
			return httperr.NotFound(c, "Project not found")
		}

		var input UpdateProjectInput
		if err := c.Bind().Body(&input); err != nil {
			return httperr.BadRequest(c, "Invalid request body", nil)
		}

		if input.Name != nil {
			newName := strings.TrimSpace(*input.Name)
			if newName != "" && !strings.EqualFold(newName, p.Name) {
				var count int64
				db.Model(&domain.Project{}).Where("lower(name) = ? AND id != ? AND deleted_at IS NULL", strings.ToLower(newName), id).Count(&count)
				if count > 0 {
					return httperr.Conflict(c, fmt.Sprintf("Project with name '%s' already exists", newName))
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

		if err := db.Save(&p).Error; err != nil {
			return httperr.Internal(c, "Failed to update project")
		}

		return c.JSON(toDTO(&p))
	}
}

func DeleteHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid project ID", nil)
		}

		res := db.Where("id = ? AND deleted_at IS NULL", id).Delete(&domain.Project{})
		if res.Error != nil {
			return httperr.Internal(c, "Failed to delete project")
		}
		if res.RowsAffected == 0 {
			return httperr.NotFound(c, "Project not found")
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

