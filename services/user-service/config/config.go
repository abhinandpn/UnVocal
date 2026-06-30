package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	cfg := Config{
		DatabaseURL: os.Getenv("DB_URL"),
		Port:        os.Getenv("PORT"),
	}

	if cfg.DatabaseURL == "" {
		panic("DB_URL is required")
	}

	if cfg.Port == "" {
		panic("PORT is required")
	}

	return cfg
}
