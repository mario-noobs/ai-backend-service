package auth

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"golang-ai-management/proto/pb"
)

type rpcClient struct {
	client pb.UserAuthServiceClient
}

func (r rpcClient) Login(ctx context.Context, in *pb.AuthEmailPassword) (*pb.TokenResponse, error) {
	resp, err := r.client.Login(ctx, in)

	if err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

func (r rpcClient) Register(ctx context.Context, in *pb.AuthRegister) (*empty.Empty, error) {
	resp, err := r.client.Register(ctx, in)

	if err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

func NewClient(client pb.UserAuthServiceClient) *rpcClient {
	return &rpcClient{client: client}
}
