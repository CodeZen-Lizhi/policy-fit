package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zhenglizhi/policy-fit/internal/domain"
)

// TaskRepository 任务仓储
type TaskRepository struct {
	db *sql.DB
}

// NewTaskRepository 创建任务仓储
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// CreateTask 创建任务
func (r *TaskRepository) CreateTask(ctx context.Context, task *domain.AnalysisTask) error {
	if task == nil {
		return errors.New("task is nil")
	}
	if task.Status == "" {
		task.Status = domain.TaskStatusPending
	}

	riskSummaryJSON, err := marshalJSON(task.RiskSummary)
	if err != nil {
		return fmt.Errorf("marshal risk summary: %w", err)
	}

	const query = `
INSERT INTO analysis_task (user_id, request_id, status, risk_summary)
VALUES ($1, NULLIF($2, ''), $3, $4)
RETURNING id, created_at, updated_at`

	if err := r.db.QueryRowContext(
		ctx,
		query,
		task.UserID,
		task.RequestID,
		task.Status,
		riskSummaryJSON,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt); err != nil {
		return fmt.Errorf("insert task: %w", err)
	}

	return nil
}

// GetTask 获取任务
func (r *TaskRepository) GetTask(ctx context.Context, id int64) (*domain.AnalysisTask, error) {
	const query = `
SELECT id, user_id, request_id, status, risk_summary, created_at, updated_at
FROM analysis_task
WHERE id = $1`

	var (
		task           domain.AnalysisTask
		requestID      sql.NullString
		riskSummaryRaw []byte
	)

	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.UserID,
		&requestID,
		&task.Status,
		&riskSummaryRaw,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("query task: %w", err)
	}

	if requestID.Valid {
		task.RequestID = requestID.String
	}

	if len(riskSummaryRaw) > 0 {
		if err := json.Unmarshal(riskSummaryRaw, &task.RiskSummary); err != nil {
			return nil, fmt.Errorf("unmarshal risk summary: %w", err)
		}
	}

	return &task, nil
}

// UpdateTaskStatus 更新任务状态
func (r *TaskRepository) UpdateTaskStatus(
	ctx context.Context,
	id int64,
	status domain.TaskStatus,
	riskSummary map[string]int,
) error {
	if status == "" {
		return errors.New("status is required")
	}

	riskSummaryJSON, err := marshalJSON(riskSummary)
	if err != nil {
		return fmt.Errorf("marshal risk summary: %w", err)
	}

	const query = `
UPDATE analysis_task
SET status = $2,
    risk_summary = COALESCE($3::jsonb, risk_summary)
WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id, status, riskSummaryJSON)
	if err != nil {
		return fmt.Errorf("update task status: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read affected rows: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteTask 删除任务
func (r *TaskRepository) DeleteTask(ctx context.Context, id int64) error {
	const query = `DELETE FROM analysis_task WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read affected rows: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func marshalJSON(value interface{}) ([]byte, error) {
	if value == nil {
		return nil, nil
	}
	return json.Marshal(value)
}
