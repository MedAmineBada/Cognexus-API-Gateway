package config

import (
	"os"
)

type Config struct {
	JWT_SECRET        string
	PORT              string
	EXAM_SERVICE_URL  string
	AUTH_SERVICE_URL  string
	CLASS_SERVICE_URL string
	REDIS_URL         string
}

func Load() *Config {
	return &Config{
		JWT_SECRET:       os.Getenv("JWT_SECRET"),
		PORT:             os.Getenv("PORT"),
		EXAM_SERVICE_URL: os.Getenv("EXAM_SERVICE_URL"),
		AUTH_SERVICE_URL: os.Getenv("AUTH_SERVICE_URL"),
		REDIS_URL:        os.Getenv("REDIS_URL"),
	}
}
