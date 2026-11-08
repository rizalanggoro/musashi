package seeder

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/model"
)

// Seed populates the database with default dev data idempotently.
func Seed(db *gorm.DB) error {
	slog.Info("starting database seeder...")

	if err := seedUsers(db); err != nil {
		return fmt.Errorf("seed users: %w", err)
	}

	unitPriceMap, err := seedUnitPrices(db)
	if err != nil {
		return fmt.Errorf("seed unit prices: %w", err)
	}

	ahspMap, err := seedAhsps(db, unitPriceMap)
	if err != nil {
		return fmt.Errorf("seed ahsps: %w", err)
	}

	if err := seedProjects(db, ahspMap); err != nil {
		return fmt.Errorf("seed projects: %w", err)
	}

	slog.Info("database seeding completed successfully")
	return nil
}

func seedUsers(db *gorm.DB) error {
	adminHash := "$2a$10$i/3zZPOBK/Bw79ayDXks9u3YblN8XnEAQqruvWzi2uwzl4u/TQerS" // admin123
	staffHash, err := bcrypt.GenerateFromPassword([]byte("staff123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users := []model.User{
		{
			Email:        "admin@musashi.local",
			Name:         "Administrator",
			PasswordHash: adminHash,
			Role:         domain.RoleAdmin,
		},
		{
			Email:        "staff@musashi.local",
			Name:         "Staff Estimator",
			PasswordHash: string(staffHash),
			Role:         domain.RoleStaff,
		},
	}

	for _, u := range users {
		var existing model.User
		err := db.Where("lower(email) = lower(?) AND deleted_at IS NULL", u.Email).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&u).Error; err != nil {
				return err
			}
			slog.Info("seeded user", "email", u.Email, "role", u.Role)
		} else if err != nil {
			return err
		}
	}

	return nil
}

