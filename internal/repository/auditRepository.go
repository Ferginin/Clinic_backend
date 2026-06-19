package repository

import (
	"Clinic_backend/internal/entity"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditRepositoryInterface interface {
	Log(ctx context.Context, log *entity.AuditLog) error
}

type AuditRepository struct {
	db *pgxpool.Pool
}

func NewAuditRepository(db *pgxpool.Pool) AuditRepositoryInterface {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Log(ctx context.Context, log *entity.AuditLog) error {
	query := `
		INSERT INTO audit_logs (admin_id, action, entity_type, entity_id, ip)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, log.AdminID, log.Action, log.EntityType, log.EntityID, log.IP)
	if err != nil {
		return fmt.Errorf("failed to write audit log: %w", err)
	}
	return nil
}
