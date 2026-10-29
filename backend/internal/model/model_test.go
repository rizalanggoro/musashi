package model

import (
	"testing"
	"time"

	"musashi/backend/internal/domain"
)

func TestUserModelMapping(t *testing.T) {
	t.Run("nil safety", func(t *testing.T) {
		var m *User
		if m.ToDomain() != nil {
			t.Error("expected nil domain for nil model")
		}
		if FromUserDomain(nil) != nil {
			t.Error("expected nil model for nil domain")
		}
	})

	t.Run("bidirectional mapping", func(t *testing.T) {
		now := time.Now()
		d := &domain.User{
			ID:           10,
			CreatedAt:    now,
			UpdatedAt:    now,
			Email:        "staff@musashi.local",
			Name:         "Staff User",
			PasswordHash: "hashedsecret",
			Role:         domain.RoleStaff,
		}

		m := FromUserDomain(d)
		if m.ID != 10 || m.Email != "staff@musashi.local" || m.Role != domain.RoleStaff {
			t.Fatalf("unexpected model: %+v", m)
		}
		if m.TableName() != "users" {
			t.Errorf("expected TableName 'users', got %s", m.TableName())
		}

		back := m.ToDomain()
		if back.ID != d.ID || back.Email != d.Email || back.PasswordHash != d.PasswordHash || back.Role != d.Role {
			t.Fatalf("roundtrip mismatch: %+v != %+v", back, d)
		}
	})
}

func TestProjectModelMapping(t *testing.T) {
	t.Run("nil safety", func(t *testing.T) {
		var m *Project
		if m.ToDomain() != nil {
			t.Error("expected nil")
		}
		if FromProjectDomain(nil) != nil {
			t.Error("expected nil")
		}
	})

	t.Run("project with child relations", func(t *testing.T) {
		planDate := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
		d := &domain.Project{
			ID:                 1,
			Name:               "Renovasi Gedung",
			Section:            "Sipil",
			Plant:              "Plant 2",
			ImplementationPlan: &planDate,
			Budget:             1000000,
			Information:        "catatan",
			Status:             domain.ProjectStatusPlanning,
			ProgressSPB:        domain.ProgressPending,
			ProgressRAB:        domain.ProgressDone,
			ProgressXXX:        domain.ProgressPending,
			Schedules: []domain.Schedule{
				{ID: 2, ProjectID: 1, Code: "SCH-01", WorkItem: "Pembersihan", Duration: 3},
			},
			BoqItems: []domain.BoqItem{
				{ID: 3, ProjectID: 1, WorkItem: "Galian", Unit: "m3", PlannedVolume: 50},
			},
			RabItems: []domain.RabItem{
				{ID: 4, ProjectID: 1, WorkItem: "Galian", Unit: "m3", Volume: 50, UnitPrice: 20000, TotalCost: 1000000},
			},
		}

		m := FromProjectDomain(d)
		if m.ID != 1 || len(m.Schedules) != 1 || len(m.BoqItems) != 1 || len(m.RabItems) != 1 {
			t.Fatalf("unexpected project model: %+v", m)
		}
		if m.TableName() != "projects" {
			t.Errorf("expected TableName 'projects', got %s", m.TableName())
		}

		back := m.ToDomain()
		if back.Name != d.Name || len(back.Schedules) != 1 || len(back.BoqItems) != 1 || len(back.RabItems) != 1 {
			t.Fatalf("roundtrip project mismatch: %+v", back)
		}
	})
}

func TestScheduleModelMapping(t *testing.T) {
	t.Run("nil safety", func(t *testing.T) {
		var m *Schedule
		if m.ToDomain() != nil {
			t.Error("expected nil")
		}
		if FromScheduleDomain(nil) != nil {
			t.Error("expected nil")
		}
	})

	t.Run("bidirectional schedule", func(t *testing.T) {
		start := time.Now()
		end := start.AddDate(0, 0, 5)
		d := &domain.Schedule{
			ID:        5,
			ProjectID: 1,
			Code:      "SCH-10",
			WorkItem:  "Pondasi",
			StartDate: start,
			EndDate:   end,
			Duration:  6,
			Weight:    15.5,
			Progress:  50.0,
			Status:    domain.ScheduleOnGoing,
		}

		m := FromScheduleDomain(d)
		if m.TableName() != "schedules" {
			t.Errorf("expected TableName 'schedules', got %s", m.TableName())
		}
		back := m.ToDomain()
		if back.ID != d.ID || back.Code != d.Code || back.Weight != d.Weight || back.Status != d.Status {
			t.Fatalf("schedule roundtrip mismatch: %+v", back)
		}
	})
}

