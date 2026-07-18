package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv      string
	HTTPAddr    string
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
	JWTTTL      time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:      getenv("APP_ENV", "development"),
		HTTPAddr:    getenv("HTTP_ADDR", ":8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisURL:    getenv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		JWTTTL:      time.Duration(getenvInt("JWT_TTL_HOURS", 72)) * time.Hour,
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" || cfg.JWTSecret == "dev-only-change-me-ilo-super-secret-key" && cfg.AppEnv == "production" {
		if cfg.AppEnv == "production" {
			return Config{}, fmt.Errorf("JWT_SECRET must be set to a strong value in production")
		}
	}
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
