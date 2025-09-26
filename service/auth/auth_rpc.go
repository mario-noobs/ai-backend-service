package auth

import (
	"context"
	"golang-ai-management/proto"
	"golang-ai-management/proto/pb"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
)

type rpcClient struct {
	client pb.UserAuthServiceClient
}

func (r rpcClient) Login(ctx context.Context, in *proto.AuthEmailPassword) (*proto.TokenResponse, error) {
	resp, err := r.client.Login(ctx, in)

	if err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

func (r rpcClient) Register(ctx context.Context, in *proto.AuthRegister) (*empty.Empty, error) {
	resp, err := r.client.Register(ctx, in)

	if err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

func (r rpcClient) Logout(ctx context.Context, in *proto.LogoutRequest) (*empty.Empty, error) {
	resp, err := r.client.Logout(ctx, in)

	if err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

func (r rpcClient) RefreshToken(ctx context.Context, in *proto.RefreshTokenRequest) (*proto.TokenResponse, error) {
	resp, err := r.client.RefreshToken(ctx, in)

	if err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

func NewClient(client pb.UserAuthServiceClient) *rpcClient {
	return &rpcClient{client: client}
}
