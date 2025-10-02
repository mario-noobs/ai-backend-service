package mysql

import (
	"context"
	"golang-ai-management/models"

	"github.com/pkg/errors"
	"github.com/viettranx/service-context/core"
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

// Audit Log Methods
func (repo *mysqlRepo) AddAuditLog(ctx context.Context, data *models.AuditLog) error {
	if err := repo.db.Table(data.TableName()).Create(data).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (repo *mysqlRepo) GetAuditLogs(ctx context.Context, limit, offset int) ([]*models.AuditLog, error) {
	var auditLogs []*models.AuditLog

	if err := repo.db.
		Table("audit_logs").
		Order("time DESC").
		Limit(limit).
		Offset(offset).
		Find(&auditLogs).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return auditLogs, nil
}

func (repo *mysqlRepo) GetAuditLogsByIPAddress(ctx context.Context, ipAddress string, limit, offset int) ([]*models.AuditLog, error) {
	var auditLogs []*models.AuditLog

	if err := repo.db.
		Table("audit_logs").
		Where("ip_address = ?", ipAddress).
		Order("time DESC").
		Limit(limit).
		Offset(offset).
		Find(&auditLogs).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return auditLogs, nil
}

func (repo *mysqlRepo) GetAuditLogsByAPICall(ctx context.Context, apiCall string, limit, offset int) ([]*models.AuditLog, error) {
	var auditLogs []*models.AuditLog

	if err := repo.db.
		Table("audit_logs").
		Where("api_call = ?", apiCall).
		Order("time DESC").
		Limit(limit).
		Offset(offset).
		Find(&auditLogs).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return auditLogs, nil
}

func (repo *mysqlRepo) GetAuditLogsByUserID(ctx context.Context, userID string, limit, offset int) ([]*models.AuditLog, error) {
	var auditLogs []*models.AuditLog

	if err := repo.db.
		Table("audit_logs").
		Where("user_id = ?", userID).
		Order("time DESC").
		Limit(limit).
		Offset(offset).
		Find(&auditLogs).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return auditLogs, nil
}
