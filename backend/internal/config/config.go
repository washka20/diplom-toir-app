package config

import (
	"fmt"
	"os"
)

// Config содержит конфигурацию приложения, загружаемую из переменных окружения.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	ServerPort string
	JWTSecret  string
}

// Load читает конфигурацию из переменных окружения.
// Возвращает ошибку, если обязательные переменные (DB_HOST, DB_NAME, JWT_SECRET) отсутствуют.
func Load() (*Config, error) {
	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     envOrDefault("DB_PORT", "5432"),
		DBUser:     envOrDefault("DB_USER", "postgres"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  envOrDefault("DB_SSLMODE", "disable"),
		ServerPort: envOrDefault("SERVER_PORT", "8080"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}

	if cfg.DBHost == "" {
		return nil, fmt.Errorf("required environment variable DB_HOST is not set")
	}
	if cfg.DBName == "" {
		return nil, fmt.Errorf("required environment variable DB_NAME is not set")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("required environment variable JWT_SECRET is not set")
	}

	return cfg, nil
}

// DSN возвращает строку подключения к PostgreSQL.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode,
	)
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
