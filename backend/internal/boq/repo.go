package boq

import (
	"strings"

	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/model"
)

type ProjectCounts struct {
	BoqCount int64
	VolCount int64
	RabCount int64
}

type Repository interface {
	ListProjects(q string) ([]domain.Project, error)
	GetProjectCounts(projectID uint) (ProjectCounts, error)
	GetBoqItems(projectID uint) ([]domain.BoqItem, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) ListProjects(q string) ([]domain.Project, error) {
	query := r.db.Model(&model.Project{}).Where("deleted_at IS NULL")
	if q != "" {
		like := "%" + strings.ToLower(q) + "%"
		query = query.Where("lower(name) LIKE ? OR lower(section) LIKE ? OR lower(plant) LIKE ?", like, like, like)
	}

	var models []model.Project
	if err := query.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	projects := make([]domain.Project, len(models))
	for i := range models {
		projects[i] = *models[i].ToDomain()
	}
	return projects, nil
}

func (r *repository) GetProjectCounts(projectID uint) (ProjectCounts, error) {
	var counts ProjectCounts

	if err := r.db.Model(&model.BoqItem{}).Where("project_id = ?", projectID).Count(&counts.BoqCount).Error; err != nil {
		return counts, err
	}

	if err := r.db.Table("volumes").
		Joins("JOIN boq_items ON boq_items.id = volumes.boq_item_id").
		Where("boq_items.project_id = ?", projectID).
		Count(&counts.VolCount).Error; err != nil {
		return counts, err
	}

	if err := r.db.Model(&model.RabItem{}).Where("project_id = ?", projectID).Count(&counts.RabCount).Error; err != nil {
		return counts, err
	}

	return counts, nil
}

func (r *repository) GetBoqItems(projectID uint) ([]domain.BoqItem, error) {
	var models []model.BoqItem
	if err := r.db.Where("project_id = ?", projectID).Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]domain.BoqItem, len(models))
	for i := range models {
		items[i] = *models[i].ToDomain()
	}
	return items, nil
}
