package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config содержит все настройки приложения
type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Опционально: параметры подключения (например, таймауты)
	DBTimeout time.Duration
}

// Load загружает конфигурацию из .env или переменных окружения
func Load() (*Config, error) {
	// Загружаем .env, если файл существует
	if err := godotenv.Load(); err != nil {
		// .env не обязателен — просто используем переменные окружения
	}

	getEnv := func(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}

	cfg := &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "subscriptions_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		DBTimeout:  10 * time.Second,
	}

	// Валидация обязательных полей
	if cfg.DBUser == "" {
		return nil, fmt.Errorf("не задан DB_USER")
	}
	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("не задан DB_PASSWORD")
	}
	if cfg.DBName == "" {
		return nil, fmt.Errorf("не задан DB_NAME")
	}

	return cfg, nil
}
