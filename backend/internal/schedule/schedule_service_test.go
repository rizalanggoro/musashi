package schedule

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type mockScheduleRepo struct {
	findByProjectIDFn func(projectID uint) ([]domain.Schedule, error)
	findByIDFn        func(id uint) (*domain.Schedule, error)
	createFn          func(s *domain.Schedule) error
	updateFn          func(s *domain.Schedule) error
	deleteFn          func(id uint) (bool, error)
}

func (m *mockScheduleRepo) FindByProjectID(projectID uint) ([]domain.Schedule, error) {
	if m.findByProjectIDFn != nil {
		return m.findByProjectIDFn(projectID)
	}
	return nil, nil
}

func (m *mockScheduleRepo) FindByID(id uint) (*domain.Schedule, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockScheduleRepo) Create(s *domain.Schedule) error {
	if m.createFn != nil {
		return m.createFn(s)
	}
	s.ID = 1
	return nil
}

func (m *mockScheduleRepo) Update(s *domain.Schedule) error {
	if m.updateFn != nil {
		return m.updateFn(s)
	}
	return nil
}

func (m *mockScheduleRepo) Delete(id uint) (bool, error) {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return true, nil
}

func TestScheduleService(t *testing.T) {
	start := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 6, 5, 0, 0, 0, 0, time.UTC)

	t.Run("ListByProject", func(t *testing.T) {
		repo := &mockScheduleRepo{
			findByProjectIDFn: func(projectID uint) ([]domain.Schedule, error) {
				return []domain.Schedule{
					{ID: 1, ProjectID: projectID, Code: "SCH-01", WorkItem: "Pekerjaan 1", StartDate: start, EndDate: end, Duration: 5},
				}, nil
			},
		}
		svc := NewService(repo)
		list, err := svc.ListByProject(1)
		if err != nil || len(list) != 1 || list[0].Code != "SCH-01" {
			t.Fatalf("unexpected list: %v, list: %+v", err, list)
		}
	})

	t.Run("Create validation & success", func(t *testing.T) {
		svc := NewService(&mockScheduleRepo{})

		// Missing required field
		_, err := svc.Create(1, CreateScheduleInput{Code: "", WorkItem: "Item"})
		if !errors.Is(err, ErrRequiredFields) {
			t.Fatalf("expected ErrRequiredFields, got %v", err)
		}

		// Invalid date format
		_, err = svc.Create(1, CreateScheduleInput{
			Code:      "SCH-01",
			WorkItem:  "Galian",
			StartDate: "invalid-date",
			EndDate:   "2026-06-05",
		})
		if !errors.Is(err, ErrInvalidDate) {
			t.Fatalf("expected ErrInvalidDate, got %v", err)
		}

		// Success
		created, err := svc.Create(1, CreateScheduleInput{
			Code:      "SCH-01",
			WorkItem:  "Galian",
			StartDate: "2026-06-01",
			EndDate:   "2026-06-05",
			Weight:    "10.5",
			Progress:  "20",
		})
		if err != nil || created.Duration != 5 {
			t.Fatalf("unexpected create result: %v, created: %+v", err, created)
		}
	})

	t.Run("Update not found and success", func(t *testing.T) {
		repo := &mockScheduleRepo{
			findByIDFn: func(id uint) (*domain.Schedule, error) {
				if id == 1 {
					return &domain.Schedule{ID: 1, ProjectID: 10, Code: "OLD", WorkItem: "Old Item", StartDate: start, EndDate: end, Duration: 5}, nil
				}
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := NewService(repo)

		// Not found
		_, err := svc.Update(99, UpdateScheduleInput{})
		if !errors.Is(err, ErrScheduleNotFound) {
			t.Fatalf("expected ErrScheduleNotFound, got %v", err)
		}

		// Success
		newCode := "NEW-01"
		newWork := "New Item"
		updated, err := svc.Update(1, UpdateScheduleInput{
			Code:     &newCode,
			WorkItem: &newWork,
		})
		if err != nil || updated.Code != "NEW-01" || updated.WorkItem != "New Item" {
			t.Fatalf("unexpected update result: %v, updated: %+v", err, updated)
		}
	})

	t.Run("Delete not found and success", func(t *testing.T) {
		repo := &mockScheduleRepo{
			deleteFn: func(id uint) (bool, error) {
				return id == 1, nil
			},
		}
		svc := NewService(repo)

		if err := svc.Delete(1); err != nil {
			t.Fatalf("expected success, got %v", err)
		}
		if err := svc.Delete(99); !errors.Is(err, ErrScheduleNotFound) {
			t.Fatalf("expected ErrScheduleNotFound, got %v", err)
		}
	})
}

func TestScheduleHandler(t *testing.T) {
	start := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 6, 5, 0, 0, 0, 0, time.UTC)
	repo := &mockScheduleRepo{
		findByProjectIDFn: func(projectID uint) ([]domain.Schedule, error) {
			return []domain.Schedule{
				{ID: 1, ProjectID: projectID, Code: "SCH-01", WorkItem: "Pondasi", StartDate: start, EndDate: end, Duration: 5},
			}, nil
		},
	}
	svc := NewService(repo)
	handler := NewHandler(svc)

	app := fiber.New()
	app.Get("/projects/:id/schedules", handler.ListByProject)
	app.Post("/projects/:id/schedules", handler.Create)
	app.Put("/schedules/:id", handler.Update)
	app.Delete("/schedules/:id", handler.Delete)

	t.Run("HTTP ListByProject", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/projects/1/schedules", nil)
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("HTTP Create success", func(t *testing.T) {
		body := `{"code":"SCH-02","work_item":"Atap","start_date":"2026-06-01","end_date":"2026-06-10"}`
		req := httptest.NewRequest(http.MethodPost, "/projects/1/schedules", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201, got %d, err: %v", resp.StatusCode, err)
		}
	})
}

