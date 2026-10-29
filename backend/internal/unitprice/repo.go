package unitprice

import (
	"strings"

	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/model"
)

type Repository interface {
	List(q, category string) ([]domain.UnitPrice, error)
	FindByID(id uint) (*domain.UnitPrice, error)
	FindByCode(code string, excludeID *uint) (*domain.UnitPrice, error)
	CountAhspUsage(id uint) (int64, error)
	Create(u *domain.UnitPrice) error
	Update(u *domain.UnitPrice) error
	Delete(id uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(q, category string) ([]domain.UnitPrice, error) {
	query := r.db.Model(&model.UnitPrice{})
	if q != "" {
		like := "%" + strings.ToLower(q) + "%"
		query = query.Where("lower(code) LIKE ? OR lower(work_item) LIKE ?", like, like)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var models []model.UnitPrice
	if err := query.Order("category ASC, code ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]domain.UnitPrice, len(models))
	for i := range models {
		items[i] = *models[i].ToDomain()
	}
	return items, nil
}

func (r *repository) FindByID(id uint) (*domain.UnitPrice, error) {
	var m model.UnitPrice
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *repository) FindByCode(code string, excludeID *uint) (*domain.UnitPrice, error) {
	var m model.UnitPrice
	query := r.db.Where("lower(code) = ?", strings.ToLower(code))
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	if err := query.First(&m).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *repository) CountAhspUsage(id uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.AhspItem{}).Where("unit_price_id = ?", id).Count(&count).Error
	return count, err
}

func (r *repository) Create(u *domain.UnitPrice) error {
	m := model.FromUnitPriceDomain(u)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	*u = *m.ToDomain()
	return nil
}

func (r *repository) Update(u *domain.UnitPrice) error {
	m := model.FromUnitPriceDomain(u)
	if err := r.db.Save(m).Error; err != nil {
		return err
	}
	*u = *m.ToDomain()
	return nil
}

func (r *repository) Delete(id uint) (bool, error) {
	res := r.db.Delete(&model.UnitPrice{}, id)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
