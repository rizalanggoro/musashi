package costestimate

import (
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/model"
)

type Repository interface {
	GetProject(projectID uint) (*domain.Project, error)
	GetTotalRab(projectID uint) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetProject(projectID uint) (*domain.Project, error) {
	var m model.Project
	if err := r.db.Where("id = ? AND deleted_at IS NULL", projectID).First(&m).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *repository) GetTotalRab(projectID uint) (int64, error) {
	var totalRab int64
	err := r.db.Model(&model.RabItem{}).Where("project_id = ?", projectID).Select("COALESCE(SUM(total_cost), 0)").Scan(&totalRab).Error
	return totalRab, err
}
