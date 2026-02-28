package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/zhenglizhi/policy-fit/internal/domain"
)

// AuditRepository 审计仓储
type AuditRepository struct {
	db *sql.DB
}

// NewAuditRepository 创建审计仓储
func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

// CreateLog 写入审计日志
func (r *AuditRepository) CreateLog(ctx context.Context, log *domain.AuditLog) error {
	if log == nil {
		return errors.New("audit log is nil")
	}
	if log.Action == "" {
		return errors.New("action is required")
	}
	if log.TargetType == "" {
		return errors.New("target_type is required")
	}

	detailJSON, err := marshalJSON(log.Detail)
	if err != nil {
		return fmt.Errorf("marshal detail: %w", err)
	}

	const query = `
INSERT INTO audit_log (task_id, actor_id, action, target_type, target_id, detail)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at`

	if err := r.db.QueryRowContext(
		ctx,
		query,
		nullInt64(log.TaskID),
		nullInt64(log.ActorID),
		log.Action,
		log.TargetType,
		nullString(log.TargetID),
		detailJSON,
	).Scan(&log.ID, &log.CreatedAt); err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}

	return nil
}

func nullInt64(value *int64) interface{} {
	if value == nil {
		return nil
	}
	return *value
}
