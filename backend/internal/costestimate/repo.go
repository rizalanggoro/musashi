package costestimate

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type Repository interface {
	GetProject(projectID uuid.UUID) (*domain.Project, error)
	GetTotalRab(projectID uuid.UUID) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetProject(projectID uuid.UUID) (*domain.Project, error) {
	var project domain.Project
	if err := r.db.Where("id = ? AND deleted_at IS NULL", projectID).First(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *repository) GetTotalRab(projectID uuid.UUID) (int64, error) {
	var totalRab int64
	err := r.db.Model(&domain.RabItem{}).Where("project_id = ?", projectID).Select("COALESCE(SUM(total_cost), 0)").Scan(&totalRab).Error
	return totalRab, err
}
