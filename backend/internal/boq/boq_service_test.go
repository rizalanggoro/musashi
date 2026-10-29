package boq

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"

	"musashi/backend/internal/domain"
)

type mockBoqRepo struct {
	listProjectsFn     func(q string) ([]domain.Project, error)
	getProjectCountsFn func(projectID uint) (ProjectCounts, error)
	getBoqItemsFn      func(projectID uint) ([]domain.BoqItem, error)
}

func (m *mockBoqRepo) ListProjects(q string) ([]domain.Project, error) {
	if m.listProjectsFn != nil {
		return m.listProjectsFn(q)
	}
	return nil, nil
}

func (m *mockBoqRepo) GetProjectCounts(projectID uint) (ProjectCounts, error) {
	if m.getProjectCountsFn != nil {
		return m.getProjectCountsFn(projectID)
	}
	return ProjectCounts{}, nil
}

func (m *mockBoqRepo) GetBoqItems(projectID uint) ([]domain.BoqItem, error) {
	if m.getBoqItemsFn != nil {
		return m.getBoqItemsFn(projectID)
	}
	return nil, nil
}

func TestCalculateStatuses(t *testing.T) {
	p := domain.Project{ID: 1}

	t.Run("empty project", func(t *testing.T) {
		vol, rab := calculateStatuses(p, ProjectCounts{})
		if vol != "empty" || rab != "empty" {
			t.Errorf("expected empty/empty, got %s/%s", vol, rab)
		}
	})

	t.Run("partial volume and complete rab", func(t *testing.T) {
		counts := ProjectCounts{
			BoqCount: 5,
			VolCount: 2,
			RabCount: 5,
		}
		vol, rab := calculateStatuses(p, counts)
		if vol != "partial" || rab != "complete" {
			t.Errorf("expected partial/complete, got %s/%s", vol, rab)
		}
	})

	t.Run("complete volume and partial rab", func(t *testing.T) {
		counts := ProjectCounts{
			BoqCount: 5,
			VolCount: 5,
			RabCount: 3,
		}
		vol, rab := calculateStatuses(p, counts)
		if vol != "complete" || rab != "partial" {
			t.Errorf("expected complete/partial, got %s/%s", vol, rab)
		}
	})
}

func TestBoqService(t *testing.T) {
	t.Run("ListSummary", func(t *testing.T) {
		repo := &mockBoqRepo{
			listProjectsFn: func(q string) ([]domain.Project, error) {
				return []domain.Project{
					{ID: 1, Name: "Proj 1", Section: "Sipil", Plant: "Plant 1"},
				}, nil
			},
			getProjectCountsFn: func(projectID uint) (ProjectCounts, error) {
				return ProjectCounts{BoqCount: 2, VolCount: 2, RabCount: 2}, nil
			},
		}
		svc := NewService(repo)
		list, err := svc.ListSummary("")
		if err != nil || len(list) != 1 || list[0].VolumeStatus != "complete" {
			t.Fatalf("unexpected summary: %v, list: %+v", err, list)
		}
	})

	t.Run("GetProjectBoq", func(t *testing.T) {
		repo := &mockBoqRepo{
			getBoqItemsFn: func(projectID uint) ([]domain.BoqItem, error) {
				return []domain.BoqItem{
					{ID: 10, ProjectID: projectID, WorkItem: "Galian", Unit: "m3", PlannedVolume: 15.5},
				}, nil
			},
		}
		svc := NewService(repo)
		items, err := svc.GetProjectBoq(1)
		if err != nil || len(items) != 1 || items[0].WorkItem != "Galian" {
			t.Fatalf("unexpected boq items: %v, items: %+v", err, items)
		}
	})
}

func TestBoqHandler(t *testing.T) {
	repo := &mockBoqRepo{
		listProjectsFn: func(q string) ([]domain.Project, error) {
			return []domain.Project{{ID: 1, Name: "P1"}}, nil
		},
		getBoqItemsFn: func(projectID uint) ([]domain.BoqItem, error) {
			return []domain.BoqItem{{ID: 10, ProjectID: projectID, WorkItem: "Item 1", Unit: "m3"}}, nil
		},
	}
	svc := NewService(repo)
	handler := NewHandler(svc)

	app := fiber.New()
	app.Get("/boq/summary", handler.ListSummary)
	app.Get("/projects/:id/boq", handler.GetProjectBoq)

	t.Run("HTTP ListSummary", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/boq/summary", nil)
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("HTTP GetProjectBoq", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/projects/1/boq", nil)
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d, err: %v", resp.StatusCode, err)
		}
	})
}

