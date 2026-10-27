package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"

	"musashi/backend/internal/platform/config"
	"musashi/backend/internal/platform/database"
	"musashi/backend/internal/server"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		slog.Error("open database", "err", err)
		os.Exit(1)
	}

	if cfg.AutoMigrate {
		if err := database.Migrate(db); err != nil {
			slog.Error("auto-migrate database", "err", err)
			os.Exit(1)
		}
		slog.Info("database migrations applied")
	}

	app := server.New(db, cfg)

	go func() {
		slog.Info("api listening", "addr", ":"+cfg.Port)
		if err := app.Listen(":"+cfg.Port, fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
			slog.Error("serve", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		slog.Error("shutdown", "err", err)
	}
}
