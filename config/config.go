package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	JWTSecret           string
	SquareAccessToken   string
	SquareApplicationID string
	SquareLocationID    string
	SquareEnvironment   string
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	config := &Config{
		Port:                getEnv("PORT", "8080"),
		JWTSecret:           getEnv("JWT_SECRET", "your-secret-key"),
		SquareAccessToken:   getEnv("SQUARE_ACCESS_TOKEN", "your-square-access-token"),
		SquareApplicationID: getEnv("SQUARE_APPLICATION_ID", "your-square-application-id"),
		SquareLocationID:    getEnv("SQUARE_LOCATION_ID", "your-square-location-id"),
		SquareEnvironment:   getEnv("SQUARE_ENVIRONMENT", "sandbox"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
