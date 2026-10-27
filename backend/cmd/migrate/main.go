// Command migrate applies embedded SQL migrations. Run it before starting the API.
package main

import (
	"log/slog"
	"os"

	"musashi/backend/internal/platform/config"
	"musashi/backend/internal/platform/database"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		slog.Error("open database", "err", err)
		os.Exit(1)
	}

	if err := database.Migrate(db); err != nil {
		slog.Error("migrate", "err", err)
		os.Exit(1)
	}

	slog.Info("migrations applied")
}
