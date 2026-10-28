package project

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type Repository interface {
	List(q string, page, perPage int) ([]domain.Project, int64, error)
	FindByID(id uuid.UUID) (*domain.Project, error)
	FindByName(name string, excludeID *uuid.UUID) (*domain.Project, error)
	Create(p *domain.Project) error
	Update(p *domain.Project) error
	Delete(id uuid.UUID) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(q string, page, perPage int) ([]domain.Project, int64, error) {
	query := r.db.Model(&domain.Project{}).Where("deleted_at IS NULL")
	if q != "" {
		like := "%" + strings.ToLower(q) + "%"
		query = query.Where("lower(name) LIKE ? OR lower(section) LIKE ? OR lower(plant) LIKE ?", like, like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var projects []domain.Project
	offset := (page - 1) * perPage
	if err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *repository) FindByID(id uuid.UUID) (*domain.Project, error) {
	var p domain.Project
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) FindByName(name string, excludeID *uuid.UUID) (*domain.Project, error) {
	var p domain.Project
	query := r.db.Where("lower(name) = ? AND deleted_at IS NULL", strings.ToLower(name))
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	if err := query.First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) Create(p *domain.Project) error {
	return r.db.Create(p).Error
}

func (r *repository) Update(p *domain.Project) error {
	return r.db.Save(p).Error
}

func (r *repository) Delete(id uuid.UUID) (bool, error) {
	res := r.db.Where("id = ? AND deleted_at IS NULL", id).Delete(&domain.Project{})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
