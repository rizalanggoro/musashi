package auth

import (
	"strings"

	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/model"
)

type Repository interface {
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByEmail(email string) (*domain.User, error) {
	var u model.User
	if err := r.db.Where("lower(email) = ?", strings.ToLower(email)).First(&u).Error; err != nil {
		return nil, err
	}
	return u.ToDomain(), nil
}

func (r *repository) FindByID(id uint) (*domain.User, error) {
	var u model.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return u.ToDomain(), nil
}
