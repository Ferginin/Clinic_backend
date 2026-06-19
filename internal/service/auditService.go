package service

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/repository"
	"context"
	"log/slog"
)

type AuditServiceInterface interface {
	Log(ctx context.Context, adminID int, action, entityType string, entityID *int, ip string)
}

type AuditService struct {
	repo repository.AuditRepositoryInterface
}

func NewAuditService(repo repository.AuditRepositoryInterface) AuditServiceInterface {
	return &AuditService{repo: repo}
}

func (s *AuditService) Log(ctx context.Context, adminID int, action, entityType string, entityID *int, ip string) {
	entry := &entity.AuditLog{
		AdminID:    adminID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		IP:         ip,
	}
	if err := s.repo.Log(ctx, entry); err != nil {
		slog.Error("audit log write failed", "error", err)
	}
}
