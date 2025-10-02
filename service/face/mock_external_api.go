package face

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"golang-ai-management/models"
)

// MockFaceRecognitionServer creates a test server that mocks the external face recognition API
func MockFaceRecognitionServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Route based on the endpoint path
		switch {
		case strings.Contains(r.URL.Path, "register-identity"):
			handleMockEnrollment(w, r)
		case strings.Contains(r.URL.Path, "recognize-identity"):
			handleMockRecognition(w, r)
		case strings.Contains(r.URL.Path, "delete-identity"):
			handleMockDeletion(w, r)
		case strings.Contains(r.URL.Path, "is-registered"):
			handleMockIsRegistered(w, r)
		case strings.Contains(r.URL.Path, "get-list"):
			handleMockGetList(w, r)
		default:
			http.NotFound(w, r)
		}
	}))
}

// handleMockEnrollment mocks the face enrollment API response
func handleMockEnrollment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var requestBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    "400",
			"message": "Invalid request body",
		})
		return
	}

	userId, exists := requestBody["userId"]
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    "400",
			"message": "userId is required",
		})
		return
	}

	// Simulate different responses based on userId
	userIdStr := userId.(string)
	switch userIdStr {
	case "test-error-user":
		// Simulate enrollment failure
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":      "400",
			"message":   "Face enrollment failed - poor image quality",
			"requestId": "req-error-001",
		})
	case "test-duplicate-user":
		// Simulate duplicate enrollment
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":      "409",
			"message":   "Face already enrolled for this user",
			"requestId": "req-duplicate-001",
		})
	default:
		// Successful enrollment
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":      models.Success,
			"message":   "Face enrolled successfully",
			"requestId": "req-enroll-" + userIdStr,
			"userId":    userIdStr,
			"rawImage":  "processed_base64_image_data",
			"flow":      "enroll",
		})
	}
}

// handleMockRecognition mocks the face recognition API response
func handleMockRecognition(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var requestBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    "400",
			"message": "Invalid request body",
		})
		return
	}

	imageBase64, exists := requestBody["imageBase64"]
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    "400",
			"message": "imageBase64 is required",
		})
		return
	}

	// Simulate different responses based on image content
	imageStr := imageBase64.(string)
	switch {
	case strings.Contains(imageStr, "unknown_face"):
		// No face found/recognized
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":      "404",
			"message":   "No matching face found",
			"requestId": "req-recognize-unknown",
			"searchData": map[string]interface{}{
				"searh_result": map[string]string{},
			},
		})
	case strings.Contains(imageStr, "low_confidence"):
		// Low confidence recognition
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":      models.Success,
			"message":   "Face recognized with low confidence",
			"requestId": "req-recognize-low",
			"userId":    "low_confidence_user",
			"rawImage":  "processed_base64_image_data",
			"flow":      "recognize",
			"searchData": map[string]interface{}{
				"searh_result": map[string]string{
					"low_confidence_user": "0.65",
				},
				"userId": "low_confidence_user",
			},
		})
	case strings.Contains(imageStr, "multiple_faces"):
		// Multiple faces detected
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":      models.Success,
			"message":   "Multiple faces recognized",
			"requestId": "req-recognize-multiple",
			"searchData": map[string]interface{}{
				"searh_result": map[string]string{
					"john_doe":  "0.95",
					"jane_doe":  "0.88",
					"bob_smith": "0.82",
				},
			},
		})
	default:
		// Successful high-confidence recognition
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":      models.Success,
			"message":   "Face recognized successfully",
			"requestId": "req-recognize-success",
			"userId":    "john_doe",
			"rawImage":  "processed_base64_image_data",
			"flow":      "recognize",
			"searchData": map[string]interface{}{
				"searh_result": map[string]string{
					"john_doe": "0.95",
				},
				"userId": "john_doe",
			},
		})
	}
}

// handleMockDeletion mocks the face deletion API response
func handleMockDeletion(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var requestBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    "400",
			"message": "Invalid request body",
		})
		return
	}

	userId, exists := requestBody["userId"]
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    "400",
			"message": "userId is required",
		})
		return
	}

	userIdStr := userId.(string)
	switch userIdStr {
	case "non-existent-user":
		// User not found
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":      "404",
			"message":   "User face identity not found",
			"requestId": "req-delete-notfound",
		})
	case "system-protected-user":
		// Protected user cannot be deleted
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":      "403",
			"message":   "Cannot delete system protected user",
			"requestId": "req-delete-forbidden",
		})
	default:
		// Successful deletion
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":      models.Success,
			"message":   "Face identity deleted successfully",
			"requestId": "req-delete-" + userIdStr,
			"userId":    userIdStr,
		})
	}
}

// handleMockIsRegistered mocks the is-registered check API response
func handleMockIsRegistered(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userId := r.URL.Query().Get("userId")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "userId parameter is required",
			"message": "Missing userId query parameter",
		})
		return
	}

	// Simulate different registration statuses
	var registered bool
	switch userId {
	case "registered-user", "john_doe", "jane_doe", "test-user":
		registered = true
	case "unregistered-user", "new-user", "unknown-user":
		registered = false
	default:
		// Default to registered for most test cases
		registered = true
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"registered": registered,
		"userId":     userId,
		"timestamp":  "2025-09-27T12:00:00Z",
	})
}

