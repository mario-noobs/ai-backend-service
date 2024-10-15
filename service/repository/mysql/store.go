package mysql

import (
	"context"
	"github.com/pkg/errors"
	"github.com/viettranx/service-context/core"
	"golang-ai-management/models"
	"gorm.io/gorm"
)

type mysqlRepo struct {
	db *gorm.DB
}

func NewMySQLRepository(db *gorm.DB) *mysqlRepo {
	return &mysqlRepo{db: db}
}

func (repo *mysqlRepo) AddNewAuth(ctx context.Context, data *models.Auth) error {
	if err := repo.db.Table(data.TableName()).Create(data).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (repo *mysqlRepo) GetAuth(ctx context.Context, email string) (*models.Auth, error) {
	var data models.Auth

	if err := repo.db.
		Table(data.TableName()).
		Where("email = ?", email).
		First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrRecordNotFound
		}

		return nil, errors.WithStack(err)
	}

	return &data, nil
}
