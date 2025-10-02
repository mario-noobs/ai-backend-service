package face

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarioFaceServiceConfig_LoadMarioFaceServiceConfig(t *testing.T) {
	// Test with default values (no environment variables set)
	t.Run("default configuration", func(t *testing.T) {
		config := MarioFaceServiceConfig{}
		result := config.LoadMarioFaceServiceConfig()
		
		assert.Equal(t, "get-list", result.listPath)
		assert.Equal(t, "register-identity", result.enrollPath)
		assert.Equal(t, "recognize-identity", result.recognizePath)
		assert.Equal(t, "delete-identity", result.deletePath)
		assert.Equal(t, "is-registered", result.registeredPath)
		assert.Equal(t, "http://face-reg-engine:8080/face/v1/api/", result.Host)
		assert.Equal(t, "info", result.LogLevel)
	})

	// Test with custom environment variables
	t.Run("custom configuration from environment", func(t *testing.T) {
		// Set environment variables
		os.Setenv("FACE_LIST_API", "custom-list")
		os.Setenv("FACE_ENROLL_API", "custom-enroll")
		os.Setenv("FACE_RECOGNIZE_API", "custom-recognize")
		os.Setenv("FACE_DELETE_API", "custom-delete")
		os.Setenv("IS_REGISTERED_API", "custom-registered")
		os.Setenv("FACE_HOST", "http://custom-host:9090/api/")
		os.Setenv("FACE_SERVICE_LOG_LEVEL", "debug")
		
		// Clean up after test
		defer func() {
			os.Unsetenv("FACE_LIST_API")
			os.Unsetenv("FACE_ENROLL_API")
			os.Unsetenv("FACE_RECOGNIZE_API")
			os.Unsetenv("FACE_DELETE_API")
			os.Unsetenv("IS_REGISTERED_API")
			os.Unsetenv("FACE_HOST")
			os.Unsetenv("FACE_SERVICE_LOG_LEVEL")
		}()
		
		config := MarioFaceServiceConfig{}
		result := config.LoadMarioFaceServiceConfig()
		
		assert.Equal(t, "custom-list", result.listPath)
		assert.Equal(t, "custom-enroll", result.enrollPath)
		assert.Equal(t, "custom-recognize", result.recognizePath)
		assert.Equal(t, "custom-delete", result.deletePath)
		assert.Equal(t, "custom-registered", result.registeredPath)
		assert.Equal(t, "http://custom-host:9090/api/", result.Host)
		assert.Equal(t, "debug", result.LogLevel)
	})

	// Test partial environment variables (some set, some default)
	t.Run("partial configuration from environment", func(t *testing.T) {
		// Set only some environment variables
		os.Setenv("FACE_HOST", "http://partial-host:8080/")
		os.Setenv("FACE_SERVICE_LOG_LEVEL", "warn")
		
		// Clean up after test
		defer func() {
			os.Unsetenv("FACE_HOST")
			os.Unsetenv("FACE_SERVICE_LOG_LEVEL")
		}()
		
		config := MarioFaceServiceConfig{}
		result := config.LoadMarioFaceServiceConfig()
		
		// Custom values
		assert.Equal(t, "http://partial-host:8080/", result.Host)
		assert.Equal(t, "warn", result.LogLevel)
		
		// Default values
		assert.Equal(t, "get-list", result.listPath)
		assert.Equal(t, "register-identity", result.enrollPath)
		assert.Equal(t, "recognize-identity", result.recognizePath)
		assert.Equal(t, "delete-identity", result.deletePath)
		assert.Equal(t, "is-registered", result.registeredPath)
	})
}

func TestGetEnv(t *testing.T) {
	// Test when environment variable exists
	t.Run("environment variable exists", func(t *testing.T) {
		key := "TEST_ENV_VAR"
		expectedValue := "test-value"
		fallback := "fallback-value"
		
		os.Setenv(key, expectedValue)
		defer os.Unsetenv(key)
		
		result := getEnv(key, fallback)
		assert.Equal(t, expectedValue, result)
	})

	// Test when environment variable does not exist
	t.Run("environment variable does not exist", func(t *testing.T) {
		key := "NON_EXISTENT_ENV_VAR"
		fallback := "fallback-value"
		
		// Ensure the environment variable doesn't exist
		os.Unsetenv(key)
		
		result := getEnv(key, fallback)
		assert.Equal(t, fallback, result)
	})

	// Test with empty environment variable
	t.Run("environment variable is empty", func(t *testing.T) {
		key := "EMPTY_ENV_VAR"
		fallback := "fallback-value"
		
		os.Setenv(key, "")
		defer os.Unsetenv(key)
		
		result := getEnv(key, fallback)
		assert.Equal(t, "", result) // Empty string should be returned, not fallback
	})

	// Test with whitespace environment variable
	t.Run("environment variable has whitespace", func(t *testing.T) {
		key := "WHITESPACE_ENV_VAR"
		expectedValue := "  whitespace value  "
		fallback := "fallback-value"
		
		os.Setenv(key, expectedValue)
		defer os.Unsetenv(key)
		
		result := getEnv(key, fallback)
		assert.Equal(t, expectedValue, result) // Should preserve whitespace
	})
}
