package service

import (
	"context"
	"database/sql"
	"fmt"
)

// DataRetentionService 数据保留清理服务
type DataRetentionService struct {
	db *sql.DB
}

// NewDataRetentionService 创建清理服务
func NewDataRetentionService(db *sql.DB) *DataRetentionService {
	return &DataRetentionService{db: db}
}

// CleanupExpired 删除超过 retentionDays 的任务数据
func (s *DataRetentionService) CleanupExpired(ctx context.Context, retentionDays int) (int64, error) {
	if retentionDays <= 0 {
		return 0, fmt.Errorf("%w: retention days must be positive", ErrInvalidArgument)
	}

	const query = `
DELETE FROM analysis_task
WHERE created_at < NOW() - ($1 || ' days')::interval`
	result, err := s.db.ExecContext(ctx, query, retentionDays)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}