func seedUnitPrices(db *gorm.DB) (map[string]model.UnitPrice, error) {
	items := []model.UnitPrice{
		// LABOR (Tenaga Kerja)
		{Code: "TK-01", Category: domain.CategoryLabor, WorkItem: "Pekerja", Unit: "OH", UnitPrice: 120000, Note: "Tenaga kerja harian standar"},
		{Code: "TK-02", Category: domain.CategoryLabor, WorkItem: "Tukang Batu", Unit: "OH", UnitPrice: 150000, Note: "Tukang batu & dinding"},
		{Code: "TK-03", Category: domain.CategoryLabor, WorkItem: "Tukang Besi / Kayu", Unit: "OH", UnitPrice: 150000, Note: "Pekerjaan fabrikasi & bekisting"},
		{Code: "TK-04", Category: domain.CategoryLabor, WorkItem: "Tukang Cat", Unit: "OH", UnitPrice: 145000, Note: "Pekerjaan finishing pengecatan"},
		{Code: "TK-05", Category: domain.CategoryLabor, WorkItem: "Kepala Tukang", Unit: "OH", UnitPrice: 175000, Note: "Koordinator tukang lapangan"},
		{Code: "TK-06", Category: domain.CategoryLabor, WorkItem: "Mandor", Unit: "OH", UnitPrice: 200000, Note: "Supervisi harian proyek"},

		// MATERIALS (Bahan Material)
		{Code: "MT-01", Category: domain.CategoryMaterials, WorkItem: "Semen Portland (50 kg)", Unit: "sak", UnitPrice: 75000, Note: "Standar SNI (Gresik/Tiga Roda)"},
		{Code: "MT-02", Category: domain.CategoryMaterials, WorkItem: "Pasir Pasang / Pasir Beton", Unit: "m3", UnitPrice: 320000, Note: "Pasir ayak kualitas baik"},
		{Code: "MT-03", Category: domain.CategoryMaterials, WorkItem: "Kerikil / Batu Split Beton", Unit: "m3", UnitPrice: 350000, Note: "Ukuran 1-2 cm"},
		{Code: "MT-04", Category: domain.CategoryMaterials, WorkItem: "Bata Ringan (Hebel) Tebal 10cm", Unit: "m3", UnitPrice: 720000, Note: "Ukuran standar SNI"},
		{Code: "MT-05", Category: domain.CategoryMaterials, WorkItem: "Semen Mortar Instan Perekat Bata Ringan", Unit: "sak", UnitPrice: 95000, Note: "Kemasan 40 kg"},
		{Code: "MT-06", Category: domain.CategoryMaterials, WorkItem: "Besi Beton Polos / Ulir", Unit: "kg", UnitPrice: 14500, Note: "Krakatau Steel SNI"},
		{Code: "MT-07", Category: domain.CategoryMaterials, WorkItem: "Kawat Beton (Bendrat)", Unit: "kg", UnitPrice: 28000, Note: "Pengikat tulangan"},
		{Code: "MT-08", Category: domain.CategoryMaterials, WorkItem: "Cat Dasar Tembok Alkali Resisting", Unit: "kg", UnitPrice: 48000, Note: "Dulux / Mowilex dasar"},
		{Code: "MT-09", Category: domain.CategoryMaterials, WorkItem: "Cat Tembok Interior Emulsi", Unit: "kg", UnitPrice: 65000, Note: "Warna standar pabrik"},
		{Code: "MT-10", Category: domain.CategoryMaterials, WorkItem: "Plamir Tembok", Unit: "kg", UnitPrice: 22000, Note: "Persiapan permukaan dinding"},
		{Code: "MT-11", Category: domain.CategoryMaterials, WorkItem: "Keramik Lantai 40x40 cm Polished", Unit: "m2", UnitPrice: 90000, Note: "Roman / Platinum motif polos"},

		// EQUIPMENT (Peralatan)
		{Code: "EQ-01", Category: domain.CategoryEquipment, WorkItem: "Concrete Mixer (Molen) 0.35 m3", Unit: "hari", UnitPrice: 250000, Note: "Termasuk BBM dan operator"},
		{Code: "EQ-02", Category: domain.CategoryEquipment, WorkItem: "Scaffolding Besi Pipa", Unit: "set/hari", UnitPrice: 35000, Note: "Komplit main frame & cross brace"},
		{Code: "EQ-03", Category: domain.CategoryEquipment, WorkItem: "Jack Hammer Listrik", Unit: "hari", UnitPrice: 180000, Note: "Untuk pembongkaran dinding & beton"},
		{Code: "EQ-04", Category: domain.CategoryEquipment, WorkItem: "Dump Truck 5 m3", Unit: "hari", UnitPrice: 850000, Note: "Pembuangan puing & material sisa"},
	}

	resultMap := make(map[string]model.UnitPrice)
	for _, item := range items {
		var existing model.UnitPrice
		err := db.Where("lower(code) = lower(?) AND deleted_at IS NULL", item.Code).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&item).Error; err != nil {
				return nil, err
			}
			resultMap[item.Code] = item
			slog.Info("seeded unit price", "code", item.Code, "workItem", item.WorkItem)
		} else if err != nil {
			return nil, err
		} else {
			resultMap[item.Code] = existing
		}
	}

	return resultMap, nil
}

type ahspItemTemplate struct {
	Category    domain.UnitPriceCategory
	Sort        int
	Description string
	Code        string
	Unit        string
	Coefficient float64
	UpCode      string // references UnitPrice Code
}

type ahspTemplate struct {
	Code              string
	WorkItem          string
	Unit              string
	OverheadProfitPct float64
	Status            domain.AhspStatus
	Note              string
	Items             []ahspItemTemplate
}

