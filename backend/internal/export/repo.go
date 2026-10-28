package export

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type Repository interface {
	GetProject(projectID uuid.UUID) (*domain.Project, error)
	GetVolumeExportRows(projectID uuid.UUID) ([]VolumeExportRow, error)
	GetRabItems(projectID uuid.UUID) ([]domain.RabItem, error)
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

func (r *repository) GetVolumeExportRows(projectID uuid.UUID) ([]VolumeExportRow, error) {
	var rows []VolumeExportRow
	query := r.db.Table("volumes").
		Select("volumes.*, boq_items.work_item, boq_items.unit, boq_items.planned_volume").
		Joins("JOIN boq_items ON boq_items.id = volumes.boq_item_id").
		Where("boq_items.project_id = ?", projectID).
		Order("boq_items.created_at ASC, volumes.created_at ASC")

	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *repository) GetRabItems(projectID uuid.UUID) ([]domain.RabItem, error) {
	var items []domain.RabItem
	if err := r.db.Where("project_id = ?", projectID).Order("created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