func TestUnitPriceModelMapping(t *testing.T) {
	t.Run("nil safety", func(t *testing.T) {
		var m *UnitPrice
		if m.ToDomain() != nil {
			t.Error("expected nil")
		}
		if FromUnitPriceDomain(nil) != nil {
			t.Error("expected nil")
		}
	})

	t.Run("bidirectional unit price", func(t *testing.T) {
		d := &domain.UnitPrice{
			ID:          7,
			Code:        "UP-01",
			Category:    domain.CategoryLabor,
			WorkItem:    "Tukang Batu",
			Unit:        "OH",
			UnitPrice:   120000,
			Information: "info",
		}

		m := FromUnitPriceDomain(d)
		if m.TableName() != "unit_prices" {
			t.Errorf("expected TableName 'unit_prices', got %s", m.TableName())
		}
		back := m.ToDomain()
		if back.ID != d.ID || back.Code != d.Code || back.UnitPrice != d.UnitPrice {
			t.Fatalf("unit price roundtrip mismatch: %+v", back)
		}
	})
}

func TestAhspModelMapping(t *testing.T) {
	t.Run("nil safety", func(t *testing.T) {
		var m *Ahsp
		if m.ToDomain() != nil {
			t.Error("expected nil")
		}
		if FromAhspDomain(nil) != nil {
			t.Error("expected nil")
		}
		var item *AhspItem
		if item.ToDomain() != nil {
			t.Error("expected nil")
		}
		if FromAhspItemDomain(nil) != nil {
			t.Error("expected nil")
		}
	})

	t.Run("bidirectional ahsp with items", func(t *testing.T) {
		upID := uint(99)
		d := &domain.Ahsp{
			ID:                8,
			Code:              "AHSP-01",
			WorkItem:          "Beton",
			Unit:              "m3",
			OverheadProfitPct: 10,
			TotalCost:         500000,
			UnitPriceOfWork:   550000,
			Status:            domain.AhspStatusDraft,
			Items: []domain.AhspItem{
				{
					ID:          1,
					AhspID:      8,
					Category:    domain.CategoryMaterials,
					Sort:        1,
					Description: "Semen",
					Code:        "MAT-01",
					Unit:        "sak",
					Coefficient: 6.5,
					UnitPrice:   70000,
					Cost:        455000,
					UnitPriceID: &upID,
				},
			},
		}

		m := FromAhspDomain(d)
		if m.TableName() != "ahsps" {
			t.Errorf("expected TableName 'ahsps', got %s", m.TableName())
		}
		if len(m.Items) != 1 || m.Items[0].TableName() != "ahsp_items" {
			t.Fatalf("unexpected items in ahsp model")
		}

		back := m.ToDomain()
		if back.Code != d.Code || len(back.Items) != 1 || back.Items[0].UnitPriceID == nil || *back.Items[0].UnitPriceID != 99 {
			t.Fatalf("ahsp roundtrip mismatch: %+v", back)
		}
	})
}

func TestBoqVolumeRabMapping(t *testing.T) {
	t.Run("boq nil safety", func(t *testing.T) {
		var b *BoqItem
		if b.ToDomain() != nil || FromBoqItemDomain(nil) != nil {
			t.Error("expected nil")
		}
		var v *Volume
		if v.ToDomain() != nil || FromVolumeDomain(nil) != nil {
			t.Error("expected nil")
		}
		var r *RabItem
		if r.ToDomain() != nil || FromRabItemDomain(nil) != nil {
			t.Error("expected nil")
		}
	})

	t.Run("boq and volume bidirectional", func(t *testing.T) {
		l, w, h := 10.0, 5.0, 2.0
		b := &domain.BoqItem{
			ID:            12,
			ProjectID:     1,
			WorkItem:      "Galian Tanah",
			Unit:          "m3",
			PlannedVolume: 100,
			Volumes: []domain.Volume{
				{
					ID:            20,
					BoqItemID:     12,
					Length:        &l,
					Width:         &w,
					Height:        &h,
					Quantity:      100,
					UnitWeight:    1,
					TotalQuantity: 100,
					Description:   "Section A",
				},
			},
		}

		bm := FromBoqItemDomain(b)
		if bm.TableName() != "boq_items" || len(bm.Volumes) != 1 || bm.Volumes[0].TableName() != "volumes" {
			t.Fatalf("unexpected boq model: %+v", bm)
		}
		back := bm.ToDomain()
		if back.WorkItem != b.WorkItem || len(back.Volumes) != 1 || *back.Volumes[0].Length != 10.0 {
			t.Fatalf("boq roundtrip mismatch: %+v", back)
		}
	})

	t.Run("rab bidirectional", func(t *testing.T) {
		r := &domain.RabItem{
			ID:        15,
			ProjectID: 1,
			WorkItem:  "Galian",
			Unit:      "m3",
			Volume:    100,
			UnitPrice: 50000,
			TotalCost: 5000000,
		}
		rm := FromRabItemDomain(r)
		if rm.TableName() != "rab_items" {
			t.Errorf("expected TableName 'rab_items', got %s", rm.TableName())
		}
		back := rm.ToDomain()
		if back.TotalCost != 5000000 {
			t.Errorf("rab roundtrip mismatch: %+v", back)
		}
	})
}
