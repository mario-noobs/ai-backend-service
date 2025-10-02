package face

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMockFaceRecognitionServer_Enrollment tests the mock enrollment endpoint
func TestMockFaceRecognitionServer_Enrollment(t *testing.T) {
	server := MockFaceRecognitionServer()
	defer server.Close()

	scenarios := GetMockScenarios()

	t.Run("successful enrollment", func(t *testing.T) {
		// Test successful enrollment response
		response := scenarios.SuccessfulEnrollment
		assert.Equal(t, "0000", response["code"]) // Changed from "200" to "0000"
		assert.Equal(t, "Face enrolled successfully", response["message"])
		assert.NotNil(t, response["userId"])
	})

	t.Run("failed enrollment", func(t *testing.T) {
		// Test failed enrollment response
		response := scenarios.FailedEnrollment
		assert.Equal(t, "400", response["code"])
		assert.Contains(t, response["message"], "Face enrollment failed")
	})

	t.Run("duplicate enrollment", func(t *testing.T) {
		// Test duplicate enrollment response
		response := scenarios.DuplicateEnrollment
		assert.Equal(t, "409", response["code"])
		assert.Contains(t, response["message"], "already enrolled")
	})
}

// TestMockFaceRecognitionServer_Recognition tests the mock recognition endpoint
func TestMockFaceRecognitionServer_Recognition(t *testing.T) {
	scenarios := GetMockScenarios()

	t.Run("high confidence recognition", func(t *testing.T) {
		response := scenarios.HighConfidenceRecognition
		assert.Equal(t, "0000", response["code"]) // Changed from "200" to "0000"
		assert.Equal(t, "Face recognized successfully", response["message"])

		searchData := response["searchData"].(map[string]interface{})
		searchResult := searchData["searh_result"].(map[string]string)
		assert.Equal(t, "0.95", searchResult["john_doe"])
	})

	t.Run("low confidence recognition", func(t *testing.T) {
		response := scenarios.LowConfidenceRecognition
		assert.Equal(t, "0000", response["code"]) // Changed from "200" to "0000"
		assert.Contains(t, response["message"], "low confidence")
	})

	t.Run("no face found", func(t *testing.T) {
		response := scenarios.NoFaceFound
		assert.Equal(t, "404", response["code"])
		assert.Contains(t, response["message"], "No matching face")
	})

	t.Run("multiple faces found", func(t *testing.T) {
		response := scenarios.MultipleFacesFound
		assert.Equal(t, "0000", response["code"]) // Changed from "200" to "0000"

		searchData := response["searchData"].(map[string]interface{})
		searchResult := searchData["searh_result"].(map[string]string)
		assert.True(t, len(searchResult) > 1)
		assert.Equal(t, "0.95", searchResult["john_doe"])
		assert.Equal(t, "0.88", searchResult["jane_doe"])
	})
}

// TestMockFaceRecognitionServer_Deletion tests the mock deletion endpoint
func TestMockFaceRecognitionServer_Deletion(t *testing.T) {
	scenarios := GetMockScenarios()

	t.Run("successful deletion", func(t *testing.T) {
		response := scenarios.SuccessfulDeletion
		assert.Equal(t, "0000", response["code"]) // Changed from "200" to "0000"
		assert.Contains(t, response["message"], "deleted successfully")
	})

	t.Run("user not found", func(t *testing.T) {
		response := scenarios.UserNotFound
		assert.Equal(t, "404", response["code"])
		assert.Contains(t, response["message"], "not found")
	})

	t.Run("protected user", func(t *testing.T) {
		response := scenarios.ProtectedUser
		assert.Equal(t, "403", response["code"])
		assert.Contains(t, response["message"], "protected user")
	})
}

// TestMockFaceRecognitionServer_Registration tests the mock registration check endpoint
func TestMockFaceRecognitionServer_Registration(t *testing.T) {
	scenarios := GetMockScenarios()

	t.Run("user registered", func(t *testing.T) {
		response := scenarios.UserRegistered
		assert.True(t, response["registered"].(bool))
		assert.NotNil(t, response["timestamp"])
	})

	t.Run("user not registered", func(t *testing.T) {
		response := scenarios.UserNotRegistered
		assert.False(t, response["registered"].(bool))
		assert.NotNil(t, response["timestamp"])
	})
}

//// Example of how to use the mock server in your existing tests
//func ExampleUsageInTests(t *testing.T) {
//	// Start the mock server
//	mockServer := MockFaceRecognitionServer()
//	defer mockServer.Close()
//
//	// Create a face business instance with the mock server URL
//	mockConfig := MarioFaceServiceConfig{
//		Host:           mockServer.URL,
//		enrollPath:     "/register-identity",
//		recognizePath:  "/recognize-identity",
//		deletePath:     "/delete-identity",
//		registeredPath: "/is-registered",
//		listPath:       "/get-list",
//	}
//
//	// Use this config in your FaceBusiness for testing
//	// business := &FaceBusiness{config: mockConfig, ...}
//
//	// Now all external API calls will hit your mock server
//	// providing predictable responses for testing
//}
