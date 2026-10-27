package config

import "os"

type Config struct {
	Port           string
	DatabaseURL    string
	AllowedOrigins string
}

func Load() Config {
	return Config{
		Port:           env("PORT", "8080"),
		DatabaseURL:    env("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/musashi?sslmode=disable"),
		AllowedOrigins: env("ALLOWED_ORIGINS", "http://localhost:5173"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
