package profile

import (
	"context"
	"github.com/pkg/errors"
	"golang-ai-management/common"
	"golang-ai-management/proto/pb"
)

type ProfileService interface {
	GetUserProfile(context.Context) (*pb.User, error)
	GetUserById(context.Context, *pb.GetUserByIdReq) (*pb.PublicUserInfoResp, error)
	GetUsersByIds(context.Context, *pb.GetUsersByIdsReq) (*pb.PublicUsersInfoResp, error)
	CreateUser(context.Context, *pb.CreateUserReq) (*pb.NewUserIdResp, error)
}

type business struct {
	profileService ProfileService
	jwtProvider    common.JWTProvider
}

func (b business) GetUserProfile(ctx context.Context) (*pb.User, error) {
	resp, err := b.profileService.GetUserProfile(ctx)

	if err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

func NewBusiness(profileService ProfileService,
	jwtProvider common.JWTProvider) *business {
	return &business{
		profileService: profileService,
		jwtProvider:    jwtProvider,
	}
}
