package volume

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

type mockVolumeRepo struct {
	listRowsFn             func(projectID uint, q string) ([]VolumeRowQueryResult, error)
	findVolumeByIDFn       func(id uint) (*domain.Volume, error)
	findBoqItemByIDFn      func(id uint) (*domain.BoqItem, error)
	findBoqItemInProjectFn func(id uint, projectID uint) (*domain.BoqItem, error)
	createBoqItemFn        func(b *domain.BoqItem) error
	createVolumeFn         func(v *domain.Volume) error
	updateVolumeFn         func(v *domain.Volume) error
	deleteVolumeFn         func(id uint) (bool, error)
}

func (m *mockVolumeRepo) ListRows(projectID uint, q string) ([]VolumeRowQueryResult, error) {
	if m.listRowsFn != nil {
		return m.listRowsFn(projectID, q)
	}
	return nil, nil
}

func (m *mockVolumeRepo) FindVolumeByID(id uint) (*domain.Volume, error) {
	if m.findVolumeByIDFn != nil {
		return m.findVolumeByIDFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockVolumeRepo) FindBoqItemByID(id uint) (*domain.BoqItem, error) {
	if m.findBoqItemByIDFn != nil {
		return m.findBoqItemByIDFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockVolumeRepo) FindBoqItemInProject(id uint, projectID uint) (*domain.BoqItem, error) {
	if m.findBoqItemInProjectFn != nil {
		return m.findBoqItemInProjectFn(id, projectID)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockVolumeRepo) CreateBoqItem(b *domain.BoqItem) error {
	if m.createBoqItemFn != nil {
		return m.createBoqItemFn(b)
	}
	b.ID = 10
	return nil
}

func (m *mockVolumeRepo) CreateVolume(v *domain.Volume) error {
	if m.createVolumeFn != nil {
		return m.createVolumeFn(v)
	}
	v.ID = 1
	return nil
}

func (m *mockVolumeRepo) UpdateVolume(v *domain.Volume) error {
	if m.updateVolumeFn != nil {
		return m.updateVolumeFn(v)
	}
	return nil
}

func (m *mockVolumeRepo) DeleteVolume(id uint) (bool, error) {
	if m.deleteVolumeFn != nil {
		return m.deleteVolumeFn(id)
	}
	return true, nil
}

func TestCalcQuantities(t *testing.T) {
	t.Run("with dimensions", func(t *testing.T) {
		l, w, h := 10.0, 2.0, 3.0
		qty, total := calcQuantities(&l, &w, &h, 0, 2.5)
		if qty != 60.0 {
			t.Errorf("expected qty 60, got %f", qty)
		}
		if total != 150.0 {
			t.Errorf("expected total 150, got %f", total)
		}
	})

	t.Run("without dimensions uses planned volume", func(t *testing.T) {
		qty, total := calcQuantities(nil, nil, nil, 45.0, 0)
		if qty != 45.0 || total != 45.0 {
			t.Errorf("expected 45/45, got %f/%f", qty, total)
		}
	})
}

func TestVolumeService(t *testing.T) {
	now := time.Now()
	l, w, h := 5.0, 2.0, 1.0

	t.Run("List", func(t *testing.T) {
		repo := &mockVolumeRepo{
			listRowsFn: func(projectID uint, q string) ([]VolumeRowQueryResult, error) {
				return []VolumeRowQueryResult{
					{
						Volume: domain.Volume{
							ID:            1,
							BoqItemID:     10,
							Date:          &now,
							Length:        &l,
							Width:         &w,
							Height:        &h,
							Quantity:      10,
							UnitWeight:    1,
							TotalQuantity: 10,
							Description:   "Section 1",
						},
						WorkItem:      "Galian",
						Unit:          "m3",
						PlannedVolume: 10,
					},
				}, nil
			},
		}
		svc := NewService(repo)
		rows, err := svc.List(1, "")
		if err != nil || len(rows) != 1 || rows[0].WorkItem != "Galian" {
			t.Fatalf("unexpected rows: %v, rows: %+v", err, rows)
		}
	})

	t.Run("Create with existing boq item", func(t *testing.T) {
		repo := &mockVolumeRepo{
			findBoqItemInProjectFn: func(id uint, projectID uint) (*domain.BoqItem, error) {
				if id == 10 {
					return &domain.BoqItem{ID: 10, ProjectID: projectID, WorkItem: "Galian", Unit: "m3"}, nil
				}
				return nil, gorm.ErrRecordNotFound
			},
		}
		svc := NewService(repo)

		// Non-existent boq item
		_, err := svc.Create(1, CreateVolumeInput{BoqItemID: strPtr("99")})
		if !errors.Is(err, ErrBoqItemNotFound) {
			t.Fatalf("expected ErrBoqItemNotFound, got %v", err)
		}

		// Valid boq item
		created, err := svc.Create(1, CreateVolumeInput{
			BoqItemID:  strPtr("10"),
			Length:     strPtr("5"),
			Width:      strPtr("2"),
			Height:     strPtr("1"),
			UnitWeight: "1",
		})
		if err != nil || created.Quantity != "10.0000" {
			t.Fatalf("unexpected create: %v, created: %+v", err, created)
		}
	})

	t.Run("Create with new boq item", func(t *testing.T) {
		svc := NewService(&mockVolumeRepo{})

		// Missing required field
		_, err := svc.Create(1, CreateVolumeInput{WorkItem: "", Unit: "m3"})
		if !errors.Is(err, ErrRequiredBoqField) {
			t.Fatalf("expected ErrRequiredBoqField, got %v", err)
		}

		// Success
		created, err := svc.Create(1, CreateVolumeInput{
			WorkItem:      "Pengecoran",
			Unit:          "m3",
			PlannedVolume: "25",
		})
		if err != nil || created.WorkItem != "Pengecoran" {
			t.Fatalf("unexpected create: %v, created: %+v", err, created)
		}
	})

	t.Run("Update not found and success", func(t *testing.T) {
		repo := &mockVolumeRepo{
			findVolumeByIDFn: func(id uint) (*domain.Volume, error) {
				if id == 1 {
					return &domain.Volume{ID: 1, BoqItemID: 10, Quantity: 10, UnitWeight: 1, TotalQuantity: 10}, nil
				}
				return nil, gorm.ErrRecordNotFound
			},
			findBoqItemByIDFn: func(id uint) (*domain.BoqItem, error) {
				return &domain.BoqItem{ID: 10, WorkItem: "Galian", Unit: "m3"}, nil
			},
		}
		svc := NewService(repo)

		// Not found
		_, err := svc.Update(99, UpdateVolumeInput{})
		if !errors.Is(err, ErrVolumeNotFound) {
			t.Fatalf("expected ErrVolumeNotFound, got %v", err)
		}

		// Success
		newDesc := "Updated Desc"
		updated, err := svc.Update(1, UpdateVolumeInput{
			Description: &newDesc,
		})
		if err != nil || updated.Description != "Updated Desc" {
			t.Fatalf("unexpected update: %v, updated: %+v", err, updated)
		}
	})

	t.Run("Delete not found and success", func(t *testing.T) {
		repo := &mockVolumeRepo{
			deleteVolumeFn: func(id uint) (bool, error) {
				return id == 1, nil
			},
		}
		svc := NewService(repo)

		if err := svc.Delete(1); err != nil {
			t.Fatalf("expected success, got %v", err)
		}
		if err := svc.Delete(99); !errors.Is(err, ErrVolumeNotFound) {
			t.Fatalf("expected ErrVolumeNotFound, got %v", err)
		}
	})
}

func TestVolumeHandler(t *testing.T) {
	repo := &mockVolumeRepo{
		listRowsFn: func(projectID uint, q string) ([]VolumeRowQueryResult, error) {
			return []VolumeRowQueryResult{{Volume: domain.Volume{ID: 1}, WorkItem: "Galian", Unit: "m3"}}, nil
		},
	}
	svc := NewService(repo)
	handler := NewHandler(svc)

	app := fiber.New()
	app.Get("/projects/:id/volumes", handler.List)
	app.Post("/projects/:id/volumes", handler.Create)
	app.Put("/volumes/:id", handler.Update)
	app.Delete("/volumes/:id", handler.Delete)

	t.Run("HTTP List", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/projects/1/volumes", nil)
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("HTTP Create success", func(t *testing.T) {
		body := `{"work_item":"Pasang Bata","unit":"m2","planned_volume":"50"}`
		req := httptest.NewRequest(http.MethodPost, "/projects/1/volumes", strings.NewReader(body))
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

