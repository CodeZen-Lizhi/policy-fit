package service

import (
	"context"
	"strings"
	"time"

	"github.com/zhenglizhi/policy-fit/internal/domain"
	"github.com/zhenglizhi/policy-fit/internal/repository"
)

type analyticsRepository interface {
	CreateEvent(ctx context.Context, event *domain.AnalyticsEvent) error
	CountByEventName(ctx context.Context) (map[string]int64, error)
	CountByEventNameSince(ctx context.Context, since time.Time) (map[string]int64, error)
}

// AnalyticsService 埋点服务
type AnalyticsService struct {
	repo analyticsRepository
	now  func() time.Time
}

// NewAnalyticsService 创建埋点服务
func NewAnalyticsService(repo *repository.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{
		repo: repo,
		now:  time.Now,
	}
}

// Track 记录埋点
func (s *AnalyticsService) Track(
	ctx context.Context,
	eventName string,
	userID *int64,
	taskID *int64,
	properties map[string]interface{},
) error {
	return s.repo.CreateEvent(ctx, &domain.AnalyticsEvent{
		UserID:     userID,
		TaskID:     taskID,
		EventName:  eventName,
		Properties: properties,
	})
}

// Funnel 漏斗数据
func (s *AnalyticsService) Funnel(ctx context.Context) (map[string]int64, error) {
	counts, err := s.countByPeriod(ctx, "all")
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

// FunnelByPeriod 漏斗数据（按周期）
func (s *AnalyticsService) FunnelByPeriod(ctx context.Context, period string) (map[string]int64, error) {
	counts, err := s.countByPeriod(ctx, period)
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
	}, nil
}

// Dashboard 看板概览（按周期）
func (s *AnalyticsService) Dashboard(ctx context.Context, period string) (map[string]interface{}, error) {
	counts, err := s.countByPeriod(ctx, period)
	if err != nil {
		return nil, err
	}

	created := counts["task_created"]
	completed := counts["task_completed"]
	viewed := counts["report_viewed"]
	exported := counts["report_exported"]

	return map[string]interface{}{
		"period":          normalizePeriod(period),
		"task_created":    created,
		"task_completed":  completed,
		"report_viewed":   viewed,
		"report_exported": exported,
		"task_deleted":    counts["task_deleted"],
		"completion_rate": safeRatio(completed, created),
		"view_rate":       safeRatio(viewed, completed),
		"export_rate":     safeRatio(exported, viewed),
		"raw_events":      counts,
	}, nil
}

func (s *AnalyticsService) countByPeriod(ctx context.Context, period string) (map[string]int64, error) {
	switch normalizePeriod(period) {
	case "all":
		return s.repo.CountByEventName(ctx)
	case "week":
		return s.repo.CountByEventNameSince(ctx, s.now().AddDate(0, 0, -7))
	case "month":
		return s.repo.CountByEventNameSince(ctx, s.now().AddDate(0, 0, -30))
	default:
		return nil, ErrInvalidArgument
	}
}

func normalizePeriod(period string) string {
	period = strings.TrimSpace(strings.ToLower(period))
	if period == "" {
		return "week"
	}
	if period == "all" || period == "week" || period == "month" {
		return period
	}
	return period
}

func safeRatio(numerator, denominator int64) float64 {
	if denominator <= 0 || numerator <= 0 {
		return 0
	}
	return float64(numerator) / float64(denominator)
}
