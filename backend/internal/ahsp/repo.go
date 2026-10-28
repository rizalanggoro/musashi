package ahsp

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type Repository interface {
	List(q string) ([]domain.Ahsp, error)
	FindByID(id uuid.UUID) (*domain.Ahsp, error)
	FindByCode(code string, excludeID *uuid.UUID) (*domain.Ahsp, error)
	Create(ahsp *domain.Ahsp, items []domain.AhspItem) error
	Update(id uuid.UUID, ahsp *domain.Ahsp, items []domain.AhspItem) error
	Delete(id uuid.UUID) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(q string) ([]domain.Ahsp, error) {
	query := r.db.Model(&domain.Ahsp{})
	if q != "" {
		like := "%" + strings.ToLower(q) + "%"
		query = query.Where("lower(code) LIKE ? OR lower(work_item) LIKE ?", like, like)
	}

	var ahsps []domain.Ahsp
	if err := query.Order("code ASC").Find(&ahsps).Error; err != nil {
		return nil, err
	}
	return ahsps, nil
}

func (r *repository) FindByID(id uuid.UUID) (*domain.Ahsp, error) {
	var a domain.Ahsp
	if err := r.db.Preload("Items", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("sort ASC, created_at ASC")
	}).First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repository) FindByCode(code string, excludeID *uuid.UUID) (*domain.Ahsp, error) {
	var a domain.Ahsp
	query := r.db.Where("lower(code) = ?", strings.ToLower(code))
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	if err := query.First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repository) Create(ahsp *domain.Ahsp, items []domain.AhspItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(ahsp).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].AhspID = ahsp.ID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *repository) Update(id uuid.UUID, ahsp *domain.Ahsp, items []domain.AhspItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing domain.Ahsp
		if err := tx.First(&existing, id).Error; err != nil {
			return err
		}

		ahsp.ID = id
		if err := tx.Model(&existing).Updates(ahsp).Error; err != nil {
			return err
		}

		// Replace items
		if err := tx.Where("ahsp_id = ?", id).Delete(&domain.AhspItem{}).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].AhspID = id
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *repository) Delete(id uuid.UUID) (bool, error) {
	res := r.db.Delete(&domain.Ahsp{}, id)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
