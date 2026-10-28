package volume

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type Repository interface {
	ListRows(projectID uuid.UUID, q string) ([]VolumeRowQueryResult, error)
	FindVolumeByID(id uuid.UUID) (*domain.Volume, error)
	FindBoqItemByID(id uuid.UUID) (*domain.BoqItem, error)
	FindBoqItemInProject(id uuid.UUID, projectID uuid.UUID) (*domain.BoqItem, error)
	CreateBoqItem(b *domain.BoqItem) error
	CreateVolume(v *domain.Volume) error
	UpdateVolume(v *domain.Volume) error
	DeleteVolume(id uuid.UUID) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) ListRows(projectID uuid.UUID, q string) ([]VolumeRowQueryResult, error) {
	var rows []VolumeRowQueryResult

	query := r.db.Table("volumes").
		Select("volumes.*, boq_items.work_item, boq_items.unit, boq_items.planned_volume").
		Joins("JOIN boq_items ON boq_items.id = volumes.boq_item_id").
		Where("boq_items.project_id = ?", projectID)

	if q != "" {
		like := "%" + strings.ToLower(q) + "%"
		query = query.Where("lower(boq_items.work_item) LIKE ? OR lower(volumes.description) LIKE ?", like, like)
	}

	if err := query.Order("boq_items.created_at ASC, volumes.created_at ASC").Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *repository) FindVolumeByID(id uuid.UUID) (*domain.Volume, error) {
	var vol domain.Volume
	if err := r.db.First(&vol, id).Error; err != nil {
		return nil, err
	}
	return &vol, nil
}

func (r *repository) FindBoqItemByID(id uuid.UUID) (*domain.BoqItem, error) {
	var boqItem domain.BoqItem
	if err := r.db.First(&boqItem, id).Error; err != nil {
		return nil, err
	}
	return &boqItem, nil
}

func (r *repository) FindBoqItemInProject(id uuid.UUID, projectID uuid.UUID) (*domain.BoqItem, error) {
	var boqItem domain.BoqItem
	if err := r.db.Where("id = ? AND project_id = ?", id, projectID).First(&boqItem).Error; err != nil {
		return nil, err
	}
	return &boqItem, nil
}

func (r *repository) CreateBoqItem(b *domain.BoqItem) error {
	return r.db.Create(b).Error
}

func (r *repository) CreateVolume(v *domain.Volume) error {
	return r.db.Create(v).Error
}

func (r *repository) UpdateVolume(v *domain.Volume) error {
	return r.db.Save(v).Error
}

func (r *repository) DeleteVolume(id uuid.UUID) (bool, error) {
	res := r.db.Delete(&domain.Volume{}, id)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
