package schedule

import (
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/model"
)

type Repository interface {
	FindByProjectID(projectID uint) ([]domain.Schedule, error)
	FindByID(id uint) (*domain.Schedule, error)
	Create(s *domain.Schedule) error
	Update(s *domain.Schedule) error
	Delete(id uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByProjectID(projectID uint) ([]domain.Schedule, error) {
	var models []model.Schedule
	if err := r.db.Where("project_id = ?", projectID).Order("start_date ASC, code ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	schedules := make([]domain.Schedule, len(models))
	for i := range models {
		schedules[i] = *models[i].ToDomain()
	}
	return schedules, nil
}

func (r *repository) FindByID(id uint) (*domain.Schedule, error) {
	var m model.Schedule
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *repository) Create(s *domain.Schedule) error {
	m := model.FromScheduleDomain(s)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	*s = *m.ToDomain()
	return nil
}

func (r *repository) Update(s *domain.Schedule) error {
	m := model.FromScheduleDomain(s)
	if err := r.db.Save(m).Error; err != nil {
		return err
	}
	*s = *m.ToDomain()
	return nil
}

func (r *repository) Delete(id uint) (bool, error) {
	res := r.db.Delete(&model.Schedule{}, id)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
