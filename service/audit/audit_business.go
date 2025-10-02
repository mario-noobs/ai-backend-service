package audit

import (
	"context"
	"golang-ai-management/models"
)

type Repository interface {
	GetAuditLogsByUserID(ctx context.Context, userID string, limit, offset int) ([]*models.AuditLog, error)
	GetAuditLogs(ctx context.Context, limit, offset int) ([]*models.AuditLog, error)
}

type AuditBusiness struct {
	repo Repository
}

func NewAuditBusiness(repo Repository) *AuditBusiness {
	return &AuditBusiness{repo: repo}
}

func (biz *AuditBusiness) GetUserAuditLogs(ctx context.Context, userID string, limit, offset int) ([]*models.AuditLog, error) {
	return biz.repo.GetAuditLogsByUserID(ctx, userID, limit, offset)
}

func (biz *AuditBusiness) GetAllAuditLogs(ctx context.Context, limit, offset int) ([]*models.AuditLog, error) {
	return biz.repo.GetAuditLogs(ctx, limit, offset)
}
