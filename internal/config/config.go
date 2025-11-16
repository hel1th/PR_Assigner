package config

import (
	"os"
)

type Config struct {
	ServerAddress string
	DatabaseURL   string
}

func Load() Config {
	return Config{
		ServerAddress: getEnv("SERVER_ADDRESS", "0.0.0.0:8080"),
		DatabaseURL: getEnv(
			"DATABASE_URL",
			"postgres://user:password@db:5432/pr_assigner?sslmode=disable",
		),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