// handleMockGetList mocks the get-list API response
func handleMockGetList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Simulate pagination
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}

	// Mock face identities list
	faceIdentities := []map[string]interface{}{
		{
			"userId":     "john_doe",
			"enrolledAt": "2025-09-25T10:30:00Z",
			"confidence": 0.95,
			"status":     "active",
		},
		{
			"userId":     "jane_doe",
			"enrolledAt": "2025-09-24T14:20:00Z",
			"confidence": 0.92,
			"status":     "active",
		},
		{
			"userId":     "bob_smith",
			"enrolledAt": "2025-09-23T09:15:00Z",
			"confidence": 0.88,
			"status":     "active",
		},
		{
			"userId":     "alice_johnson",
			"enrolledAt": "2025-09-22T16:45:00Z",
			"confidence": 0.94,
			"status":     "inactive",
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    models.Success,
		"message": "Face identities retrieved successfully",
		"data": map[string]interface{}{
			"identities": faceIdentities,
			"pagination": map[string]interface{}{
				"page":       page,
				"limit":      limit,
				"total":      len(faceIdentities),
				"totalPages": 1,
			},
		},
		"timestamp": "2025-09-27T12:00:00Z",
	})
}

// MockAPIResponseScenarios provides predefined test scenarios
type MockAPIResponseScenarios struct {
	// Face enrollment scenarios
	SuccessfulEnrollment map[string]interface{}
	FailedEnrollment     map[string]interface{}
	DuplicateEnrollment  map[string]interface{}

	// Face recognition scenarios
	HighConfidenceRecognition map[string]interface{}
	LowConfidenceRecognition  map[string]interface{}
	NoFaceFound               map[string]interface{}
	MultipleFacesFound        map[string]interface{}

	// Face deletion scenarios
	SuccessfulDeletion map[string]interface{}
	UserNotFound       map[string]interface{}
	ProtectedUser      map[string]interface{}

	// Registration check scenarios
	UserRegistered    map[string]interface{}
	UserNotRegistered map[string]interface{}
}

// GetMockScenarios returns predefined mock API response scenarios
func GetMockScenarios() MockAPIResponseScenarios {
	return MockAPIResponseScenarios{
		// Enrollment scenarios
		SuccessfulEnrollment: map[string]interface{}{
			"code":      models.Success, // This is "0000"
			"message":   "Face enrolled successfully",
			"requestId": "req-enroll-success",
			"userId":    "test-user",
			"rawImage":  "base64_processed_image",
			"flow":      "enroll",
		},
		FailedEnrollment: map[string]interface{}{
			"code":      "400",
			"message":   "Face enrollment failed - poor image quality",
			"requestId": "req-enroll-failed",
		},
		DuplicateEnrollment: map[string]interface{}{
			"code":      "409",
			"message":   "Face already enrolled for this user",
			"requestId": "req-enroll-duplicate",
		},

		// Recognition scenarios
		HighConfidenceRecognition: map[string]interface{}{
			"code":      models.Success, // This is "0000"
			"message":   "Face recognized successfully",
			"requestId": "req-recognize-success",
			"userId":    "john_doe",
			"rawImage":  "base64_processed_image",
			"flow":      "recognize",
			"searchData": map[string]interface{}{
				"searh_result": map[string]string{
					"john_doe": "0.95",
				},
				"userId": "john_doe",
			},
		},
		LowConfidenceRecognition: map[string]interface{}{
			"code":      models.Success, // This is "0000"
			"message":   "Face recognized with low confidence",
			"requestId": "req-recognize-low",
			"userId":    "uncertain_user",
			"searchData": map[string]interface{}{
				"searh_result": map[string]string{
					"uncertain_user": "0.65",
				},
			},
		},
		NoFaceFound: map[string]interface{}{
			"code":      "404",
			"message":   "No matching face found",
			"requestId": "req-recognize-notfound",
			"searchData": map[string]interface{}{
				"searh_result": map[string]string{},
			},
		},
		MultipleFacesFound: map[string]interface{}{
			"code":      models.Success, // This is "0000"
			"message":   "Multiple faces recognized",
			"requestId": "req-recognize-multiple",
			"searchData": map[string]interface{}{
				"searh_result": map[string]string{
					"john_doe":  "0.95",
					"jane_doe":  "0.88",
					"bob_smith": "0.82",
				},
			},
		},

		// Deletion scenarios
		SuccessfulDeletion: map[string]interface{}{
			"code":      models.Success, // This is "0000"
			"message":   "Face identity deleted successfully",
			"requestId": "req-delete-success",
			"userId":    "test-user",
		},
		UserNotFound: map[string]interface{}{
			"code":      "404",
			"message":   "User face identity not found",
			"requestId": "req-delete-notfound",
		},
		ProtectedUser: map[string]interface{}{
			"code":      "403",
			"message":   "Cannot delete system protected user",
			"requestId": "req-delete-forbidden",
		},

		// Registration check scenarios
		UserRegistered: map[string]interface{}{
			"registered": true,
			"userId":     "test-user",
			"timestamp":  "2025-09-27T12:00:00Z",
		},
		UserNotRegistered: map[string]interface{}{
			"registered": false,
			"userId":     "new-user",
			"timestamp":  "2025-09-27T12:00:00Z",
		},
	}
}
