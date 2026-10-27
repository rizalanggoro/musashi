package server

import (
	"testing"

	"musashi/backend/internal/platform/config"
)

func TestHealthRouteRegistered(t *testing.T) {
	app := New(nil, config.Config{AllowedOrigins: "http://localhost:5173"})

	for _, r := range app.GetRoutes(true) {
		if r.Method == "GET" && r.Path == "/api/v1/health" {
			return
		}
	}
	t.Fatal("GET /api/v1/health is not registered")
}
