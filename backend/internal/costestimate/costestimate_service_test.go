package costestimate

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
)

type mockCostEstimateRepo struct {
	getProjectFn  func(projectID uint) (*domain.Project, error)
	getTotalRabFn func(projectID uint) (int64, error)
}

func (m *mockCostEstimateRepo) GetProject(projectID uint) (*domain.Project, error) {
	if m.getProjectFn != nil {
		return m.getProjectFn(projectID)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockCostEstimateRepo) GetTotalRab(projectID uint) (int64, error) {
	if m.getTotalRabFn != nil {
		return m.getTotalRabFn(projectID)
	}
	return 0, nil
}

func TestCostEstimateService(t *testing.T) {
	t.Run("project not found", func(t *testing.T) {
		svc := NewService(&mockCostEstimateRepo{})
		_, err := svc.GetSummary(99)
		if !errors.Is(err, ErrProjectNotFound) {
			t.Fatalf("expected ErrProjectNotFound, got %v", err)
		}
	})

	t.Run("summary with SPB pending", func(t *testing.T) {
		repo := &mockCostEstimateRepo{
			getProjectFn: func(projectID uint) (*domain.Project, error) {
				return &domain.Project{
					ID:          1,
					Budget:      1000000,
					ProgressSPB: domain.ProgressPending,
				}, nil
			},
			getTotalRabFn: func(projectID uint) (int64, error) {
				return 800000, nil
			},
		}
		svc := NewService(repo)
		summary, err := svc.GetSummary(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if summary.Budget != "1000000" || summary.Rab != "800000" || summary.Spb != "0" || summary.Variance != "200000" || summary.Sisa != "1000000" {
			t.Fatalf("unexpected summary values: %+v", summary)
		}
	})

	t.Run("summary with SPB done", func(t *testing.T) {
		repo := &mockCostEstimateRepo{
			getProjectFn: func(projectID uint) (*domain.Project, error) {
				return &domain.Project{
					ID:          1,
					Budget:      1000000,
					ProgressSPB: domain.ProgressDone,
				}, nil
			},
			getTotalRabFn: func(projectID uint) (int64, error) {
				return 850000, nil
			},
		}
		svc := NewService(repo)
		summary, err := svc.GetSummary(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if summary.Spb != "850000" || summary.Sisa != "150000" {
			t.Fatalf("unexpected summary for SPB done: %+v", summary)
		}
	})
}

func TestCostEstimateHandler(t *testing.T) {
	repo := &mockCostEstimateRepo{
		getProjectFn: func(projectID uint) (*domain.Project, error) {
			return &domain.Project{ID: 1, Budget: 500000}, nil
		},
	}
	svc := NewService(repo)
	handler := NewHandler(svc)

	app := fiber.New()
	app.Get("/projects/:id/cost-estimate", handler.GetSummary)

	t.Run("HTTP GetSummary", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/projects/1/cost-estimate", nil)
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d, err: %v", resp.StatusCode, err)
		}
	})
}

