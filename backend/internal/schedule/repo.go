package schedule

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type Repository interface {
	FindByProjectID(projectID uuid.UUID) ([]domain.Schedule, error)
	FindByID(id uuid.UUID) (*domain.Schedule, error)
	Create(s *domain.Schedule) error
	Update(s *domain.Schedule) error
	Delete(id uuid.UUID) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByProjectID(projectID uuid.UUID) ([]domain.Schedule, error) {
	var schedules []domain.Schedule
	if err := r.db.Where("project_id = ?", projectID).Order("start_date ASC, code ASC").Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *repository) FindByID(id uuid.UUID) (*domain.Schedule, error) {
	var s domain.Schedule
	if err := r.db.First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) Create(s *domain.Schedule) error {
	return r.db.Create(s).Error
}

func (r *repository) Update(s *domain.Schedule) error {
	return r.db.Save(s).Error
}

func (r *repository) Delete(id uuid.UUID) (bool, error) {
	res := r.db.Delete(&domain.Schedule{}, id)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
