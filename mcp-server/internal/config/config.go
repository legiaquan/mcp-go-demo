package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
}

func Load() *Config {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Fallback for docker environment since inspector proxy strips env vars
		connStr = "postgres://demo:demo123@db:5432/demodb?sslmode=disable"
	}

	return &Config{
		DatabaseURL: connStr,
	}
}
