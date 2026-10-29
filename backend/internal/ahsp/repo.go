package ahsp

import (
	"strings"

	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/model"
)

type Repository interface {
	List(q string) ([]domain.Ahsp, error)
	FindByID(id uint) (*domain.Ahsp, error)
	FindByCode(code string, excludeID *uint) (*domain.Ahsp, error)
	Create(ahsp *domain.Ahsp, items []domain.AhspItem) error
	Update(id uint, ahsp *domain.Ahsp, items []domain.AhspItem) error
	Delete(id uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(q string) ([]domain.Ahsp, error) {
	query := r.db.Model(&model.Ahsp{})
	if q != "" {
		like := "%" + strings.ToLower(q) + "%"
		query = query.Where("lower(code) LIKE ? OR lower(work_item) LIKE ?", like, like)
	}

	var models []model.Ahsp
	if err := query.Order("code ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	ahsps := make([]domain.Ahsp, len(models))
	for i := range models {
		ahsps[i] = *models[i].ToDomain()
	}
	return ahsps, nil
}

func (r *repository) FindByID(id uint) (*domain.Ahsp, error) {
	var a model.Ahsp
	if err := r.db.Preload("Items", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("sort ASC, created_at ASC")
	}).First(&a, id).Error; err != nil {
		return nil, err
	}
	return a.ToDomain(), nil
}

func (r *repository) FindByCode(code string, excludeID *uint) (*domain.Ahsp, error) {
	var a model.Ahsp
	query := r.db.Where("lower(code) = ?", strings.ToLower(code))
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	if err := query.First(&a).Error; err != nil {
		return nil, err
	}
	return a.ToDomain(), nil
}

func (r *repository) Create(ahsp *domain.Ahsp, items []domain.AhspItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		ahspModel := model.FromAhspDomain(ahsp)
		if err := tx.Create(ahspModel).Error; err != nil {
			return err
		}
		ahsp.ID = ahspModel.ID
		ahsp.CreatedAt = ahspModel.CreatedAt
		ahsp.UpdatedAt = ahspModel.UpdatedAt

		for i := range items {
			items[i].AhspID = ahspModel.ID
			itemModel := model.FromAhspItemDomain(&items[i])
			if err := tx.Create(itemModel).Error; err != nil {
				return err
			}
			items[i].ID = itemModel.ID
			items[i].CreatedAt = itemModel.CreatedAt
			items[i].UpdatedAt = itemModel.UpdatedAt
		}
		return nil
	})
}

func (r *repository) Update(id uint, ahsp *domain.Ahsp, items []domain.AhspItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing model.Ahsp
		if err := tx.First(&existing, id).Error; err != nil {
			return err
		}

		ahsp.ID = id
		ahspModel := model.FromAhspDomain(ahsp)
		if err := tx.Model(&existing).Updates(ahspModel).Error; err != nil {
			return err
		}

		// Replace items
		if err := tx.Where("ahsp_id = ?", id).Delete(&model.AhspItem{}).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].AhspID = id
			itemModel := model.FromAhspItemDomain(&items[i])
			if err := tx.Create(itemModel).Error; err != nil {
				return err
			}
			items[i].ID = itemModel.ID
			items[i].CreatedAt = itemModel.CreatedAt
			items[i].UpdatedAt = itemModel.UpdatedAt
		}
		return nil
	})
}

func (r *repository) Delete(id uint) (bool, error) {
	res := r.db.Delete(&model.Ahsp{}, id)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
