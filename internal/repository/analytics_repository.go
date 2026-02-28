package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zhenglizhi/policy-fit/internal/domain"
)

// AnalyticsRepository 埋点仓储
type AnalyticsRepository struct {
	db *sql.DB
}

// NewAnalyticsRepository 创建埋点仓储
func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

// CreateEvent 写入埋点
func (r *AnalyticsRepository) CreateEvent(ctx context.Context, event *domain.AnalyticsEvent) error {
	properties, err := marshalJSON(event.Properties)
	if err != nil {
		return err
	}

	const query = `
INSERT INTO analytics_event(user_id, task_id, event_name, properties)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at`
	return r.db.QueryRowContext(
		ctx,
		query,
		nullInt64(event.UserID),
		nullInt64(event.TaskID),
		event.EventName,
		properties,
	).Scan(&event.ID, &event.CreatedAt)
}

// CountByEventName 统计埋点数量
func (r *AnalyticsRepository) CountByEventName(ctx context.Context) (map[string]int64, error) {
	return r.countByEventNameWithCondition(ctx, "", nil)
}

// CountByEventNameSince 统计某个时间点后的埋点数量
func (r *AnalyticsRepository) CountByEventNameSince(ctx context.Context, since time.Time) (map[string]int64, error) {
	return r.countByEventNameWithCondition(ctx, "WHERE created_at >= $1", []interface{}{since})
}

func (r *AnalyticsRepository) countByEventNameWithCondition(
	ctx context.Context,
	condition string,
	args []interface{},
) (map[string]int64, error) {
	const query = `
SELECT event_name, COUNT(1)
FROM analytics_event
%s
GROUP BY event_name`
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(query, condition), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]int64{}
	for rows.Next() {
		var name string
		var count int64
		if err := rows.Scan(&name, &count); err != nil {
			return nil, err
		}
		result[name] = count
	}
	return result, rows.Err()
}

// FunnelOverview 输出漏斗统计
func (r *AnalyticsRepository) FunnelOverview(ctx context.Context) (map[string]int64, error) {
	counts, err := r.CountByEventName(ctx)
	if err != nil {
		return nil, err
	}
	return map[string]int64{
		"task_created":      counts["task_created"],
		"document_uploaded": counts["document_uploaded"],
		"task_run":          counts["task_run"],
		"task_completed":    counts["task_completed"],
		"report_viewed":     counts["report_viewed"],
		"report_exported":   counts["report_exported"],
		"task_deleted":      counts["task_deleted"],
	}, nil
}

// ListRecentEvents 查询最近埋点
func (r *AnalyticsRepository) ListRecentEvents(ctx context.Context, limit int) ([]domain.AnalyticsEvent, error) {
	if limit <= 0 {
		limit = 100
	}

	const query = `
SELECT id, user_id, task_id, event_name, properties, created_at
FROM analytics_event
ORDER BY id DESC
LIMIT $1`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]domain.AnalyticsEvent, 0, limit)
	for rows.Next() {
		var (
			event          domain.AnalyticsEvent
			userID, taskID sql.NullInt64
			raw            []byte
		)
		if err := rows.Scan(&event.ID, &userID, &taskID, &event.EventName, &raw, &event.CreatedAt); err != nil {
			return nil, err
		}
		if userID.Valid {
			event.UserID = &userID.Int64
		}
		if taskID.Valid {
			event.TaskID = &taskID.Int64
		}
		if len(raw) > 0 {
			if err := json.Unmarshal(raw, &event.Properties); err != nil {
				return nil, fmt.Errorf("unmarshal analytics properties: %w", err)
			}
		}
		events = append(events, event)
	}
	return events, rows.Err()
}
