package config

import (
	"os"
	"time"

	"routrapp-api/internal/utils/constants"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", constants.DefaultPort),
			ReadTimeout:  constants.DefaultReadTimeout,
			WriteTimeout: constants.DefaultWriteTimeout,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 