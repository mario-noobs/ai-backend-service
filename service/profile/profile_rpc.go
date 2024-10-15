package profile

import (
	"context"
	"github.com/pkg/errors"
	"golang-ai-management/proto/pb"
)

type rpcProfileClient struct {
	client pb.UserServiceClient
}

func (r rpcProfileClient) GetUserById(ctx context.Context, req *pb.GetUserByIdReq) (*pb.PublicUserInfoResp, error) {
	//TODO implement me
	panic("implement me")
}

func (r rpcProfileClient) GetUsersByIds(ctx context.Context, req *pb.GetUsersByIdsReq) (*pb.PublicUsersInfoResp, error) {
	//TODO implement me
	panic("implement me")
}

func (r rpcProfileClient) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.NewUserIdResp, error) {
	//TODO implement me
	panic("implement me")
}

func (r rpcProfileClient) GetUserProfile(ctx context.Context) (*pb.User, error) {
	resp, err := r.client.GetUserProfile(ctx)

	if err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

func NewClient(client pb.UserServiceClient) *rpcProfileClient {
	return &rpcProfileClient{client: client}
}
