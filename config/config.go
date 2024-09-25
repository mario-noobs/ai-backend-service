package config

import (
	"os"
)

type ServerConfig struct {
	Port      string
	Host      string
	LogLevel  string
	LogFormat string
	LogType   string
}

var Config = loadConfig()

func loadConfig() ServerConfig {
	port := getEnv("SERVER_PORT", "8080")
	host := getEnv("SERVER_HOST", "localhost")
	logLevel := getEnv("LOG_LEVEL", "debug")
	LogFormat := getEnv("LOG_FORMAT", "json")
	LogType := getEnv("LOG_TYPE", "logrus")

	return ServerConfig{
		Port:      port,
		Host:      host,
		LogLevel:  logLevel,
		LogFormat: LogFormat,
		LogType:   LogType,
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