func seedAhsps(db *gorm.DB, upMap map[string]model.UnitPrice) (map[string]model.Ahsp, error) {
	templates := []ahspTemplate{
		{
			Code:              "AHSP-BR-10",
			WorkItem:          "Pemasangan 1 m2 Dinding Bata Ringan Tebal 10 cm",
			Unit:              "m2",
			OverheadProfitPct: 10.0,
			Status:            domain.AhspStatusFinal,
			Note:              "Standar SNI Pekerjaan Dinding Pabrik",
			Items: []ahspItemTemplate{
				{Category: domain.CategoryLabor, Sort: 1, Description: "Pekerja", Code: "TK-01", Unit: "OH", Coefficient: 0.2000, UpCode: "TK-01"},
				{Category: domain.CategoryLabor, Sort: 2, Description: "Tukang Batu", Code: "TK-02", Unit: "OH", Coefficient: 0.1000, UpCode: "TK-02"},
				{Category: domain.CategoryLabor, Sort: 3, Description: "Kepala Tukang", Code: "TK-05", Unit: "OH", Coefficient: 0.0100, UpCode: "TK-05"},
				{Category: domain.CategoryLabor, Sort: 4, Description: "Mandor", Code: "TK-06", Unit: "OH", Coefficient: 0.0100, UpCode: "TK-06"},
				{Category: domain.CategoryMaterials, Sort: 5, Description: "Bata Ringan (Hebel) Tebal 10cm", Code: "MT-04", Unit: "m3", Coefficient: 0.1000, UpCode: "MT-04"},
				{Category: domain.CategoryMaterials, Sort: 6, Description: "Semen Mortar Instan Perekat Bata Ringan", Code: "MT-05", Unit: "sak", Coefficient: 0.1250, UpCode: "MT-05"},
			},
		},
		{
			Code:              "AHSP-CAT-INT",
			WorkItem:          "Pengecatan 1 m2 Dinding Interior Baru (1 Dasar, 2 Penutup)",
			Unit:              "m2",
			OverheadProfitPct: 10.0,
			Status:            domain.AhspStatusFinal,
			Note:              "Finishing dinding ruang produksi dan kantor",
			Items: []ahspItemTemplate{
				{Category: domain.CategoryLabor, Sort: 1, Description: "Pekerja", Code: "TK-01", Unit: "OH", Coefficient: 0.0200, UpCode: "TK-01"},
				{Category: domain.CategoryLabor, Sort: 2, Description: "Tukang Cat", Code: "TK-04", Unit: "OH", Coefficient: 0.0630, UpCode: "TK-04"},
				{Category: domain.CategoryLabor, Sort: 3, Description: "Kepala Tukang", Code: "TK-05", Unit: "OH", Coefficient: 0.0063, UpCode: "TK-05"},
				{Category: domain.CategoryLabor, Sort: 4, Description: "Mandor", Code: "TK-06", Unit: "OH", Coefficient: 0.0030, UpCode: "TK-06"},
				{Category: domain.CategoryMaterials, Sort: 5, Description: "Plamir Tembok", Code: "MT-10", Unit: "kg", Coefficient: 0.1000, UpCode: "MT-10"},
				{Category: domain.CategoryMaterials, Sort: 6, Description: "Cat Dasar Tembok Alkali Resisting", Code: "MT-08", Unit: "kg", Coefficient: 0.1000, UpCode: "MT-08"},
				{Category: domain.CategoryMaterials, Sort: 7, Description: "Cat Tembok Interior Emulsi", Code: "MT-09", Unit: "kg", Coefficient: 0.2600, UpCode: "MT-09"},
				{Category: domain.CategoryEquipment, Sort: 8, Description: "Scaffolding Besi Pipa", Code: "EQ-02", Unit: "set/hari", Coefficient: 0.0500, UpCode: "EQ-02"},
			},
		},
		{
			Code:              "AHSP-KR-40",
			WorkItem:          "Pemasangan 1 m2 Lantai Keramik 40x40 cm Polished",
			Unit:              "m2",
			OverheadProfitPct: 10.0,
			Status:            domain.AhspStatusFinal,
			Note:              "Lantai area koridor dan ruang kerja",
			Items: []ahspItemTemplate{
				{Category: domain.CategoryLabor, Sort: 1, Description: "Pekerja", Code: "TK-01", Unit: "OH", Coefficient: 0.2500, UpCode: "TK-01"},
				{Category: domain.CategoryLabor, Sort: 2, Description: "Tukang Batu", Code: "TK-02", Unit: "OH", Coefficient: 0.1200, UpCode: "TK-02"},
				{Category: domain.CategoryLabor, Sort: 3, Description: "Kepala Tukang", Code: "TK-05", Unit: "OH", Coefficient: 0.0120, UpCode: "TK-05"},
				{Category: domain.CategoryLabor, Sort: 4, Description: "Mandor", Code: "TK-06", Unit: "OH", Coefficient: 0.0125, UpCode: "TK-06"},
				{Category: domain.CategoryMaterials, Sort: 5, Description: "Keramik Lantai 40x40 cm Polished", Code: "MT-11", Unit: "m2", Coefficient: 1.0500, UpCode: "MT-11"},
				{Category: domain.CategoryMaterials, Sort: 6, Description: "Semen Portland (50 kg)", Code: "MT-01", Unit: "sak", Coefficient: 0.2000, UpCode: "MT-01"},
				{Category: domain.CategoryMaterials, Sort: 7, Description: "Pasir Pasang / Pasir Beton", Code: "MT-02", Unit: "m3", Coefficient: 0.0450, UpCode: "MT-02"},
			},
		},
		{
			Code:              "AHSP-BOK-01",
			WorkItem:          "Pembongkaran 1 m2 Dinding Bata / Partisi Lama",
			Unit:              "m2",
			OverheadProfitPct: 10.0,
			Status:            domain.AhspStatusFinal,
			Note:              "Pekerjaan persiapan renovasi dan pembuangan puing",
			Items: []ahspItemTemplate{
				{Category: domain.CategoryLabor, Sort: 1, Description: "Pekerja", Code: "TK-01", Unit: "OH", Coefficient: 0.2500, UpCode: "TK-01"},
				{Category: domain.CategoryLabor, Sort: 2, Description: "Mandor", Code: "TK-06", Unit: "OH", Coefficient: 0.0250, UpCode: "TK-06"},
				{Category: domain.CategoryEquipment, Sort: 3, Description: "Jack Hammer Listrik", Code: "EQ-03", Unit: "hari", Coefficient: 0.0500, UpCode: "EQ-03"},
			},
		},
	}

	resultMap := make(map[string]model.Ahsp)

	for _, tmpl := range templates {
		var existing model.Ahsp
		err := db.Where("lower(code) = lower(?) AND deleted_at IS NULL", tmpl.Code).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var totalLabor int64
			var totalMaterial int64
			var totalEquipment int64

			var items []model.AhspItem
			for _, it := range tmpl.Items {
				up, ok := upMap[it.UpCode]
				unitPriceVal := int64(0)
				var upID *uint
				if ok {
					unitPriceVal = up.UnitPrice
					upID = &up.ID
				}
				cost := int64(math.Round(it.Coefficient * float64(unitPriceVal)))

				switch it.Category {
				case domain.CategoryLabor:
					totalLabor += cost
				case domain.CategoryMaterials:
					totalMaterial += cost
				case domain.CategoryEquipment:
					totalEquipment += cost
				}

				items = append(items, model.AhspItem{
					Category:    it.Category,
					Sort:        it.Sort,
					Description: it.Description,
					Code:        it.Code,
					Unit:        it.Unit,
					Coefficient: it.Coefficient,
					UnitPrice:   unitPriceVal,
					Cost:        cost,
					UnitPriceID: upID,
				})
			}

			totalCost := totalLabor + totalMaterial + totalEquipment
			overheadCost := int64(math.Round((tmpl.OverheadProfitPct / 100.0) * float64(totalCost)))
			unitPriceOfWork := totalCost + overheadCost

			ahsp := model.Ahsp{
				Code:              tmpl.Code,
				WorkItem:          tmpl.WorkItem,
				Unit:              tmpl.Unit,
				OverheadProfitPct: tmpl.OverheadProfitPct,
				TotalLabor:        totalLabor,
				TotalMaterial:     totalMaterial,
				TotalEquipment:    totalEquipment,
				TotalCost:         totalCost,
				UnitPriceOfWork:   unitPriceOfWork,
				Status:            tmpl.Status,
				Note:              tmpl.Note,
				Items:             items,
			}

			if err := db.Create(&ahsp).Error; err != nil {
				return nil, err
			}
			resultMap[ahsp.Code] = ahsp
			slog.Info("seeded ahsp", "code", ahsp.Code, "workItem", ahsp.WorkItem, "unitPrice", ahsp.UnitPriceOfWork)
		} else if err != nil {
			return nil, err
		} else {
			resultMap[tmpl.Code] = existing
		}
	}

	return resultMap, nil
}

