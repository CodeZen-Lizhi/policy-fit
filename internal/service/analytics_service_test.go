package service

import (
	"context"
	"testing"
	"time"

	"github.com/zhenglizhi/policy-fit/internal/domain"
)

func TestAnalyticsServiceDashboardByPeriod(t *testing.T) {
	repo := &fakeAnalyticsRepo{
		all: map[string]int64{
			"task_created":    12,
			"task_completed":  9,
			"report_viewed":   7,
			"report_exported": 3,
		},
		since: map[string]int64{
			"task_created":    6,
			"task_completed":  5,
			"report_viewed":   4,
			"report_exported": 2,
		},
	}
	svc := &AnalyticsService{
		repo: repo,
		now: func() time.Time {
			return time.Date(2026, 2, 28, 10, 0, 0, 0, time.UTC)
		},
	}

	week, err := svc.Dashboard(context.Background(), "week")
	if err != nil {
		t.Fatalf("dashboard week: %v", err)
	}
	if week["period"] != "week" {
		t.Fatalf("period should be week, got %v", week["period"])
	}
	if week["completion_rate"] == nil {
		t.Fatalf("completion_rate should exist")
	}
	if !repo.calledSince {
		t.Fatalf("expected CountByEventNameSince to be called")
	}

	_, err = svc.Dashboard(context.Background(), "invalid")
	if err != ErrInvalidArgument {
		t.Fatalf("expected ErrInvalidArgument, got %v", err)
	}
}

type fakeAnalyticsRepo struct {
	all         map[string]int64
	since       map[string]int64
	calledSince bool
}

func (f *fakeAnalyticsRepo) CreateEvent(_ context.Context, _ *domain.AnalyticsEvent) error {
	return nil
}

func (f *fakeAnalyticsRepo) CountByEventName(_ context.Context) (map[string]int64, error) {
	return f.all, nil
}

func (f *fakeAnalyticsRepo) CountByEventNameSince(_ context.Context, _ time.Time) (map[string]int64, error) {
	f.calledSince = true
	return f.since, nil
}
