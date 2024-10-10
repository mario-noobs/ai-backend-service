package service

import (
	"os"
)

type MarioFaceServiceConfig struct {
	Host          string
	LogLevel      string
	listPath      string
	enrollPath    string
	recognizePath string
}

func (m MarioFaceServiceConfig) LoadMarioFaceServiceConfig() MarioFaceServiceConfig {
	return MarioFaceServiceConfig{
		listPath:      getEnv("FACE_LIST_API", "get-list"),
		enrollPath:    getEnv("FACE_ENROLL_API", "add-identity"),
		recognizePath: getEnv("FACE_RECOGNIZE_API", "recognize"),
		Host:          getEnv("FACE_HOST", "http://75.119.149.223:5000/face/"),
		LogLevel:      getEnv("FACE_SERVICE_LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
