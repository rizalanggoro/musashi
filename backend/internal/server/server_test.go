package server

import (
	"testing"

	"musashi/backend/internal/platform/config"
)

func TestHealthRouteRegistered(t *testing.T) {
	app := New(nil, config.Config{AllowedOrigins: "http://localhost:5173"})

	hasRootHealth := false
	hasV1Health := false

	for _, r := range app.GetRoutes(true) {
		if r.Method == "GET" && r.Path == "/health" {
			hasRootHealth = true
		}
		if r.Method == "GET" && r.Path == "/api/v1/health" {
			hasV1Health = true
		}
	}

	if !hasRootHealth {
		t.Fatal("GET /health is not registered")
	}
	if !hasV1Health {
		t.Fatal("GET /api/v1/health is not registered")
	}
}
