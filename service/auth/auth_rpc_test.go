package auth

import (
	"context"
	"errors"
	"golang-ai-management/proto"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockUserAuthServiceClient is a mock implementation of pb.UserAuthServiceClient
type MockUserAuthServiceClient struct {
	mock.Mock
}

func (m *MockUserAuthServiceClient) Login(ctx context.Context, in *proto.AuthEmailPassword, opts ...grpc.CallOption) (*proto.TokenResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*proto.TokenResponse), args.Error(1)
}

func (m *MockUserAuthServiceClient) Register(ctx context.Context, in *proto.AuthRegister, opts ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*empty.Empty), args.Error(1)
}

func (m *MockUserAuthServiceClient) Logout(ctx context.Context, in *proto.LogoutRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*empty.Empty), args.Error(1)
}

func (m *MockUserAuthServiceClient) RefreshToken(ctx context.Context, in *proto.RefreshTokenRequest, opts ...grpc.CallOption) (*proto.TokenResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*proto.TokenResponse), args.Error(1)
}

func TestRpcClient_Login(t *testing.T) {
	mockClient := new(MockUserAuthServiceClient)
	rpcClient := rpcClient{client: mockClient}

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
		mockClient.On("Login", ctx, loginData).Return(expectedResponse, nil).Once()

		result, err := rpcClient.Login(ctx, loginData)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		mockClient.AssertExpectations(t)
	})

	// Test login with RPC error
	t.Run("login with RPC error", func(t *testing.T) {
		rpcError := errors.New("RPC connection failed")
		mockClient.On("Login", ctx, loginData).Return((*proto.TokenResponse)(nil), rpcError).Once()

		result, err := rpcClient.Login(ctx, loginData)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "RPC connection failed")
		mockClient.AssertExpectations(t)
	})
}

func TestRpcClient_Register(t *testing.T) {
	mockClient := new(MockUserAuthServiceClient)
	rpcClient := rpcClient{client: mockClient}

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
		mockClient.On("Register", ctx, registerData).Return(expectedResponse, nil).Once()

		result, err := rpcClient.Register(ctx, registerData)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		mockClient.AssertExpectations(t)
	})

	// Test registration with RPC error
	t.Run("registration with RPC error", func(t *testing.T) {
		rpcError := errors.New("user already exists")
		mockClient.On("Register", ctx, registerData).Return((*empty.Empty)(nil), rpcError).Once()

		result, err := rpcClient.Register(ctx, registerData)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "user already exists")
		mockClient.AssertExpectations(t)
	})
}

func TestRpcClient_Logout(t *testing.T) {
	mockClient := new(MockUserAuthServiceClient)
	rpcClient := rpcClient{client: mockClient}

	ctx := context.Background()
	logoutData := &proto.LogoutRequest{
		AccessToken: "access-token",
	}

	expectedResponse := &empty.Empty{}

	// Test successful logout
	t.Run("successful logout", func(t *testing.T) {
		mockClient.On("Logout", ctx, logoutData).Return(expectedResponse, nil).Once()

		result, err := rpcClient.Logout(ctx, logoutData)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		mockClient.AssertExpectations(t)
	})

	// Test logout with RPC error
	t.Run("logout with RPC error", func(t *testing.T) {
		rpcError := errors.New("invalid token")
		mockClient.On("Logout", ctx, logoutData).Return((*empty.Empty)(nil), rpcError).Once()

		result, err := rpcClient.Logout(ctx, logoutData)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid token")
		mockClient.AssertExpectations(t)
	})
}

func TestRpcClient_RefreshToken(t *testing.T) {
	mockClient := new(MockUserAuthServiceClient)
	rpcClient := rpcClient{client: mockClient}

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
		mockClient.On("RefreshToken", ctx, refreshData).Return(expectedResponse, nil).Once()

		result, err := rpcClient.RefreshToken(ctx, refreshData)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		mockClient.AssertExpectations(t)
	})

	// Test token refresh with RPC error
	t.Run("token refresh with RPC error", func(t *testing.T) {
		rpcError := errors.New("refresh token expired")
		mockClient.On("RefreshToken", ctx, refreshData).Return((*proto.TokenResponse)(nil), rpcError).Once()

		result, err := rpcClient.RefreshToken(ctx, refreshData)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "refresh token expired")
		mockClient.AssertExpectations(t)
	})
}
