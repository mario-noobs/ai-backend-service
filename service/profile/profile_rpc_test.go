package profile

import (
	"context"
	"errors"
	"golang-ai-management/proto/pb"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockUserServiceClient is a mock implementation of pb.UserServiceClient
type MockUserServiceClient struct {
	mock.Mock
}

func (m *MockUserServiceClient) GetUserById(ctx context.Context, in *pb.GetUserByIdReq, _ ...grpc.CallOption) (*pb.PublicUserInfoResp, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.PublicUserInfoResp), args.Error(1)
}

func (m *MockUserServiceClient) GetUsersByIds(ctx context.Context, in *pb.GetUsersByIdsReq, _ ...grpc.CallOption) (*pb.PublicUsersInfoResp, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.PublicUsersInfoResp), args.Error(1)
}

func (m *MockUserServiceClient) CreateUser(ctx context.Context, in *pb.CreateUserReq, _ ...grpc.CallOption) (*pb.NewUserIdResp, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.NewUserIdResp), args.Error(1)
}

func (m *MockUserServiceClient) GetUserProfile(ctx context.Context, _ ...grpc.CallOption) (*pb.User, error) {
	args := m.Called(ctx)
	return args.Get(0).(*pb.User), args.Error(1)
}

func TestRpcProfileClient_GetUserProfile(t *testing.T) {
	mockClient := new(MockUserServiceClient)
	rpcClient := NewClient(mockClient)

	ctx := context.Background()

	expectedUser := &pb.User{
		Id:    123,
		Email: "john@example.com",
	}

	// Test successful get user profile
	t.Run("successful get user profile", func(t *testing.T) {
		mockClient.On("GetUserProfile", ctx).Return(expectedUser, nil).Once()

		result, err := rpcClient.GetUserProfile(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		assert.Equal(t, int32(123), result.Id)
		assert.Equal(t, "john@example.com", result.Email)
		mockClient.AssertExpectations(t)
	})

	// Test get user profile with RPC error
	t.Run("get user profile with RPC error", func(t *testing.T) {
		rpcError := errors.New("user not found")
		mockClient.On("GetUserProfile", ctx).Return((*pb.User)(nil), rpcError).Once()

		result, err := rpcClient.GetUserProfile(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "user not found")
		mockClient.AssertExpectations(t)
	})

	// Test get user profile with network error
	t.Run("get user profile with network error", func(t *testing.T) {
		networkError := errors.New("RPC connection failed")
		mockClient.On("GetUserProfile", ctx).Return((*pb.User)(nil), networkError).Once()

		result, err := rpcClient.GetUserProfile(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "RPC connection failed")
		mockClient.AssertExpectations(t)
	})
}

func TestRpcProfileClient_GetUserById(t *testing.T) {
	mockClient := new(MockUserServiceClient)
	rpcClient := NewClient(mockClient)

	ctx := context.Background()
	request := &pb.GetUserByIdReq{
		Id: 123,
	}

	// Test that GetUserById panics (not implemented)
	t.Run("GetUserById panics - not implemented", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = rpcClient.GetUserById(ctx, request)
		})
	})
}

func TestRpcProfileClient_GetUsersByIds(t *testing.T) {
	mockClient := new(MockUserServiceClient)
	rpcClient := NewClient(mockClient)

	ctx := context.Background()
	request := &pb.GetUsersByIdsReq{
		Ids: []int32{123, 456},
	}

	// Test that GetUsersByIds panics (not implemented)
	t.Run("GetUsersByIds panics - not implemented", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = rpcClient.GetUsersByIds(ctx, request)
		})
	})
}

func TestRpcProfileClient_CreateUser(t *testing.T) {
	mockClient := new(MockUserServiceClient)
	rpcClient := NewClient(mockClient)

	ctx := context.Background()
	request := &pb.CreateUserReq{
		Email: "newuser@example.com",
	}

	// Test that CreateUser panics (not implemented)
	t.Run("CreateUser panics - not implemented", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = rpcClient.CreateUser(ctx, request)
		})
	})
}

func TestNewClient(t *testing.T) {
	mockClient := new(MockUserServiceClient)

	rpcClient := NewClient(mockClient)

	assert.NotNil(t, rpcClient)
	assert.Equal(t, mockClient, rpcClient.client)
}
