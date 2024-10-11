package coordinate

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"golang-ai-management/common"
	"golang-ai-management/proto/pb"
)

type AuthService interface {
	Login(ctx context.Context, in *pb.AuthEmailPassword) (*pb.TokenResponse, error)
	Register(ctx context.Context, in *pb.AuthRegister) (*empty.Empty, error)
}

type Hasher interface {
	RandomStr(length int) (string, error)
	HashPassword(salt, password string) (string, error)
	CompareHashPassword(hashedPassword, salt, password string) bool
}

type business struct {
	authService AuthService
	jwtProvider common.JWTProvider
	hasher      Hasher
}

func (b business) Login(ctx context.Context, data *pb.AuthEmailPassword) (*pb.TokenResponse, error) {
	resp, err := b.authService.Login(ctx, data)

	if err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

func (b business) Register(ctx context.Context, in *pb.AuthRegister) (*empty.Empty, error) {
	resp, err := b.authService.Register(ctx, in)

	if err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

func NewBusiness(authService AuthService,
	jwtProvider common.JWTProvider, hasher Hasher) *business {
	return &business{
		authService: authService,
		jwtProvider: jwtProvider,
		hasher:      hasher,
	}
}
