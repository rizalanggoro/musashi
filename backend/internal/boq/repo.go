package boq

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type ProjectCounts struct {
	BoqCount int64
	VolCount int64
	RabCount int64
}

type Repository interface {
	ListProjects(q string) ([]domain.Project, error)
	GetProjectCounts(projectID uuid.UUID) (ProjectCounts, error)
	GetBoqItems(projectID uuid.UUID) ([]domain.BoqItem, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) ListProjects(q string) ([]domain.Project, error) {
	query := r.db.Model(&domain.Project{}).Where("deleted_at IS NULL")
	if q != "" {
		like := "%" + strings.ToLower(q) + "%"
		query = query.Where("lower(name) LIKE ? OR lower(section) LIKE ? OR lower(plant) LIKE ?", like, like, like)
	}

	var projects []domain.Project
	if err := query.Order("created_at DESC").Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *repository) GetProjectCounts(projectID uuid.UUID) (ProjectCounts, error) {
	var counts ProjectCounts

	if err := r.db.Model(&domain.BoqItem{}).Where("project_id = ?", projectID).Count(&counts.BoqCount).Error; err != nil {
		return counts, err
	}

	if err := r.db.Table("volumes").
		Joins("JOIN boq_items ON boq_items.id = volumes.boq_item_id").
		Where("boq_items.project_id = ?", projectID).
		Count(&counts.VolCount).Error; err != nil {
		return counts, err
	}

	if err := r.db.Model(&domain.RabItem{}).Where("project_id = ?", projectID).Count(&counts.RabCount).Error; err != nil {
		return counts, err
	}

	return counts, nil
}

func (r *repository) GetBoqItems(projectID uuid.UUID) ([]domain.BoqItem, error) {
	var items []domain.BoqItem
	if err := r.db.Where("project_id = ?", projectID).Order("created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
