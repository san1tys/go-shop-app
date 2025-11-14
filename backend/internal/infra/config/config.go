package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerPort string
	DBDSN      string
	JWTSecret  string
}

// getEnv возвращает значение переменной окружения или дефолт
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// Load загружает конфиг из переменных окружения
func Load() (*Config, error) {
	cfg := &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBDSN:      os.Getenv("DB_DSN"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}

	if cfg.DBDSN == "" {
		return nil, fmt.Errorf("DB_DSN is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}
