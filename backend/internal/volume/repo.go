package volume

import (
	"strings"

	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/model"
)

type Repository interface {
	ListRows(projectID uint, q string) ([]VolumeRowQueryResult, error)
	FindVolumeByID(id uint) (*domain.Volume, error)
	FindBoqItemByID(id uint) (*domain.BoqItem, error)
	FindBoqItemInProject(id uint, projectID uint) (*domain.BoqItem, error)
	CreateBoqItem(b *domain.BoqItem) error
	CreateVolume(v *domain.Volume) error
	UpdateVolume(v *domain.Volume) error
	DeleteVolume(id uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) ListRows(projectID uint, q string) ([]VolumeRowQueryResult, error) {
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

func (r *repository) FindVolumeByID(id uint) (*domain.Volume, error) {
	var v model.Volume
	if err := r.db.First(&v, id).Error; err != nil {
		return nil, err
	}
	return v.ToDomain(), nil
}

func (r *repository) FindBoqItemByID(id uint) (*domain.BoqItem, error) {
	var b model.BoqItem
	if err := r.db.First(&b, id).Error; err != nil {
		return nil, err
	}
	return b.ToDomain(), nil
}

func (r *repository) FindBoqItemInProject(id uint, projectID uint) (*domain.BoqItem, error) {
	var b model.BoqItem
	if err := r.db.Where("id = ? AND project_id = ?", id, projectID).First(&b).Error; err != nil {
		return nil, err
	}
	return b.ToDomain(), nil
}

func (r *repository) CreateBoqItem(b *domain.BoqItem) error {
	m := model.FromBoqItemDomain(b)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	*b = *m.ToDomain()
	return nil
}

func (r *repository) CreateVolume(v *domain.Volume) error {
	m := model.FromVolumeDomain(v)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	*v = *m.ToDomain()
	return nil
}

func (r *repository) UpdateVolume(v *domain.Volume) error {
	m := model.FromVolumeDomain(v)
	if err := r.db.Save(m).Error; err != nil {
		return err
	}
	*v = *m.ToDomain()
	return nil
}

func (r *repository) DeleteVolume(id uint) (bool, error) {
	res := r.db.Delete(&model.Volume{}, id)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
