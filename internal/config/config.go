package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	ShutdownTimeout time.Duration
	Environment     string
	LogLevel        string
	DatabaseURL     string
	BucketName      string
	Region          string
}

func New() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	defaultPort := "8080"
	
	return &Config{
		Port:            getEnv("PORT", defaultPort),
		ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
		Environment:     getEnv("ENVIRONMENT", "development"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		BucketName:      getEnv("CONVERSE_AWS_BUCKET", ""),
		Region: 		 getEnv("AWS_REGION", "us-east-2"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		duration, err := time.ParseDuration(value)
		if err != nil {
			log.Printf("Warning: invalid duration format for %s: %s", key, value)
			return defaultValue
		}
		return duration
	}
	return defaultValue
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
