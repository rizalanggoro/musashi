package ahsp

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

type mockAhspRepo struct {
	listFn       func(q string) ([]domain.Ahsp, error)
	findByIDFn   func(id uint) (*domain.Ahsp, error)
	findByCodeFn func(code string, excludeID *uint) (*domain.Ahsp, error)
	createFn     func(ahsp *domain.Ahsp, items []domain.AhspItem) error
	updateFn     func(id uint, ahsp *domain.Ahsp, items []domain.AhspItem) error
	deleteFn     func(id uint) (bool, error)
}

func (m *mockAhspRepo) List(q string) ([]domain.Ahsp, error) {
	if m.listFn != nil {
		return m.listFn(q)
	}
	return nil, nil
}

func (m *mockAhspRepo) FindByID(id uint) (*domain.Ahsp, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAhspRepo) FindByCode(code string, excludeID *uint) (*domain.Ahsp, error) {
	if m.findByCodeFn != nil {
		return m.findByCodeFn(code, excludeID)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAhspRepo) Create(ahsp *domain.Ahsp, items []domain.AhspItem) error {
	if m.createFn != nil {
		return m.createFn(ahsp, items)
	}
	ahsp.ID = 1
	return nil
}

func (m *mockAhspRepo) Update(id uint, ahsp *domain.Ahsp, items []domain.AhspItem) error {
	if m.updateFn != nil {
		return m.updateFn(id, ahsp, items)
	}
	return nil
}

func (m *mockAhspRepo) Delete(id uint) (bool, error) {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return true, nil
}

func TestAhspService(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		repo := &mockAhspRepo{
			listFn: func(q string) ([]domain.Ahsp, error) {
				return []domain.Ahsp{
					{ID: 1, Code: "AHSP-01", WorkItem: "Galian", Unit: "m3"},
				}, nil
			},
		}
		svc := NewService(repo)
		list, err := svc.List("")
		if err != nil || len(list) != 1 || list[0].Code != "AHSP-01" {
			t.Fatalf("unexpected list: %v, list: %+v", err, list)
		}
	})

	t.Run("Get found and not found", func(t *testing.T) {
		repo := &mockAhspRepo{
			findByIDFn: func(id uint) (*domain.Ahsp, error) {
				if id == 1 {
					return &domain.Ahsp{ID: 1, Code: "AHSP-01", WorkItem: "Galian", Unit: "m3"}, nil
				}
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := NewService(repo)

		got, err := svc.Get(1)
		if err != nil || got.Code != "AHSP-01" {
			t.Fatalf("unexpected get: %v, got: %+v", err, got)
		}

		_, err = svc.Get(99)
		if !errors.Is(err, ErrAhspNotFound) {
			t.Fatalf("expected ErrAhspNotFound, got %v", err)
		}
	})

	t.Run("Create validation & duplicate code", func(t *testing.T) {
		repo := &mockAhspRepo{
			findByCodeFn: func(code string, excludeID *uint) (*domain.Ahsp, error) {
				if code == "EXISTING" {
					return &domain.Ahsp{ID: 2, Code: "EXISTING"}, nil
				}
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := NewService(repo)

		// Required fields
		_, err := svc.Create(CreateAhspInput{Code: "", WorkItem: "Item", Unit: "m3"})
		if !errors.Is(err, ErrRequiredFields) {
			t.Fatalf("expected ErrRequiredFields, got %v", err)
		}

		// Duplicate code
		_, err = svc.Create(CreateAhspInput{Code: "EXISTING", WorkItem: "Item", Unit: "m3"})
		if !errors.Is(err, ErrCodeAlreadyExists) {
			t.Fatalf("expected ErrCodeAlreadyExists, got %v", err)
		}

		// Success
		created, err := svc.Create(CreateAhspInput{
			Code:     "NEW-AHSP",
			WorkItem: "Pengecatan",
			Unit:     "m2",
		})
		if err != nil || created.Code != "NEW-AHSP" {
			t.Fatalf("unexpected create: %v, created: %+v", err, created)
		}
	})

	t.Run("Update not found and duplicate code", func(t *testing.T) {
		repo := &mockAhspRepo{
			findByIDFn: func(id uint) (*domain.Ahsp, error) {
				if id == 1 {
					return &domain.Ahsp{ID: 1, Code: "AHSP-01"}, nil
				}
				return nil, gorm.ErrRecordNotFound
			},
			findByCodeFn: func(code string, excludeID *uint) (*domain.Ahsp, error) {
				if code == "DUP" {
					return &domain.Ahsp{ID: 2, Code: "DUP"}, nil
				}
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := NewService(repo)

		// Not found
		_, err := svc.Update(99, CreateAhspInput{Code: "A", WorkItem: "B", Unit: "C"})
		if !errors.Is(err, ErrAhspNotFound) {
			t.Fatalf("expected ErrAhspNotFound, got %v", err)
		}

		// Duplicate code
		_, err = svc.Update(1, CreateAhspInput{Code: "DUP", WorkItem: "B", Unit: "C"})
		if !errors.Is(err, ErrCodeAlreadyExists) {
			t.Fatalf("expected ErrCodeAlreadyExists, got %v", err)
		}

		// Success update
		updated, err := svc.Update(1, CreateAhspInput{Code: "UPDATED", WorkItem: "Item 2", Unit: "m3"})
		if err != nil || updated.Code != "UPDATED" {
			t.Fatalf("unexpected update: %v, updated: %+v", err, updated)
		}
	})

	t.Run("Delete not found and success", func(t *testing.T) {
		repo := &mockAhspRepo{
			deleteFn: func(id uint) (bool, error) {
				return id == 1, nil
			},
		}
		svc := NewService(repo)

		if err := svc.Delete(1); err != nil {
			t.Fatalf("expected success, got %v", err)
		}
		if err := svc.Delete(99); !errors.Is(err, ErrAhspNotFound) {
			t.Fatalf("expected ErrAhspNotFound, got %v", err)
		}
	})
}

func TestAhspHandler(t *testing.T) {
	repo := &mockAhspRepo{
		listFn: func(q string) ([]domain.Ahsp, error) {
			return []domain.Ahsp{{ID: 1, Code: "AHSP-01", WorkItem: "Batu Kali", Unit: "m3"}}, nil
		},
		findByIDFn: func(id uint) (*domain.Ahsp, error) {
			if id == 1 {
				return &domain.Ahsp{ID: 1, Code: "AHSP-01", WorkItem: "Batu Kali", Unit: "m3"}, nil
			}
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := NewService(repo)
	handler := NewHandler(svc)

	app := fiber.New()
	app.Get("/ahsps", handler.List)
	app.Get("/ahsps/:id", handler.Get)
	app.Post("/ahsps", handler.Create)
	app.Put("/ahsps/:id", handler.Update)
	app.Delete("/ahsps/:id", handler.Delete)

	t.Run("HTTP List", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ahsps", nil)
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("HTTP Get", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ahsps/1", nil)
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("HTTP Create success", func(t *testing.T) {
		body := `{"code":"AHSP-NEW","work_item":"Batu","unit":"m3","overhead_profit_pct":"10","items":[]}`
		req := httptest.NewRequest(http.MethodPost, "/ahsps", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201, got %d, err: %v", resp.StatusCode, err)
		}
	})
}

