package export

import (
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/model"
)

type Repository interface {
	GetProject(projectID uint) (*domain.Project, error)
	GetVolumeExportRows(projectID uint) ([]VolumeExportRow, error)
	GetRabItems(projectID uint) ([]domain.RabItem, error)
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

func (r *repository) GetVolumeExportRows(projectID uint) ([]VolumeExportRow, error) {
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

func (r *repository) GetRabItems(projectID uint) ([]domain.RabItem, error) {
	var models []model.RabItem
	if err := r.db.Where("project_id = ?", projectID).Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]domain.RabItem, len(models))
	for i := range models {
		items[i] = *models[i].ToDomain()
	}
	return items, nil
}
