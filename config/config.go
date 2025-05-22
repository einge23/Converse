package config

import "os"

type Config struct {
	Port            string
	ShutdownTimeout int
	Environment     string
	LogLevel        string
}

func New() *Config {
	cfg := &Config{
		Port:            os.Getenv("PORT"),
		ShutdownTimeout: 10,
		Environment:     os.Getenv("ENVIRONMENT"),
		LogLevel:        os.Getenv("LOG_LEVEL"),
	}

	return cfg
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
