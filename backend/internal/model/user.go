package model

import (
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type User struct {
	gorm.Model
	Email        string          `gorm:"uniqueIndex;not null"`
	Name         string          `gorm:"not null"`
	PasswordHash string          `gorm:"not null"`
	Role         domain.UserRole `gorm:"type:text;default:'staff'"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) ToDomain() *domain.User {
	if u == nil {
		return nil
	}
	return &domain.User{
		ID:           u.ID,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		Email:        u.Email,
		Name:         u.Name,
		PasswordHash: u.PasswordHash,
		Role:         u.Role,
	}
}

func FromUserDomain(d *domain.User) *User {
	if d == nil {
		return nil
	}
	u := &User{
		Email:        d.Email,
		Name:         d.Name,
		PasswordHash: d.PasswordHash,
		Role:         d.Role,
	}
	u.ID = d.ID
	u.CreatedAt = d.CreatedAt
	u.UpdatedAt = d.UpdatedAt
	return u
}

