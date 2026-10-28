package rab

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type Repository interface {
	ListByProjectID(projectID uuid.UUID, q string) ([]domain.RabItem, error)
	FindByID(id uuid.UUID) (*domain.RabItem, error)
	Create(item *domain.RabItem) error
	Update(item *domain.RabItem) error
	Delete(id uuid.UUID) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) ListByProjectID(projectID uuid.UUID, q string) ([]domain.RabItem, error) {
	query := r.db.Where("project_id = ?", projectID)
	if q != "" {
		like := "%" + strings.ToLower(q) + "%"
		query = query.Where("lower(work_item) LIKE ?", like)
	}

	var items []domain.RabItem
	if err := query.Order("created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) FindByID(id uuid.UUID) (*domain.RabItem, error) {
	var item domain.RabItem
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) Create(item *domain.RabItem) error {
	return r.db.Create(item).Error
}

func (r *repository) Update(item *domain.RabItem) error {
	return r.db.Save(item).Error
}

func (r *repository) Delete(id uuid.UUID) (bool, error) {
	res := r.db.Delete(&domain.RabItem{}, id)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
