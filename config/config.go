package config

import (
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    Port      string
    JWTSecret string
}

func LoadConfig() (*Config, error) {
    godotenv.Load()
    
    config := &Config{
        Port:      getEnv("PORT", "8080"),
        JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
    }
    
    return config, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}