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

	// Auth
	api.Post("/auth/login", auth.LoginHandler(db))
	api.Post("/auth/logout", auth.LogoutHandler())
	api.Get("/auth/me", auth.AuthMiddleware(), auth.MeHandler())

	// Protected routes
	authReq := auth.AuthMiddleware()
	adminOnly := auth.RequireRole(domain.RoleAdmin)

	// Projects
	api.Get("/projects", authReq, project.ListHandler(db))
	api.Post("/projects", authReq, adminOnly, project.CreateHandler(db))
	api.Get("/projects/:id", authReq, project.GetHandler(db))
	api.Put("/projects/:id", authReq, adminOnly, project.UpdateHandler(db))
	api.Delete("/projects/:id", authReq, adminOnly, project.DeleteHandler(db))

	// Schedules
	api.Get("/projects/:id/schedules", authReq, schedule.ListByProjectHandler(db))
	api.Post("/projects/:id/schedules", authReq, adminOnly, schedule.CreateHandler(db))
	api.Put("/schedules/:id", authReq, adminOnly, schedule.UpdateHandler(db))
	api.Delete("/schedules/:id", authReq, adminOnly, schedule.DeleteHandler(db))

	// Unit Prices
	api.Get("/unit-prices", authReq, unitprice.ListHandler(db))
	api.Post("/unit-prices", authReq, adminOnly, unitprice.CreateHandler(db))
	api.Put("/unit-prices/:id", authReq, adminOnly, unitprice.UpdateHandler(db))
	api.Delete("/unit-prices/:id", authReq, adminOnly, unitprice.DeleteHandler(db))

	// AHSP
	api.Get("/ahsps", authReq, ahsp.ListHandler(db))
	api.Post("/ahsps", authReq, adminOnly, ahsp.CreateHandler(db))
	api.Get("/ahsps/:id", authReq, ahsp.GetHandler(db))
	api.Put("/ahsps/:id", authReq, adminOnly, ahsp.UpdateHandler(db))
	api.Delete("/ahsps/:id", authReq, adminOnly, ahsp.DeleteHandler(db))

	// BOQ
	api.Get("/boq", authReq, boq.ListSummaryHandler(db))
	api.Get("/projects/:id/boq", authReq, boq.GetProjectBoqHandler(db))

	// Volume (take-off)
	api.Get("/projects/:id/volumes", authReq, volume.ListHandler(db))
	api.Post("/projects/:id/volumes", authReq, volume.CreateHandler(db))
	api.Put("/volumes/:id", authReq, volume.UpdateHandler(db))
	api.Delete("/volumes/:id", authReq, adminOnly, volume.DeleteHandler(db))

	// RAB
	api.Get("/projects/:id/rab", authReq, rab.GetProjectRabHandler(db))
	api.Post("/projects/:id/rab", authReq, adminOnly, rab.CreateHandler(db))
	api.Put("/rab-items/:id", authReq, adminOnly, rab.UpdateHandler(db))
	api.Delete("/rab-items/:id", authReq, adminOnly, rab.DeleteHandler(db))

	// Cost Estimate
	api.Get("/projects/:id/cost-estimate", authReq, costestimate.GetSummaryHandler(db))

	// Exports
	api.Get("/projects/:id/exports/volume.xlsx", authReq, export.VolumeExcelHandler(db))
	api.Get("/projects/:id/exports/rab.xlsx", authReq, export.RabExcelHandler(db))
	api.Get("/projects/:id/exports/rab.pdf", authReq, export.RabPdfHandler(db))

	return app
}

func health(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
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
