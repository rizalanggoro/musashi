// Command seed populates development seed data into the database.
package main

import (
	"log/slog"
	"os"

	"musashi/backend/internal/platform/config"
	"musashi/backend/internal/platform/database"
	"musashi/backend/internal/platform/database/seeder"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		slog.Error("open database", "err", err)
		os.Exit(1)
	}

	if err := seeder.Seed(db); err != nil {
		slog.Error("seed database", "err", err)
		os.Exit(1)
	}

	slog.Info("database seeded successfully")
}

