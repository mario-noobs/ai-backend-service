package models

import (
	"os"
)

type ServerConfig struct {
	Port     string
	Host     string
	LogLevel string
}

var Config = loadConfig()

func loadConfig() ServerConfig {
	port := getEnv("SERVER_PORT", "8080")
	host := getEnv("SERVER_HOST", "localhost")
	logLevel := getEnv("LOG_LEVEL", "info")

	return ServerConfig{
		Port:     port,
		Host:     host,
		LogLevel: logLevel,
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
