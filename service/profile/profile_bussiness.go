package profile

import (
	"context"
	"github.com/pkg/errors"
	"golang-ai-management/common"
	helper "golang-ai-management/helpers"
	"golang-ai-management/proto/pb"
	"log/slog"
	"os"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type ProfileService interface {
	GetUserProfile(context.Context) (*pb.User, error)
	GetUserById(context.Context, *pb.GetUserByIdReq) (*pb.PublicUserInfoResp, error)
	GetUsersByIds(context.Context, *pb.GetUsersByIdsReq) (*pb.PublicUsersInfoResp, error)
	CreateUser(context.Context, *pb.CreateUserReq) (*pb.NewUserIdResp, error)
}

type business struct {
	profileService ProfileService
	jwtProvider    common.JWTProvider
	time           helper.Timer
}

func (b business) GetUserProfile(ctx context.Context) (*pb.User, error) {
	b.time.Start()
	logger.Info("request", "method", "Business_GetUserProfile")
	resp, err := b.profileService.GetUserProfile(ctx)
	if err != nil {
		logger.Error("response", "method", "Business_GetUserProfile", "err", err, "ms", b.time.End())
		return nil, errors.WithStack(err)
	}
	logger.Info("response", "method", "Business_GetUserProfile", "data", resp, "ms", b.time.End())
	return resp, nil
}

func NewBusiness(profileService ProfileService,
	jwtProvider common.JWTProvider) *business {
	return &business{
		profileService: profileService,
		jwtProvider:    jwtProvider,
	}
}
