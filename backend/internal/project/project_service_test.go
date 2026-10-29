package project

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

type mockProjectRepo struct {
	listFn       func(q string, page, perPage int) ([]domain.Project, int64, error)
	findByIDFn   func(id uint) (*domain.Project, error)
	findByNameFn func(name string, excludeID *uint) (*domain.Project, error)
	createFn     func(p *domain.Project) error
	updateFn     func(p *domain.Project) error
	deleteFn     func(id uint) (bool, error)
}

func (m *mockProjectRepo) List(q string, page, perPage int) ([]domain.Project, int64, error) {
	if m.listFn != nil {
		return m.listFn(q, page, perPage)
	}
	return nil, 0, nil
}

func (m *mockProjectRepo) FindByID(id uint) (*domain.Project, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockProjectRepo) FindByName(name string, excludeID *uint) (*domain.Project, error) {
	if m.findByNameFn != nil {
		return m.findByNameFn(name, excludeID)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockProjectRepo) Create(p *domain.Project) error {
	if m.createFn != nil {
		return m.createFn(p)
	}
	p.ID = 1
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

func (m *mockProjectRepo) Update(p *domain.Project) error {
	if m.updateFn != nil {
		return m.updateFn(p)
	}
	p.UpdatedAt = time.Now()
	return nil
}

func (m *mockProjectRepo) Delete(id uint) (bool, error) {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return true, nil
}

func TestProjectService(t *testing.T) {
	t.Run("List default pagination and query", func(t *testing.T) {
		repo := &mockProjectRepo{
			listFn: func(q string, page, perPage int) ([]domain.Project, int64, error) {
				if page != 1 || perPage != 25 {
					t.Errorf("expected page=1, perPage=25, got %d, %d", page, perPage)
				}
				return []domain.Project{
					{ID: 1, Name: "Proj 1"},
					{ID: 2, Name: "Proj 2"},
				}, 2, nil
			},
		}
		svc := NewService(repo)
		list, total, err := svc.List("", 0, 0)
		if err != nil || total != 2 || len(list) != 2 {
			t.Fatalf("unexpected list result: %v, total: %d", err, total)
		}
	})

	t.Run("Get found and not found", func(t *testing.T) {
		repo := &mockProjectRepo{
			findByIDFn: func(id uint) (*domain.Project, error) {
				if id == 1 {
					return &domain.Project{ID: 1, Name: "Proj 1"}, nil
				}
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := NewService(repo)

		got, err := svc.Get(1)
		if err != nil || got.Name != "Proj 1" {
			t.Fatalf("unexpected get: %v, got: %+v", err, got)
		}

		_, err = svc.Get(999)
		if !errors.Is(err, ErrProjectNotFound) {
			t.Fatalf("expected ErrProjectNotFound, got %v", err)
		}
	})

	t.Run("Create validation & duplicate check", func(t *testing.T) {
		repo := &mockProjectRepo{
			findByNameFn: func(name string, excludeID *uint) (*domain.Project, error) {
				if name == "Duplicate" {
					return &domain.Project{ID: 2, Name: "Duplicate"}, nil
				}
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := NewService(repo)

		// Empty name
		_, err := svc.Create(CreateProjectInput{Name: "   "})
		if !errors.Is(err, ErrProjectNameRequired) {
			t.Fatalf("expected ErrProjectNameRequired, got %v", err)
		}

		// Duplicate name
		_, err = svc.Create(CreateProjectInput{Name: "Duplicate"})
		if !errors.Is(err, ErrProjectAlreadyExists) {
			t.Fatalf("expected ErrProjectAlreadyExists, got %v", err)
		}

		// Success
		plan := "2026-06-01"
		created, err := svc.Create(CreateProjectInput{
			Name:               "New Project",
			Section:            "Sipil",
			Plant:              "Plant 1",
			Budget:             "1000000",
			ImplementationPlan: &plan,
		})
		if err != nil || created.Name != "New Project" {
			t.Fatalf("unexpected error creating project: %v", err)
		}
	})

	t.Run("Update not found & duplicate name", func(t *testing.T) {
		repo := &mockProjectRepo{
			findByIDFn: func(id uint) (*domain.Project, error) {
				if id == 1 {
					return &domain.Project{ID: 1, Name: "Original"}, nil
				}
				return nil, gorm.ErrRecordNotFound
			},
			findByNameFn: func(name string, excludeID *uint) (*domain.Project, error) {
				if name == "Existing Other" {
					return &domain.Project{ID: 2, Name: "Existing Other"}, nil
				}
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := NewService(repo)

		// Not found
		_, err := svc.Update(99, UpdateProjectInput{Name: strPtr("New")})
		if !errors.Is(err, ErrProjectNotFound) {
			t.Fatalf("expected ErrProjectNotFound, got %v", err)
		}

		// Duplicate name
		_, err = svc.Update(1, UpdateProjectInput{Name: strPtr("Existing Other")})
		if !errors.Is(err, ErrProjectAlreadyExists) {
			t.Fatalf("expected ErrProjectAlreadyExists, got %v", err)
		}

		// Success update
		updated, err := svc.Update(1, UpdateProjectInput{
			Name:        strPtr("Updated Name"),
			Budget:      strPtr("2000000"),
			ProgressSPB: strPtr("done"),
			Status:      strPtr("on_going"),
		})
		if err != nil || updated.Name != "Updated Name" || updated.ProgressSPB != "done" {
			t.Fatalf("unexpected update error: %v", err)
		}
	})

	t.Run("Delete success & not found", func(t *testing.T) {
		repo := &mockProjectRepo{
			deleteFn: func(id uint) (bool, error) {
				return id == 1, nil
			},
		}
		svc := NewService(repo)

		if err := svc.Delete(1); err != nil {
			t.Fatalf("expected success, got %v", err)
		}

		if err := svc.Delete(99); !errors.Is(err, ErrProjectNotFound) {
			t.Fatalf("expected ErrProjectNotFound, got %v", err)
		}
	})
}

func TestProjectHandler(t *testing.T) {
	repo := &mockProjectRepo{
		listFn: func(q string, page, perPage int) ([]domain.Project, int64, error) {
			return []domain.Project{{ID: 1, Name: "Gedung A"}}, 1, nil
		},
		findByIDFn: func(id uint) (*domain.Project, error) {
			if id == 1 {
				return &domain.Project{ID: 1, Name: "Gedung A"}, nil
			}
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := NewService(repo)
	handler := NewHandler(svc)

	app := fiber.New()
	app.Get("/projects", handler.List)
	app.Get("/projects/:id", handler.Get)
	app.Post("/projects", handler.Create)
	app.Put("/projects/:id", handler.Update)
	app.Delete("/projects/:id", handler.Delete)

	t.Run("HTTP List", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/projects?page=1&per_page=10", nil)
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("HTTP Get existing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/projects/1", nil)
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("HTTP Get not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/projects/99", nil)
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("HTTP Create success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/projects", strings.NewReader(`{"name":"New Proj"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201, got %d, err: %v", resp.StatusCode, err)
		}
	})
}

func strPtr(s string) *string {
	return &s
}