func seedProjects(db *gorm.DB, ahspMap map[string]model.Ahsp) error {
	now := time.Now()
	planDate1 := now.AddDate(0, -1, 0) // 1 month ago
	planDate2 := now.AddDate(0, 1, 15) // next month
	planDate3 := now.AddDate(0, 0, -10)

	projects := []struct {
		project   model.Project
		schedules []model.Schedule
		boqs      []struct {
			boq    model.BoqItem
			volume model.Volume
			rab    model.RabItem
		}
	}{
		{
			project: model.Project{
				Name:                 "Renovasi Area Produksi & Koridor Gedung A",
				Section:              string(domain.ProjectSectionGABuilding),
				Plant:                string(domain.ProjectPlantCikarang),
				ImplementationPlan:   &planDate1,
				Budget:               350000000,
				Note:                 "Peremajaan dinding partisi dan lantai koridor Plant Cikarang",
				Status:               domain.ProjectStatusOnGoing,
				ProgressPPT:          domain.ProgressDone,
				ProgressDrawing:      domain.ProgressDone,
				ProgressRAB:          domain.ProgressDone,
				ProgressVendor:       domain.ProgressDone,
				ProgressSPB:          domain.ProgressDone,
				ProgressPR:           domain.ProgressDone,
				ProgressRINGI:        domain.ProgressDone,
				ProgressPO:           domain.ProgressDone,
				ProgressConstruction: domain.ProgressPending,
				ProgressBAST:         domain.ProgressPending,
			},
			schedules: []model.Schedule{
				{
					Code:      "SCH-01",
					WorkItem:  "Pekerjaan Pembongkaran Dinding & Pembersihan",
					StartDate: now.AddDate(0, 0, -20),
					EndDate:   now.AddDate(0, 0, -14),
					Duration:  7,
					Weight:    15.0,
					Progress:  100.0,
					Status:    domain.ScheduleDone,
				},
				{
					Code:      "SCH-02",
					WorkItem:  "Pemasangan Dinding Bata Ringan Baru",
					StartDate: now.AddDate(0, 0, -13),
					EndDate:   now.AddDate(0, 0, -3),
					Duration:  11,
					Weight:    35.0,
					Progress:  85.0,
					Status:    domain.ScheduleOnGoing,
				},
				{
					Code:      "SCH-03",
					WorkItem:  "Pemasangan Lantai Keramik Koridor",
					StartDate: now.AddDate(0, 0, -2),
					EndDate:   now.AddDate(0, 0, 8),
					Duration:  11,
					Weight:    25.0,
					Progress:  20.0,
					Status:    domain.ScheduleOnGoing,
				},
				{
					Code:      "SCH-04",
					WorkItem:  "Pengecatan Dinding & Finishing Akhir",
					StartDate: now.AddDate(0, 0, 9),
					EndDate:   now.AddDate(0, 0, 16),
					Duration:  8,
					Weight:    20.0,
					Progress:  0.0,
					Status:    domain.SchedulePlanning,
				},
				{
					Code:      "SCH-05",
					WorkItem:  "Serah Terima Pertama (BAST-1) & General Cleaning",
					StartDate: now.AddDate(0, 0, 17),
					EndDate:   now.AddDate(0, 0, 18),
					Duration:  2,
					Weight:    5.0,
					Progress:  0.0,
					Status:    domain.SchedulePlanning,
				},
			},
			boqs: []struct {
				boq    model.BoqItem
				volume model.Volume
				rab    model.RabItem
			}{
				{
					boq: model.BoqItem{
						WorkItem:      "Pembongkaran 1 m2 Dinding Bata / Partisi Lama",
						Unit:          "m2",
						PlannedVolume: 120.0,
					},
					volume: model.Volume{
						Quantity:      120.0,
						UnitWeight:    1.0,
						TotalQuantity: 120.0,
						Description:   "Sisi Barat Koridor Utama",
					},
					rab: model.RabItem{
						WorkItem:  "Pembongkaran 1 m2 Dinding Bata / Partisi Lama",
						Unit:      "m2",
						Volume:    120.0,
						UnitPrice: 45000,
						TotalCost: 5400000,
					},
				},
				{
					boq: model.BoqItem{
						WorkItem:      "Pemasangan 1 m2 Dinding Bata Ringan Tebal 10 cm",
						Unit:          "m2",
						PlannedVolume: 240.0,
					},
					volume: model.Volume{
						Quantity:      240.0,
						UnitWeight:    1.0,
						TotalQuantity: 240.0,
						Description:   "Dinding Partisi Koridor Baru",
					},
					rab: model.RabItem{
						WorkItem:  "Pemasangan 1 m2 Dinding Bata Ringan Tebal 10 cm",
						Unit:      "m2",
						Volume:    240.0,
						UnitPrice: 185000,
						TotalCost: 44400000,
					},
				},
				{
					boq: model.BoqItem{
						WorkItem:      "Pemasangan 1 m2 Lantai Keramik 40x40 cm Polished",
						Unit:          "m2",
						PlannedVolume: 180.0,
					},
					volume: model.Volume{
						Quantity:      180.0,
						UnitWeight:    1.0,
						TotalQuantity: 180.0,
						Description:   "Lantai Koridor Pejalan Kaki",
					},
					rab: model.RabItem{
						WorkItem:  "Pemasangan 1 m2 Lantai Keramik 40x40 cm Polished",
						Unit:      "m2",
						Volume:    180.0,
						UnitPrice: 195000,
						TotalCost: 35100000,
					},
				},
				{
					boq: model.BoqItem{
						WorkItem:      "Pengecatan 1 m2 Dinding Interior Baru (1 Dasar, 2 Penutup)",
						Unit:          "m2",
						PlannedVolume: 480.0,
					},
					volume: model.Volume{
						Quantity:      240.0,
						UnitWeight:    2.0, // 2 sisi dinding
						TotalQuantity: 480.0,
						Description:   "Dinding Sisi Dalam dan Sisi Luar",
					},
					rab: model.RabItem{
						WorkItem:  "Pengecatan 1 m2 Dinding Interior Baru (1 Dasar, 2 Penutup)",
						Unit:      "m2",
						Volume:    480.0,
						UnitPrice: 38000,
						TotalCost: 18240000,
					},
				},
			},
		},
		{
			project: model.Project{
				Name:                 "Pembangunan Instalasi Pengolahan Air Limbah (IPAL)",
				Section:              string(domain.ProjectSectionEHS),
				Plant:                string(domain.ProjectPlantKarawang),
				ImplementationPlan:   &planDate2,
				Budget:               620000000,
				Note:                 "Peningkatan kapasitas IPAL dan saluran drainase terpadu Plant Karawang",
				Status:               domain.ProjectStatusPlanning,
				ProgressPPT:          domain.ProgressDone,
				ProgressDrawing:      domain.ProgressDone,
				ProgressRAB:          domain.ProgressDone,
				ProgressVendor:       domain.ProgressPending,
				ProgressSPB:          domain.ProgressPending,
				ProgressPR:           domain.ProgressPending,
				ProgressRINGI:        domain.ProgressPending,
				ProgressPO:           domain.ProgressPending,
				ProgressConstruction: domain.ProgressPending,
				ProgressBAST:         domain.ProgressPending,
			},
		},
		{
			project: model.Project{
				Name:                 "Renovasi Ruang Istirahat & Kantin Karyawan",
				Section:              string(domain.ProjectSectionWellfare),
				Plant:                string(domain.ProjectPlantCikarang),
				ImplementationPlan:   &planDate3,
				Budget:               180000000,
				Note:                 "Fasilitas kesejahteraan karyawan shift pabrik",
				Status:               domain.ProjectStatusOnGoing,
				ProgressPPT:          domain.ProgressDone,
				ProgressDrawing:      domain.ProgressDone,
				ProgressRAB:          domain.ProgressDone,
				ProgressVendor:       domain.ProgressDone,
				ProgressSPB:          domain.ProgressPending,
				ProgressPR:           domain.ProgressPending,
				ProgressRINGI:        domain.ProgressPending,
				ProgressPO:           domain.ProgressPending,
				ProgressConstruction: domain.ProgressPending,
				ProgressBAST:         domain.ProgressPending,
			},
		},
	}

	for _, entry := range projects {
		var existing model.Project
		err := db.Where("lower(name) = lower(?) AND deleted_at IS NULL", entry.project.Name).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p := entry.project
			if err := db.Create(&p).Error; err != nil {
				return err
			}
			slog.Info("seeded project", "name", p.Name, "id", p.ID)

			// Seed schedules
			for _, s := range entry.schedules {
				s.ProjectID = p.ID
				if err := db.Create(&s).Error; err != nil {
					return err
				}
			}

			// Seed BoQ, Volume, RAB
			for _, b := range entry.boqs {
				boq := b.boq
				boq.ProjectID = p.ID

				// Link AHSP if matched
				for _, ah := range ahspMap {
					if ah.WorkItem == boq.WorkItem {
						boq.AhspID = &ah.ID
						break
					}
				}

				if err := db.Create(&boq).Error; err != nil {
					return err
				}

				// Volume
				vol := b.volume
				vol.BoqItemID = boq.ID
				d := now.AddDate(0, 0, -10)
				vol.Date = &d
				if err := db.Create(&vol).Error; err != nil {
					return err
				}

				// RAB
				rab := b.rab
				rab.ProjectID = p.ID
				rab.BoqItemID = &boq.ID
				if err := db.Create(&rab).Error; err != nil {
					return err
				}
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}
