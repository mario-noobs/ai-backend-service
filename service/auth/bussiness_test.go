package auth

import (
	"context"
	"errors"
	"golang-ai-management/proto"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, in *proto.AuthEmailPassword) (*proto.TokenResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*proto.TokenResponse), args.Error(1)
}

func (m *MockAuthService) Register(ctx context.Context, in *proto.AuthRegister) (*empty.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*empty.Empty), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, in *proto.LogoutRequest) (*empty.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*empty.Empty), args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, in *proto.RefreshTokenRequest) (*proto.TokenResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*proto.TokenResponse), args.Error(1)
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

// MockHasher is a mock implementation of Hasher
type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) RandomStr(length int) (string, error) {
	args := m.Called(length)
	return args.String(0), args.Error(1)
}

func (m *MockHasher) HashPassword(salt, password string) (string, error) {
	args := m.Called(salt, password)
	return args.String(0), args.Error(1)
}

func (m *MockHasher) CompareHashPassword(hashedPassword, salt, password string) bool {
	args := m.Called(hashedPassword, salt, password)
	return args.Bool(0)
}

func TestBusiness_Login(t *testing.T) {
	// Setup
	mockAuthService := new(MockAuthService)
	mockJWTProvider := new(MockJWTProvider)
	mockHasher := new(MockHasher)

	business := NewBusiness(mockAuthService, mockJWTProvider, mockHasher)

	ctx := context.Background()
	loginData := &proto.AuthEmailPassword{
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedResponse := &proto.TokenResponse{
		AccessToken: &proto.Token{
			Token:     "access-token",
			ExpiredIn: 3600,
		},
		RefreshToken: &proto.Token{
			Token:     "refresh-token",
			ExpiredIn: 86400,
		},
	}

	// Test successful login
	t.Run("successful login", func(t *testing.T) {
		mockAuthService.On("Login", ctx, loginData).Return(expectedResponse, nil).Once()

		result, err := business.Login(ctx, loginData)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		mockAuthService.AssertExpectations(t)
	})

	// Test login with error
	t.Run("login with error", func(t *testing.T) {
		expectedError := errors.New("authentication failed")
		mockAuthService.On("Login", ctx, loginData).Return((*proto.TokenResponse)(nil), expectedError).Once()

		result, err := business.Login(ctx, loginData)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "authentication failed")
		mockAuthService.AssertExpectations(t)
	})
}

func TestBusiness_Register(t *testing.T) {
	// Setup
	mockAuthService := new(MockAuthService)
	mockJWTProvider := new(MockJWTProvider)
	mockHasher := new(MockHasher)

	business := NewBusiness(mockAuthService, mockJWTProvider, mockHasher)

	ctx := context.Background()
	registerData := &proto.AuthRegister{
		FirstName: "Test",
		LastName:  "User",
		AuthEmailPassword: &proto.AuthEmailPassword{
			Email:    "test@example.com",
			Password: "password123",
		},
	}

	expectedResponse := &empty.Empty{}

	// Test successful registration
	t.Run("successful registration", func(t *testing.T) {
		mockAuthService.On("Register", ctx, registerData).Return(expectedResponse, nil).Once()

		result, err := business.Register(ctx, registerData)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		mockAuthService.AssertExpectations(t)
	})

	// Test registration with error
	t.Run("registration with error", func(t *testing.T) {
		expectedError := errors.New("email already exists")
		mockAuthService.On("Register", ctx, registerData).Return((*empty.Empty)(nil), expectedError).Once()

		result, err := business.Register(ctx, registerData)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "email already exists")
		mockAuthService.AssertExpectations(t)
	})
}

func TestBusiness_Logout(t *testing.T) {
	// Setup
	mockAuthService := new(MockAuthService)
	mockJWTProvider := new(MockJWTProvider)
	mockHasher := new(MockHasher)

	business := NewBusiness(mockAuthService, mockJWTProvider, mockHasher)

	ctx := context.Background()
	logoutData := &proto.LogoutRequest{
		AccessToken: "access-token",
	}

	expectedResponse := &empty.Empty{}

	// Test successful logout
	t.Run("successful logout", func(t *testing.T) {
		mockAuthService.On("Logout", ctx, logoutData).Return(expectedResponse, nil).Once()

		result, err := business.Logout(ctx, logoutData)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		mockAuthService.AssertExpectations(t)
	})

	// Test logout with error
	t.Run("logout with error", func(t *testing.T) {
		expectedError := errors.New("invalid token")
		mockAuthService.On("Logout", ctx, logoutData).Return((*empty.Empty)(nil), expectedError).Once()

		result, err := business.Logout(ctx, logoutData)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid token")
		mockAuthService.AssertExpectations(t)
	})
}

func TestBusiness_RefreshToken(t *testing.T) {
	// Setup
	mockAuthService := new(MockAuthService)
	mockJWTProvider := new(MockJWTProvider)
	mockHasher := new(MockHasher)

	business := NewBusiness(mockAuthService, mockJWTProvider, mockHasher)

	ctx := context.Background()
	refreshData := &proto.RefreshTokenRequest{
		RefreshToken: "refresh-token",
	}

	expectedResponse := &proto.TokenResponse{
		AccessToken: &proto.Token{
			Token:     "new-access-token",
			ExpiredIn: 3600,
		},
		RefreshToken: &proto.Token{
			Token:     "new-refresh-token",
			ExpiredIn: 86400,
		},
	}

	// Test successful token refresh
	t.Run("successful token refresh", func(t *testing.T) {
		mockAuthService.On("RefreshToken", ctx, refreshData).Return(expectedResponse, nil).Once()

		result, err := business.RefreshToken(ctx, refreshData)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		mockAuthService.AssertExpectations(t)
	})

	// Test token refresh with error
	t.Run("token refresh with error", func(t *testing.T) {
		expectedError := errors.New("refresh token expired")
		mockAuthService.On("RefreshToken", ctx, refreshData).Return((*proto.TokenResponse)(nil), expectedError).Once()

		result, err := business.RefreshToken(ctx, refreshData)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "refresh token expired")
		mockAuthService.AssertExpectations(t)
	})
}

func TestNewBusiness(t *testing.T) {
	mockAuthService := new(MockAuthService)
	mockJWTProvider := new(MockJWTProvider)
	mockHasher := new(MockHasher)

	business := NewBusiness(mockAuthService, mockJWTProvider, mockHasher)

	assert.NotNil(t, business)
	assert.Equal(t, mockAuthService, business.authService)
	assert.Equal(t, mockJWTProvider, business.jwtProvider)
	assert.Equal(t, mockHasher, business.hasher)
}
