package rab

import (
	"strings"

	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/model"
)

type Repository interface {
	ListByProjectID(projectID uint, q string) ([]domain.RabItem, error)
	FindByID(id uint) (*domain.RabItem, error)
	Create(item *domain.RabItem) error
	Update(item *domain.RabItem) error
	Delete(id uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) ListByProjectID(projectID uint, q string) ([]domain.RabItem, error) {
	query := r.db.Model(&model.RabItem{}).Where("project_id = ?", projectID)
	if q != "" {
		like := "%" + strings.ToLower(q) + "%"
		query = query.Where("lower(work_item) LIKE ?", like)
	}

	var models []model.RabItem
	if err := query.Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]domain.RabItem, len(models))
	for i := range models {
		items[i] = *models[i].ToDomain()
	}
	return items, nil
}

func (r *repository) FindByID(id uint) (*domain.RabItem, error) {
	var m model.RabItem
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *repository) Create(item *domain.RabItem) error {
	m := model.FromRabItemDomain(item)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	*item = *m.ToDomain()
	return nil
}

func (r *repository) Update(item *domain.RabItem) error {
	m := model.FromRabItemDomain(item)
	if err := r.db.Save(m).Error; err != nil {
		return err
	}
	*item = *m.ToDomain()
	return nil
}

func (r *repository) Delete(id uint) (bool, error) {
	res := r.db.Delete(&model.RabItem{}, id)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
