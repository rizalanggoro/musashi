package server

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"gorm.io/gorm"

	"musashi/backend/internal/ahsp"
	"musashi/backend/internal/auth"
	"musashi/backend/internal/boq"
	"musashi/backend/internal/costestimate"
	"musashi/backend/internal/domain"
	"musashi/backend/internal/export"
	"musashi/backend/internal/platform/config"
	"musashi/backend/internal/project"
	"musashi/backend/internal/rab"
	"musashi/backend/internal/schedule"
	"musashi/backend/internal/unitprice"
	"musashi/backend/internal/volume"
	"musashi/backend/utils"
)

func New(db *gorm.DB, cfg config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "musashi",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	app.Use(requestid.New(), recover.New(), logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.AllowedOrigins},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Root health check
	app.Get("/health", health(db))

	api := app.Group("/api/v1")

	// System & Health
	api.Get("/health", health(db))

	// Feature handlers (Repository -> Service -> Handler)
	authSvc := auth.NewService(auth.NewRepository(db))
	authHandler := auth.NewHandler(authSvc)
	projectHandler := project.NewHandler(project.NewService(project.NewRepository(db)))
	scheduleHandler := schedule.NewHandler(schedule.NewService(schedule.NewRepository(db)))
	unitPriceHandler := unitprice.NewHandler(unitprice.NewService(unitprice.NewRepository(db)))
	ahspHandler := ahsp.NewHandler(ahsp.NewService(ahsp.NewRepository(db)))
	boqHandler := boq.NewHandler(boq.NewService(boq.NewRepository(db)))
	volumeHandler := volume.NewHandler(volume.NewService(volume.NewRepository(db)))
	rabHandler := rab.NewHandler(rab.NewService(rab.NewRepository(db)))
	costEstimateHandler := costestimate.NewHandler(costestimate.NewService(costestimate.NewRepository(db)))
	exportHandler := export.NewHandler(export.NewService(export.NewRepository(db)))

	// Auth
	api.Post("/auth/login", authHandler.Login)
	api.Post("/auth/logout", authHandler.Logout)
	api.Get("/auth/me", utils.AuthMiddleware(authSvc), authHandler.Me)

	// Protected routes
	authReq := utils.AuthMiddleware(authSvc)
	adminOnly := utils.RequireRole(domain.RoleAdmin)

	// Projects
	api.Get("/projects", authReq, projectHandler.List)
	api.Post("/projects", authReq, adminOnly, projectHandler.Create)
	api.Get("/projects/:id", authReq, projectHandler.Get)
	api.Put("/projects/:id", authReq, adminOnly, projectHandler.Update)
	api.Delete("/projects/:id", authReq, adminOnly, projectHandler.Delete)

	// Schedules
	api.Get("/projects/:id/schedules", authReq, scheduleHandler.ListByProject)
	api.Post("/projects/:id/schedules", authReq, adminOnly, scheduleHandler.Create)
	api.Put("/schedules/:id", authReq, adminOnly, scheduleHandler.Update)
	api.Delete("/schedules/:id", authReq, adminOnly, scheduleHandler.Delete)

	// Unit Prices
	api.Get("/unit-prices", authReq, unitPriceHandler.List)
	api.Post("/unit-prices", authReq, adminOnly, unitPriceHandler.Create)
	api.Put("/unit-prices/:id", authReq, adminOnly, unitPriceHandler.Update)
	api.Delete("/unit-prices/:id", authReq, adminOnly, unitPriceHandler.Delete)

	// AHSP
	api.Get("/ahsps", authReq, ahspHandler.List)
	api.Post("/ahsps", authReq, adminOnly, ahspHandler.Create)
	api.Get("/ahsps/:id", authReq, ahspHandler.Get)
	api.Put("/ahsps/:id", authReq, adminOnly, ahspHandler.Update)
	api.Delete("/ahsps/:id", authReq, adminOnly, ahspHandler.Delete)

	// BOQ
	api.Get("/boq", authReq, boqHandler.ListSummary)
	api.Get("/projects/:id/boq", authReq, boqHandler.GetProjectBoq)

	// Volume (take-off)
	api.Get("/projects/:id/volumes", authReq, volumeHandler.List)
	api.Post("/projects/:id/volumes", authReq, volumeHandler.Create)
	api.Put("/volumes/:id", authReq, volumeHandler.Update)
	api.Delete("/volumes/:id", authReq, adminOnly, volumeHandler.Delete)

	// RAB
	api.Get("/projects/:id/rab", authReq, rabHandler.GetProjectRab)
	api.Post("/projects/:id/rab", authReq, adminOnly, rabHandler.Create)
	api.Put("/rab-items/:id", authReq, adminOnly, rabHandler.Update)
	api.Delete("/rab-items/:id", authReq, adminOnly, rabHandler.Delete)

	// Cost Estimate
	api.Get("/projects/:id/cost-estimate", authReq, costEstimateHandler.GetSummary)

	// Exports
	api.Get("/projects/:id/exports/volume.xlsx", authReq, exportHandler.VolumeExcel)
	api.Get("/projects/:id/exports/rab.xlsx", authReq, exportHandler.RabExcel)
	api.Get("/projects/:id/exports/rab.pdf", authReq, exportHandler.RabPdf)

	return app
}

func health(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		if db == nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":   "degraded",
				"database": "down",
			})
		}
		sqlDB, err := db.DB()
		if err != nil || sqlDB.PingContext(c.Context()) != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":   "degraded",
				"database": "down",
			})
		}

		return c.JSON(fiber.Map{
			"status":   "ok",
			"database": "up",
			"time":     time.Now().UTC().Format(time.RFC3339),
		})
	}
}
