package project

import (
	"strings"

	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/model"
)

type Repository interface {
	List(q string, page, perPage int) ([]domain.Project, int64, error)
	FindByID(id uint) (*domain.Project, error)
	FindByName(name string, excludeID *uint) (*domain.Project, error)
	Create(p *domain.Project) error
	Update(p *domain.Project) error
	Delete(id uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(q string, page, perPage int) ([]domain.Project, int64, error) {
	query := r.db.Model(&model.Project{}).Where("deleted_at IS NULL")
	if q != "" {
		like := "%" + strings.ToLower(q) + "%"
		query = query.Where("lower(name) LIKE ? OR lower(section) LIKE ? OR lower(plant) LIKE ?", like, like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []model.Project
	offset := (page - 1) * perPage
	if err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	projects := make([]domain.Project, len(models))
	for i := range models {
		projects[i] = *models[i].ToDomain()
	}

	return projects, total, nil
}

func (r *repository) FindByID(id uint) (*domain.Project, error) {
	var m model.Project
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&m).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *repository) FindByName(name string, excludeID *uint) (*domain.Project, error) {
	var m model.Project
	query := r.db.Where("lower(name) = ? AND deleted_at IS NULL", strings.ToLower(name))
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	if err := query.First(&m).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *repository) Create(p *domain.Project) error {
	m := model.FromProjectDomain(p)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	*p = *m.ToDomain()
	return nil
}

func (r *repository) Update(p *domain.Project) error {
	m := model.FromProjectDomain(p)
	if err := r.db.Save(m).Error; err != nil {
		return err
	}
	*p = *m.ToDomain()
	return nil
}

func (r *repository) Delete(id uint) (bool, error) {
	res := r.db.Where("id = ? AND deleted_at IS NULL", id).Delete(&model.Project{})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
