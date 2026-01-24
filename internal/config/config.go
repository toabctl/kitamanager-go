package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	ServerPort    string
	RBACModelPath string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "kitamanager"),
		DBPassword:    getEnv("DB_PASSWORD", "kitamanager"),
		DBName:        getEnv("DB_NAME", "kitamanager"),
		JWTSecret:     getEnv("JWT_SECRET", "default-secret-key"),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		RBACModelPath: getEnv("RBAC_MODEL_PATH", "configs/rbac_model.conf"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
