package unitprice

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

type mockUnitPriceRepo struct {
	listFn           func(q, category string) ([]domain.UnitPrice, error)
	findByIDFn       func(id uint) (*domain.UnitPrice, error)
	findByCodeFn     func(code string, excludeID *uint) (*domain.UnitPrice, error)
	countAhspUsageFn func(id uint) (int64, error)
	createFn         func(u *domain.UnitPrice) error
	updateFn         func(u *domain.UnitPrice) error
	deleteFn         func(id uint) (bool, error)
}

func (m *mockUnitPriceRepo) List(q, category string) ([]domain.UnitPrice, error) {
	if m.listFn != nil {
		return m.listFn(q, category)
	}
	return nil, nil
}

func (m *mockUnitPriceRepo) FindByID(id uint) (*domain.UnitPrice, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUnitPriceRepo) FindByCode(code string, excludeID *uint) (*domain.UnitPrice, error) {
	if m.findByCodeFn != nil {
		return m.findByCodeFn(code, excludeID)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUnitPriceRepo) CountAhspUsage(id uint) (int64, error) {
	if m.countAhspUsageFn != nil {
		return m.countAhspUsageFn(id)
	}
	return 0, nil
}

func (m *mockUnitPriceRepo) Create(u *domain.UnitPrice) error {
	if m.createFn != nil {
		return m.createFn(u)
	}
	u.ID = 1
	return nil
}

func (m *mockUnitPriceRepo) Update(u *domain.UnitPrice) error {
	if m.updateFn != nil {
		return m.updateFn(u)
	}
	return nil
}

func (m *mockUnitPriceRepo) Delete(id uint) (bool, error) {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return true, nil
}

func TestUnitPriceService(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		repo := &mockUnitPriceRepo{
			listFn: func(q, category string) ([]domain.UnitPrice, error) {
				return []domain.UnitPrice{
					{ID: 1, Code: "UP-01", Category: domain.CategoryLabor, WorkItem: "Tukang", Unit: "OH", UnitPrice: 100000},
				}, nil
			},
		}
		svc := NewService(repo)
		list, err := svc.List("", "")
		if err != nil || len(list) != 1 || list[0].Code != "UP-01" {
			t.Fatalf("unexpected list: %v, list: %+v", err, list)
		}
	})

	t.Run("Create validation & duplicate check", func(t *testing.T) {
		repo := &mockUnitPriceRepo{
			findByCodeFn: func(code string, excludeID *uint) (*domain.UnitPrice, error) {
				if code == "DUP-01" {
					return &domain.UnitPrice{ID: 2, Code: "DUP-01"}, nil
				}
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := NewService(repo)

		// Required fields
		_, err := svc.Create(CreateUnitPriceInput{Code: "", WorkItem: "Item", Unit: "m3"})
		if !errors.Is(err, ErrRequiredFields) {
			t.Fatalf("expected ErrRequiredFields, got %v", err)
		}

		// Duplicate code
		_, err = svc.Create(CreateUnitPriceInput{Code: "DUP-01", WorkItem: "Item", Unit: "m3"})
		if !errors.Is(err, ErrCodeAlreadyExists) {
			t.Fatalf("expected ErrCodeAlreadyExists, got %v", err)
		}

		// Success
		created, err := svc.Create(CreateUnitPriceInput{
			Code:      "UP-NEW",
			Category:  "LABOR",
			WorkItem:  "Mandor",
			Unit:      "OH",
			UnitPrice: "150000",
		})
		if err != nil || created.Code != "UP-NEW" {
			t.Fatalf("unexpected create: %v, created: %+v", err, created)
		}
	})

	t.Run("Update not found and duplicate code", func(t *testing.T) {
		repo := &mockUnitPriceRepo{
			findByIDFn: func(id uint) (*domain.UnitPrice, error) {
				if id == 1 {
					return &domain.UnitPrice{ID: 1, Code: "UP-01", Category: domain.CategoryLabor, WorkItem: "Tukang", Unit: "OH", UnitPrice: 100000}, nil
				}
				return nil, gorm.ErrRecordNotFound
			},
			findByCodeFn: func(code string, excludeID *uint) (*domain.UnitPrice, error) {
				if code == "EXISTING" {
					return &domain.UnitPrice{ID: 2, Code: "EXISTING"}, nil
				}
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := NewService(repo)

		// Not found
		_, err := svc.Update(99, UpdateUnitPriceInput{})
		if !errors.Is(err, ErrUnitPriceNotFound) {
			t.Fatalf("expected ErrUnitPriceNotFound, got %v", err)
		}

		// Duplicate code
		dup := "EXISTING"
		_, err = svc.Update(1, UpdateUnitPriceInput{Code: &dup})
		if !errors.Is(err, ErrCodeAlreadyExists) {
			t.Fatalf("expected ErrCodeAlreadyExists, got %v", err)
		}

		// Success update
		newItem := "Tukang Kayu"
		newPrice := "120000"
		updated, err := svc.Update(1, UpdateUnitPriceInput{
			WorkItem:  &newItem,
			UnitPrice: &newPrice,
		})
		if err != nil || updated.WorkItem != "Tukang Kayu" || updated.UnitPrice != "120000" {
			t.Fatalf("unexpected update error: %v", err)
		}
	})

	t.Run("Delete referenced by AHSP error and success", func(t *testing.T) {
		repo := &mockUnitPriceRepo{
			countAhspUsageFn: func(id uint) (int64, error) {
				if id == 10 {
					return 3, nil // in use
				}
				return 0, nil
			},
			deleteFn: func(id uint) (bool, error) {
				return id == 1, nil
			},
		}
		svc := NewService(repo)

		// In use
		err := svc.Delete(10)
		if !errors.Is(err, ErrReferencedByAhspItem) {
			t.Fatalf("expected ErrReferencedByAhspItem, got %v", err)
		}

		// Not found
		err = svc.Delete(99)
		if !errors.Is(err, ErrUnitPriceNotFound) {
			t.Fatalf("expected ErrUnitPriceNotFound, got %v", err)
		}

		// Success
		err = svc.Delete(1)
		if err != nil {
			t.Fatalf("expected success, got %v", err)
		}
	})
}

func TestUnitPriceHandler(t *testing.T) {
	repo := &mockUnitPriceRepo{
		listFn: func(q, category string) ([]domain.UnitPrice, error) {
			return []domain.UnitPrice{
				{ID: 1, Code: "UP-01", Category: domain.CategoryLabor, WorkItem: "Tukang", Unit: "OH", UnitPrice: 100000},
			}, nil
		},
	}
	svc := NewService(repo)
	handler := NewHandler(svc)

	app := fiber.New()
	app.Get("/unit-prices", handler.List)
	app.Post("/unit-prices", handler.Create)
	app.Put("/unit-prices/:id", handler.Update)
	app.Delete("/unit-prices/:id", handler.Delete)

	t.Run("HTTP List", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/unit-prices?category=LABOR", nil)
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("HTTP Create success", func(t *testing.T) {
		body := `{"code":"UP-02","category":"MATERIALS","work_item":"Semen","unit":"sak","unit_price":"65000"}`
		req := httptest.NewRequest(http.MethodPost, "/unit-prices", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201, got %d, err: %v", resp.StatusCode, err)
		}
	})
}

