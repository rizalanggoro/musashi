package rab

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type mockRabRepo struct {
	listByProjectIDFn func(projectID uint, q string) ([]domain.RabItem, error)
	findByIDFn        func(id uint) (*domain.RabItem, error)
	createFn          func(item *domain.RabItem) error
	updateFn          func(item *domain.RabItem) error
	deleteFn          func(id uint) (bool, error)
}

func (m *mockRabRepo) ListByProjectID(projectID uint, q string) ([]domain.RabItem, error) {
	if m.listByProjectIDFn != nil {
		return m.listByProjectIDFn(projectID, q)
	}
	return nil, nil
}

func (m *mockRabRepo) FindByID(id uint) (*domain.RabItem, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRabRepo) Create(item *domain.RabItem) error {
	if m.createFn != nil {
		return m.createFn(item)
	}
	item.ID = 1
	return nil
}

func (m *mockRabRepo) Update(item *domain.RabItem) error {
	if m.updateFn != nil {
		return m.updateFn(item)
	}
	return nil
}

func (m *mockRabRepo) Delete(id uint) (bool, error) {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return true, nil
}

func TestRoundHalfUp(t *testing.T) {
	if got := roundHalfUp(10.5); got != 11 {
		t.Errorf("roundHalfUp(10.5) = %d, want 11", got)
	}
	if got := roundHalfUp(10.4); got != 10 {
		t.Errorf("roundHalfUp(10.4) = %d, want 10", got)
	}
	if got := roundHalfUp(-5.5); got != -6 {
		t.Errorf("roundHalfUp(-5.5) = %d, want -6", got)
	}
}

func TestRabService(t *testing.T) {
	t.Run("GetProjectRab", func(t *testing.T) {
		repo := &mockRabRepo{
			listByProjectIDFn: func(projectID uint, q string) ([]domain.RabItem, error) {
				return []domain.RabItem{
					{ID: 1, ProjectID: projectID, WorkItem: "Item 1", Unit: "m3", Volume: 10, UnitPrice: 50000, TotalCost: 500000},
					{ID: 2, ProjectID: projectID, WorkItem: "Item 2", Unit: "m2", Volume: 20, UnitPrice: 10000, TotalCost: 200000},
				}, nil
			},
		}
		svc := NewService(repo)
		res, err := svc.GetProjectRab(1, "")
		if err != nil || res.TotalRab != "700000" || len(res.Items) != 2 {
			t.Fatalf("unexpected rab result: %v, total: %s", err, res.TotalRab)
		}
	})

	t.Run("Create validation & calculation", func(t *testing.T) {
		svc := NewService(&mockRabRepo{})

		// Missing required field
		_, err := svc.Create(1, CreateRabItemInput{WorkItem: "", Unit: "m3"})
		if !errors.Is(err, ErrRequiredRabField) {
			t.Fatalf("expected ErrRequiredRabField, got %v", err)
		}

		// Success
		created, err := svc.Create(1, CreateRabItemInput{
			WorkItem:  "Galian",
			Unit:      "m3",
			Volume:    "12.5",
			UnitPrice: "100000",
		})
		if err != nil || created.TotalCost != "1250000" {
			t.Fatalf("unexpected create: %v, created: %+v", err, created)
		}
	})

	t.Run("Update not found and recalculation", func(t *testing.T) {
		repo := &mockRabRepo{
			findByIDFn: func(id uint) (*domain.RabItem, error) {
				if id == 1 {
					return &domain.RabItem{ID: 1, ProjectID: 10, WorkItem: "Galian", Unit: "m3", Volume: 10, UnitPrice: 100000, TotalCost: 1000000}, nil
				}
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := NewService(repo)

		// Not found
		_, err := svc.Update(99, UpdateRabItemInput{})
		if !errors.Is(err, ErrRabItemNotFound) {
			t.Fatalf("expected ErrRabItemNotFound, got %v", err)
		}

		// Success with recalculation
		newVol := "20"
		newPrice := "150000"
		updated, err := svc.Update(1, UpdateRabItemInput{
			Volume:    &newVol,
			UnitPrice: &newPrice,
		})
		if err != nil || updated.TotalCost != "3000000" {
			t.Fatalf("unexpected update: %v, updated: %+v", err, updated)
		}
	})

	t.Run("Delete not found and success", func(t *testing.T) {
		repo := &mockRabRepo{
			deleteFn: func(id uint) (bool, error) {
				return id == 1, nil
			},
		}
		svc := NewService(repo)

		if err := svc.Delete(1); err != nil {
			t.Fatalf("expected success, got %v", err)
		}
		if err := svc.Delete(99); !errors.Is(err, ErrRabItemNotFound) {
			t.Fatalf("expected ErrRabItemNotFound, got %v", err)
		}
	})
}

func TestRabHandler(t *testing.T) {
	repo := &mockRabRepo{
		listByProjectIDFn: func(projectID uint, q string) ([]domain.RabItem, error) {
			return []domain.RabItem{{ID: 1, ProjectID: projectID, WorkItem: "Pondasi", Unit: "m3", Volume: 5, UnitPrice: 100000, TotalCost: 500000}}, nil
		},
	}
	svc := NewService(repo)
	handler := NewHandler(svc)

	app := fiber.New()
	app.Get("/projects/:id/rab", handler.GetProjectRab)
	app.Post("/projects/:id/rab", handler.Create)
	app.Put("/rab/:id", handler.Update)
	app.Delete("/rab/:id", handler.Delete)

	t.Run("HTTP GetProjectRab", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/projects/1/rab", nil)
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("HTTP Create success", func(t *testing.T) {
		body := `{"work_item":"Dinding","unit":"m2","volume":"20","unit_price":"50000"}`
		req := httptest.NewRequest(http.MethodPost, "/projects/1/rab", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201, got %d, err: %v", resp.StatusCode, err)
		}
	})
}

