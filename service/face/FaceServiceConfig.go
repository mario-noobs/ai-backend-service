package face

import (
	"os"
)

type MarioFaceServiceConfig struct {
	Host           string
	LogLevel       string
	listPath       string
	enrollPath     string
	recognizePath  string
	deletePath     string
	registeredPath string
}

func (m MarioFaceServiceConfig) LoadMarioFaceServiceConfig() MarioFaceServiceConfig {
	return MarioFaceServiceConfig{
		listPath:       getEnv("FACE_LIST_API", "get-list"),
		enrollPath:     getEnv("FACE_ENROLL_API", "register-identity"),
		recognizePath:  getEnv("FACE_RECOGNIZE_API", "recognize-identity"),
		deletePath:     getEnv("FACE_DELETE_API", "delete-identity"),
		registeredPath: getEnv("IS_REGISTERED_API", "is-registered"),
		Host:           getEnv("FACE_HOST", "http://face-reg-engine:8080/face/v1/api/"),
		LogLevel:       getEnv("FACE_SERVICE_LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
