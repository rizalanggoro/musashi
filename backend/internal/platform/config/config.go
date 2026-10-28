package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	Port           string
	DatabaseURL    string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	AllowedOrigins string
	AutoMigrate    bool
}

func loadDotEnv(files ...string) {
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				val = strings.Trim(val, `"'`)
				if os.Getenv(key) == "" {
					os.Setenv(key, val)
				}
			}
		}
	}
}

func Load() Config {
	loadDotEnv(".env", "backend/.env")

	host := env("DB_HOST", "localhost")
	port := env("DB_PORT", "5432")
	user := env("DB_USER", "postgres")
	password := env("DB_PASSWORD", "postgres")
	dbName := env("DB_NAME", "musashi")
	sslMode := env("DB_SSLMODE", "disable")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		if password != "" {
			dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, url.QueryEscape(password), host, port, dbName, sslMode)
		} else {
			dbURL = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s", user, host, port, dbName, sslMode)
		}
	}

	return Config{
		Port:           env("PORT", "8080"),
		DatabaseURL:    dbURL,
		DBHost:         host,
		DBPort:         port,
		DBUser:         user,
		DBPassword:     password,
		DBName:         dbName,
		DBSSLMode:      sslMode,
		AllowedOrigins: env("ALLOWED_ORIGINS", "http://localhost:5173"),
		AutoMigrate:    env("AUTO_MIGRATE", "true") == "true",
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
