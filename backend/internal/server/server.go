package server

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"gorm.io/gorm"

	"musashi/backend/internal/platform/config"
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

	api := app.Group("/api/v1")
	api.Get("/health", health(db))

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
