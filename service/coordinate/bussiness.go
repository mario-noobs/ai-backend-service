package coordinate

import (
	"context"
	"golang-ai-management/common"
	"golang-ai-management/models"
)

type AuthRepository interface {
	AddNewAuth(ctx context.Context, data *models.Auth) error
	GetAuth(ctx context.Context, email string) (*models.Auth, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, firstName, lastName, email string) (newId int, err error)
}

type Hasher interface {
	RandomStr(length int) (string, error)
	HashPassword(salt, password string) (string, error)
	CompareHashPassword(hashedPassword, salt, password string) bool
}

type business struct {
	repository     AuthRepository
	userRepository UserRepository
	jwtProvider    common.JWTProvider
	hasher         Hasher
}

func (b *business) Login(ctx context.Context, data *models.AuthEmailPassword) (*models.TokenResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (b *business) Register(ctx context.Context, data *models.AuthRegister) error {
	//TODO implement me
	panic("implement me")
}

func NewBusiness(repository AuthRepository,
	jwtProvider common.JWTProvider, hasher Hasher) *business {
	return &business{
		repository:  repository,
		jwtProvider: jwtProvider,
		hasher:      hasher,
	}
}
