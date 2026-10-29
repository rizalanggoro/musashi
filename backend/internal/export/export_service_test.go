package export

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type mockExportRepo struct {
	getProjectFn            func(projectID uint) (*domain.Project, error)
	getVolumeExportRowsFn   func(projectID uint) ([]VolumeExportRow, error)
	getRabItemsFn           func(projectID uint) ([]domain.RabItem, error)
}

func (m *mockExportRepo) GetProject(projectID uint) (*domain.Project, error) {
	if m.getProjectFn != nil {
		return m.getProjectFn(projectID)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockExportRepo) GetVolumeExportRows(projectID uint) ([]VolumeExportRow, error) {
	if m.getVolumeExportRowsFn != nil {
		return m.getVolumeExportRowsFn(projectID)
	}
	return nil, nil
}

func (m *mockExportRepo) GetRabItems(projectID uint) ([]domain.RabItem, error) {
	if m.getRabItemsFn != nil {
		return m.getRabItemsFn(projectID)
	}
	return nil, nil
}

func TestFormatNumber(t *testing.T) {
	if got := formatNumber(1000000); got != "1.000.000" {
		t.Errorf("formatNumber(1000000) = %s, want '1.000.000'", got)
	}
	if got := formatNumber(500); got != "500" {
		t.Errorf("formatNumber(500) = %s, want '500'", got)
	}
	if got := formatNumber(0); got != "0" {
		t.Errorf("formatNumber(0) = %s, want '0'", got)
	}
}

func TestExportService(t *testing.T) {
	now := time.Now()
	l, w, h := 2.0, 3.0, 4.0

	mockProj := &domain.Project{
		ID:      1,
		Name:    "Gedung Musashi",
		Section: "Sipil",
		Plant:   "Plant 1",
	}

	repo := &mockExportRepo{
		getProjectFn: func(projectID uint) (*domain.Project, error) {
			if projectID == 1 {
				return mockProj, nil
			}
			return nil, gorm.ErrRecordNotFound
		},
		getVolumeExportRowsFn: func(projectID uint) ([]VolumeExportRow, error) {
			return []VolumeExportRow{
				{
					Volume: domain.Volume{
						ID:            1,
						Date:          &now,
						Length:        &l,
						Width:         &w,
						Height:        &h,
						Quantity:      24,
						UnitWeight:    1,
						TotalQuantity: 24,
						Description:   "Section A",
					},
					WorkItem:      "Galian",
					Unit:          "m3",
				},
			}, nil
		},
		getRabItemsFn: func(projectID uint) ([]domain.RabItem, error) {
			return []domain.RabItem{
				{
					ID:        1,
					ProjectID: projectID,
					WorkItem:  "Galian",
					Unit:      "m3",
					Volume:    24,
					UnitPrice: 50000,
					TotalCost: 1200000,
				},
			}, nil
		},
	}

	svc := NewService(repo)

	t.Run("GenerateVolumeExcel success & not found", func(t *testing.T) {
		name, bytes, err := svc.GenerateVolumeExcel(1)
		if err != nil || len(bytes) == 0 || name == "" {
			t.Fatalf("unexpected excel export: %v, len: %d, name: %s", err, len(bytes), name)
		}

		_, _, err = svc.GenerateVolumeExcel(99)
		if !errors.Is(err, ErrProjectNotFound) {
			t.Fatalf("expected ErrProjectNotFound, got %v", err)
		}
	})

	t.Run("GenerateRabExcel success & not found", func(t *testing.T) {
		name, bytes, err := svc.GenerateRabExcel(1)
		if err != nil || len(bytes) == 0 || name == "" {
			t.Fatalf("unexpected rab excel export: %v, len: %d", err, len(bytes))
		}

		_, _, err = svc.GenerateRabExcel(99)
		if !errors.Is(err, ErrProjectNotFound) {
			t.Fatalf("expected ErrProjectNotFound, got %v", err)
		}
	})

	t.Run("GenerateRabPdf success & not found", func(t *testing.T) {
		name, bytes, err := svc.GenerateRabPdf(1)
		if err != nil || len(bytes) == 0 || name == "" {
			t.Fatalf("unexpected rab pdf export: %v, len: %d", err, len(bytes))
		}

		_, _, err = svc.GenerateRabPdf(99)
		if !errors.Is(err, ErrProjectNotFound) {
			t.Fatalf("expected ErrProjectNotFound, got %v", err)
		}
	})
}

func TestExportHandler(t *testing.T) {
	mockProj := &domain.Project{ID: 1, Name: "Proj 1"}
	repo := &mockExportRepo{
		getProjectFn: func(projectID uint) (*domain.Project, error) {
			return mockProj, nil
		},
	}
	svc := NewService(repo)
	handler := NewHandler(svc)

	app := fiber.New()
	app.Get("/export/volume/:id/excel", handler.VolumeExcel)
	app.Get("/export/rab/:id/excel", handler.RabExcel)
	app.Get("/export/rab/:id/pdf", handler.RabPdf)

	t.Run("HTTP ExportVolumeExcel", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/export/volume/1/excel", nil)
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("HTTP ExportRabExcel", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/export/rab/1/excel", nil)
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("HTTP ExportRabPdf", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/export/rab/1/pdf", nil)
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d, err: %v", resp.StatusCode, err)
		}
	})
}
