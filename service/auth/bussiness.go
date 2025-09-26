package auth

import (
	"context"
	"golang-ai-management/common"
	"golang-ai-management/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
)

type AuthService interface {
	Login(ctx context.Context, in *proto.AuthEmailPassword) (*proto.TokenResponse, error)
	Register(ctx context.Context, in *proto.AuthRegister) (*empty.Empty, error)
	Logout(ctx context.Context, in *proto.LogoutRequest) (*empty.Empty, error)
	RefreshToken(ctx context.Context, in *proto.RefreshTokenRequest) (*proto.TokenResponse, error)
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

func (b business) Login(ctx context.Context, data *proto.AuthEmailPassword) (*proto.TokenResponse, error) {
	resp, err := b.authService.Login(ctx, data)

	if err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

func (b business) Register(ctx context.Context, in *proto.AuthRegister) (*empty.Empty, error) {
	resp, err := b.authService.Register(ctx, in)

	if err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

func (b business) Logout(ctx context.Context, in *proto.LogoutRequest) (*empty.Empty, error) {
	resp, err := b.authService.Logout(ctx, in)

	if err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

func (b *business) RefreshToken(ctx context.Context, in *proto.RefreshTokenRequest) (*proto.TokenResponse, error) {
	resp, err := b.authService.RefreshToken(ctx, in)

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
