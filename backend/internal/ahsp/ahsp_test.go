package ahsp

import (
	"testing"
)

func TestAhspCalculations(t *testing.T) {
	// Case 1: Labor only
	t.Run("Labor only", func(t *testing.T) {
		input := CreateAhspInput{
			Code:              "AHSP-LAB-01",
			WorkItem:          "Tukang Galian",
			Unit:              "OH",
			OverheadProfitPct: "10",
			Items: []CreateAhspItemInput{
				{
					Category:    "LABOR",
					Description: "Pekerja",
					Unit:        "OH",
					Coefficient: "1.0000",
					UnitPrice:   "100000",
				},
				{
					Category:    "LABOR",
					Description: "Mandor",
					Unit:        "OH",
					Coefficient: "0.1000",
					UnitPrice:   "150000",
				},
			},
		}

		ahsp, _, err := calculateAndBuild(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// 1.0 * 100,000 = 100,000; 0.1 * 150,000 = 15,000 => Labor = 115,000
		if ahsp.TotalLabor != 115000 {
			t.Errorf("expected TotalLabor 115000, got %d", ahsp.TotalLabor)
		}
		if ahsp.TotalCost != 115000 {
			t.Errorf("expected TotalCost 115000, got %d", ahsp.TotalCost)
		}
		// Overhead: 10% of 115,000 = 11,500 => F = 126,500
		if ahsp.UnitPriceOfWork != 126500 {
			t.Errorf("expected UnitPriceOfWork 126500, got %d", ahsp.UnitPriceOfWork)
		}
	})

	// Case 2: Material only
	t.Run("Material only", func(t *testing.T) {
		input := CreateAhspInput{
			Code:              "AHSP-MAT-01",
			WorkItem:          "Pasir Urug",
			Unit:              "m3",
			OverheadProfitPct: "15",
			Items: []CreateAhspItemInput{
				{
					Category:    "MATERIALS",
					Description: "Pasir",
					Unit:        "m3",
					Coefficient: "1.2000",
					UnitPrice:   "250000",
				},
			},
		}

		ahsp, _, err := calculateAndBuild(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// 1.2 * 250,000 = 300,000
		if ahsp.TotalMaterial != 300000 {
			t.Errorf("expected TotalMaterial 300000, got %d", ahsp.TotalMaterial)
		}
		// Overhead: 15% of 300,000 = 45,000 => F = 345,000
		if ahsp.UnitPriceOfWork != 345000 {
			t.Errorf("expected UnitPriceOfWork 345000, got %d", ahsp.UnitPriceOfWork)
		}
	})

	// Case 3: Mixed (Labor + Material + Equipment)
	t.Run("Mixed", func(t *testing.T) {
		input := CreateAhspInput{
			Code:              "AHSP-MIX-01",
			WorkItem:          "Beton Cor K-250",
			Unit:              "m3",
			OverheadProfitPct: "10",
			Items: []CreateAhspItemInput{
				{
					Category:    "LABOR",
					Description: "Tukang",
					Unit:        "OH",
					Coefficient: "0.2500",
					UnitPrice:   "120000",
				},
				{
					Category:    "MATERIALS",
					Description: "Semen Portland",
					Unit:        "sak",
					Coefficient: "6.5000",
					UnitPrice:   "70000",
				},
				{
					Category:    "EQUIPMENT",
					Description: "Molen Mixer",
					Unit:        "hari",
					Coefficient: "0.1000",
					UnitPrice:   "200000",
				},
			},
		}

		ahsp, _, err := calculateAndBuild(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Labor: 0.25 * 120,000 = 30,000
		// Material: 6.5 * 70,000 = 455,000
		// Equipment: 0.1 * 200,000 = 20,000
		// Total Cost D = 30,000 + 455,000 + 20,000 = 505,000
		// Overhead E = 10% * 505,000 = 50,500
		// F = 505,000 + 50,500 = 555,500
		if ahsp.TotalLabor != 30000 {
			t.Errorf("expected TotalLabor 30000, got %d", ahsp.TotalLabor)
		}
		if ahsp.TotalMaterial != 455000 {
			t.Errorf("expected TotalMaterial 455000, got %d", ahsp.TotalMaterial)
		}
		if ahsp.TotalEquipment != 20000 {
			t.Errorf("expected TotalEquipment 20000, got %d", ahsp.TotalEquipment)
		}
		if ahsp.TotalCost != 505000 {
			t.Errorf("expected TotalCost 505000, got %d", ahsp.TotalCost)
		}
		if ahsp.UnitPriceOfWork != 555500 {
			t.Errorf("expected UnitPriceOfWork 555500, got %d", ahsp.UnitPriceOfWork)
		}
	})
}

