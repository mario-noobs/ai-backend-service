package face

import (
	"context"
	helper "golang-ai-management/helpers"
	"golang-ai-management/models"
	"golang-ai-management/models/response"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFaceService is a mock implementation of FaceService
type MockFaceService struct {
	mock.Mock
}

func (m *MockFaceService) Enroll(ctx context.Context, face models.Face, jwt string) response.FaceRegResponse {
	args := m.Called(ctx, face, jwt)
	return args.Get(0).(response.FaceRegResponse)
}

func (m *MockFaceService) Recognize(ctx context.Context, face models.Face, jwt string) response.FaceRegResponse {
	args := m.Called(ctx, face, jwt)
	return args.Get(0).(response.FaceRegResponse)
}

func (m *MockFaceService) Delete(ctx context.Context, face models.Face, jwt string) response.FaceRegResponse {
	args := m.Called(ctx, face, jwt)
	return args.Get(0).(response.FaceRegResponse)
}

// MockMarioFaceServiceConfig is a mock implementation that returns our test configuration
type MockMarioFaceServiceConfig struct {
	config MarioFaceServiceConfig
}

func (m *MockMarioFaceServiceConfig) LoadMarioFaceServiceConfig() MarioFaceServiceConfig {
	return m.config
}

// MockConfigProvider implements ConfigProvider for testing
type MockConfigProvider struct {
	config MarioFaceServiceConfig
}

func (m *MockConfigProvider) LoadMarioFaceServiceConfig() MarioFaceServiceConfig {
	return m.config
}

func TestNewFaceBusiness(t *testing.T) {
	mockFaceService := new(MockFaceService)
	mockConfigProvider := &MockConfigProvider{config: MarioFaceServiceConfig{}}
	mockTimer := helper.Timer{}

	business := NewFaceBusiness(mockFaceService, mockConfigProvider, mockTimer)

	assert.NotNil(t, business)
	assert.Equal(t, mockFaceService, business.FaceBussiness)
}

// Create a business instance that uses our mock config interface
type TestFaceBusiness struct {
	*FaceBussiness
	mockConfig *MockMarioFaceServiceConfig
}

func (t *TestFaceBusiness) GetConfig() MarioFaceServiceConfig {
	return t.mockConfig.LoadMarioFaceServiceConfig()
}

func TestFaceBusiness_Enroll_Success(t *testing.T) {
	// Use the new mock server instead of creating individual test servers
	mockServer := MockFaceRecognitionServer()
	defer mockServer.Close()

	// Setup with mock config provider that uses the comprehensive mock server
	mockConfigProvider := &MockConfigProvider{
		config: MarioFaceServiceConfig{
			Host:       mockServer.URL,
			enrollPath: "/register-identity",
		},
	}
	mockTimer := helper.Timer{}

	business := &FaceBussiness{
		config: mockConfigProvider, // Use ConfigProvider interface
		time:   mockTimer,
	}

	ctx := context.Background()
	face := models.Face{
		TransactionId: "test-transaction",
		Name:          stringPtr("test-user"),
		Image:         stringPtr("base64-image-data"),
	}
	jwt := "Bearer test-token"

	// Execute
	result := business.Enroll(ctx, face, jwt)

	// Verify using the mock scenarios
	assert.Equal(t, "test-user", result.UserId)
	assert.Equal(t, models.Success, result.BasicResponse.Code)
	assert.NotNil(t, result.Data.CreatedAt)
}

// New test for error scenarios using the mock system
func TestFaceBusiness_Enroll_ErrorScenarios(t *testing.T) {
	mockServer := MockFaceRecognitionServer()
	defer mockServer.Close()

	mockConfigProvider := &MockConfigProvider{
		config: MarioFaceServiceConfig{
			Host:       mockServer.URL,
			enrollPath: "/register-identity",
		},
	}
	mockTimer := helper.Timer{}

	business := &FaceBussiness{
		config: mockConfigProvider, // Use ConfigProvider interface
		time:   mockTimer,
	}

	ctx := context.Background()
	jwt := "Bearer test-token"

	t.Run("enrollment failure", func(t *testing.T) {
		face := models.Face{
			Name:  stringPtr("test-error-user"), // Triggers error response
			Image: stringPtr("base64-image-data"),
		}

		result := business.Enroll(ctx, face, jwt)
		assert.Equal(t, "400", result.BasicResponse.Code)
		assert.Contains(t, result.BasicResponse.Message, "enrollment failed")
	})

	t.Run("duplicate user", func(t *testing.T) {
		face := models.Face{
			Name:  stringPtr("test-duplicate-user"), // Triggers duplicate response
			Image: stringPtr("base64-image-data"),
		}

		result := business.Enroll(ctx, face, jwt)
		assert.Equal(t, "409", result.BasicResponse.Code)
		assert.Contains(t, result.BasicResponse.Message, "already enrolled")
	})
}

func TestFaceBusiness_Recognize_Success(t *testing.T) {
	// Use the new mock server instead of creating individual test servers
	mockServer := MockFaceRecognitionServer()
	defer mockServer.Close()

	// Setup with mock config provider that uses the comprehensive mock server
	mockConfigProvider := &MockConfigProvider{
		config: MarioFaceServiceConfig{
			Host:          mockServer.URL,
			recognizePath: "/recognize-identity",
		},
	}
	mockTimer := helper.Timer{}

	business := &FaceBussiness{
		config: mockConfigProvider, // Use ConfigProvider interface
		time:   mockTimer,
	}

	ctx := context.Background()
	face := models.Face{
		TransactionId: "test-transaction",
		Image:         stringPtr("base64-image-data"),
	}
	jwt := "Bearer test-token"

	// Execute
	result := business.Recognize(ctx, face, jwt)

	// Verify
	assert.Equal(t, models.Success, result.BasicResponse.Code)
	assert.Equal(t, "john_doe", result.UserId)
	assert.NotNil(t, result.Data.Name)
	assert.NotNil(t, result.Data.Probability)
}

func TestFaceBusiness_Delete_Success(t *testing.T) {
	// Use the new mock server instead of creating individual test servers
	mockServer := MockFaceRecognitionServer()
	defer mockServer.Close()

	// Setup with mock config provider that uses the comprehensive mock server
	mockConfigProvider := &MockConfigProvider{
		config: MarioFaceServiceConfig{
			Host:       mockServer.URL,
			deletePath: "/delete-identity",
		},
	}
	mockTimer := helper.Timer{}

	business := &FaceBussiness{
		config: mockConfigProvider, // Use ConfigProvider interface
		time:   mockTimer,
	}

	ctx := context.Background()
	face := models.Face{
		Name:  stringPtr("user-to-delete"),
		Image: stringPtr("base64-image-data"),
	}
	jwt := "Bearer test-token"

	// Execute
	result := business.Delete(ctx, face, jwt)

	// Verify
	assert.Equal(t, models.Success, result.BasicResponse.Code)
	assert.Equal(t, "Face identity deleted successfully", result.BasicResponse.Message)
}

func TestFaceBusiness_IsUserRegistered_True(t *testing.T) {
	// Use the new mock server instead of creating individual test servers
	mockServer := MockFaceRecognitionServer()
	defer mockServer.Close()

	// Setup with mock config provider that uses the comprehensive mock server
	mockConfigProvider := &MockConfigProvider{
		config: MarioFaceServiceConfig{
			Host:           mockServer.URL,
			registeredPath: "/is-registered",
		},
	}
	mockTimer := helper.Timer{}

	business := &FaceBussiness{
		config: mockConfigProvider, // Use ConfigProvider interface
		time:   mockTimer,
	}

	ctx := context.Background()
	userId := "test-user"
	jwt := "Bearer test-token"

	// Execute
	result, err := business.IsUserRegistered(ctx, userId, jwt)

	// Verify
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestFaceBusiness_IsUserRegistered_False(t *testing.T) {
	// Use the new mock server instead of creating individual test servers
	mockServer := MockFaceRecognitionServer()
	defer mockServer.Close()

	// Setup with mock config provider that uses the comprehensive mock server
	mockConfigProvider := &MockConfigProvider{
		config: MarioFaceServiceConfig{
			Host:           mockServer.URL,
			registeredPath: "/is-registered",
		},
	}
	mockTimer := helper.Timer{}

	business := &FaceBussiness{
		config: mockConfigProvider, // Use ConfigProvider interface
		time:   mockTimer,
	}

	ctx := context.Background()
	userId := "unregistered-user"
	jwt := "Bearer test-token"

	// Execute
	result, err := business.IsUserRegistered(ctx, userId, jwt)

	// Verify
	assert.NoError(t, err)
	assert.False(t, result)
}

func TestMapResponse_Success(t *testing.T) {
	// Test data
	testJSON := `{
		"code": "200",
		"message": "Success",
		"requestId": "test-request-id",
		"userId": "test-user",
		"rawImage": "base64-image-data",
		"flow": "enroll",
		"searchData": {
			"searh_result": {
				"john_doe": "0.95"
			},
			"userId": "john_doe"
		}
	}`

	// Execute
	result, err := MapResponse([]byte(testJSON))

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, "test-user", result.UserId)
	assert.Equal(t, "200", result.BasicResponse.Code)
	assert.Equal(t, "Success", result.BasicResponse.Message)
	assert.Equal(t, "test-request-id", result.BasicResponse.RequestId)
	assert.NotNil(t, result.Data.Image)
	assert.Equal(t, "base64-image-data", *result.Data.Image)
	assert.NotNil(t, result.Data.Name)
	assert.Equal(t, "john_doe", *result.Data.Name)
	assert.NotNil(t, result.Data.Probability)
	assert.Equal(t, 0.95, *result.Data.Probability)
}

func TestMapResponse_InvalidJSON(t *testing.T) {
	// Test with invalid JSON
	invalidJSON := `{"invalid": json}`

	// Execute
	_, err := MapResponse([]byte(invalidJSON))

	// Verify
	assert.Error(t, err)
}

func TestMapResponse_InvalidProbability(t *testing.T) {
	// Test with invalid probability value
	testJSON := `{
		"code": "200",
		"message": "Success",
		"searchData": {
			"searh_result": {
				"john_doe": "invalid_number"
			}
		}
	}`

	// Execute
	result, err := MapResponse([]byte(testJSON))

	// Verify - should not error, just skip the invalid probability
	assert.NoError(t, err)
	assert.Equal(t, "200", result.Code)
	assert.Nil(t, result.Data.Probability) // Should be nil since parsing failed
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
