package unitprice

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type Repository interface {
	List(q, category string) ([]domain.UnitPrice, error)
	FindByID(id uuid.UUID) (*domain.UnitPrice, error)
	FindByCode(code string, excludeID *uuid.UUID) (*domain.UnitPrice, error)
	CountAhspUsage(id uuid.UUID) (int64, error)
	Create(u *domain.UnitPrice) error
	Update(u *domain.UnitPrice) error
	Delete(id uuid.UUID) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(q, category string) ([]domain.UnitPrice, error) {
	query := r.db.Model(&domain.UnitPrice{})
	if q != "" {
		like := "%" + strings.ToLower(q) + "%"
		query = query.Where("lower(code) LIKE ? OR lower(work_item) LIKE ?", like, like)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var items []domain.UnitPrice
	if err := query.Order("category ASC, code ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) FindByID(id uuid.UUID) (*domain.UnitPrice, error) {
	var item domain.UnitPrice
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) FindByCode(code string, excludeID *uuid.UUID) (*domain.UnitPrice, error) {
	var item domain.UnitPrice
	query := r.db.Where("lower(code) = ?", strings.ToLower(code))
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	if err := query.First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) CountAhspUsage(id uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&domain.AhspItem{}).Where("unit_price_id = ?", id).Count(&count).Error
	return count, err
}

func (r *repository) Create(u *domain.UnitPrice) error {
	return r.db.Create(u).Error
}

func (r *repository) Update(u *domain.UnitPrice) error {
	return r.db.Save(u).Error
}

func (r *repository) Delete(id uuid.UUID) (bool, error) {
	res := r.db.Delete(&domain.UnitPrice{}, id)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
