package profile

import (
	"context"
	"errors"
	"golang-ai-management/proto/pb"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProfileService is a mock implementation of ProfileService
type MockProfileService struct {
	mock.Mock
}

func (m *MockProfileService) GetUserProfile(ctx context.Context) (*pb.User, error) {
	args := m.Called(ctx)
	return args.Get(0).(*pb.User), args.Error(1)
}

func (m *MockProfileService) GetUserById(ctx context.Context, req *pb.GetUserByIdReq) (*pb.PublicUserInfoResp, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pb.PublicUserInfoResp), args.Error(1)
}

func (m *MockProfileService) GetUsersByIds(ctx context.Context, req *pb.GetUsersByIdsReq) (*pb.PublicUsersInfoResp, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pb.PublicUsersInfoResp), args.Error(1)
}

func (m *MockProfileService) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.NewUserIdResp, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pb.NewUserIdResp), args.Error(1)
}

// MockJWTProvider is a mock implementation of common.JWTProvider
type MockJWTProvider struct {
	mock.Mock
}

func (m *MockJWTProvider) IssueToken(ctx context.Context, id, sub string) (token string, expSecs int, err error) {
	args := m.Called(ctx, id, sub)
	return args.String(0), args.Int(1), args.Error(2)
}

func (m *MockJWTProvider) ParseToken(ctx context.Context, tokenString string) (*jwt.RegisteredClaims, error) {
	args := m.Called(ctx, tokenString)
	return args.Get(0).(*jwt.RegisteredClaims), args.Error(1)
}

func TestBusiness_GetUserProfile(t *testing.T) {
	// Setup
	mockProfileService := new(MockProfileService)
	mockJWTProvider := new(MockJWTProvider)

	business := NewBusiness(mockProfileService, mockJWTProvider)

	ctx := context.Background()

	expectedUser := &pb.User{
		Id:        123,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}

	// Test successful get user profile
	t.Run("successful get user profile", func(t *testing.T) {
		mockProfileService.On("GetUserProfile", ctx).Return(expectedUser, nil).Once()

		result, err := business.GetUserProfile(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		assert.Equal(t, int32(123), result.Id)
		assert.Equal(t, "John", result.FirstName)
		assert.Equal(t, "Doe", result.LastName)
		assert.Equal(t, "john@example.com", result.Email)
		mockProfileService.AssertExpectations(t)
	})

	// Test get user profile with service error
	t.Run("get user profile with service error", func(t *testing.T) {
		expectedError := errors.New("user not found")
		mockProfileService.On("GetUserProfile", ctx).Return((*pb.User)(nil), expectedError).Once()

		result, err := business.GetUserProfile(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "user not found")
		mockProfileService.AssertExpectations(t)
	})

	// Test get user profile with database error
	t.Run("get user profile with database error", func(t *testing.T) {
		dbError := errors.New("database connection failed")
		mockProfileService.On("GetUserProfile", ctx).Return((*pb.User)(nil), dbError).Once()

		result, err := business.GetUserProfile(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database connection failed")
		mockProfileService.AssertExpectations(t)
	})
}

func TestNewBusiness(t *testing.T) {
	mockProfileService := new(MockProfileService)
	mockJWTProvider := new(MockJWTProvider)

	business := NewBusiness(mockProfileService, mockJWTProvider)

	assert.NotNil(t, business)
	assert.Equal(t, mockProfileService, business.profileService)
	assert.Equal(t, mockJWTProvider, business.jwtProvider)
}

func TestBusiness_Timer_Integration(t *testing.T) {
	// This test verifies that the timer is properly used in the business logic
	mockProfileService := new(MockProfileService)
	mockJWTProvider := new(MockJWTProvider)

	business := NewBusiness(mockProfileService, mockJWTProvider)

	ctx := context.Background()
	expectedUser := &pb.User{
		Id:        456,
		FirstName: "Timer",
		LastName:  "Test",
		Email:     "timer@test.com",
	}

	mockProfileService.On("GetUserProfile", ctx).Return(expectedUser, nil).Once()

	// Execute the method that uses timer
	result, err := business.GetUserProfile(ctx)

	// Verify results
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)

	// Note: The timer functionality is tested indirectly here.
	// The actual timer logic is tested in the helpers/timer_test.go
	mockProfileService.AssertExpectations(t)
}
