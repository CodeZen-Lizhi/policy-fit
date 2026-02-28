package service

import (
	"context"
	"testing"

	"github.com/zhenglizhi/policy-fit/internal/domain"
)

func TestComparisonServiceCompareTasks(t *testing.T) {
	taskSvc := &comparisonTaskService{}
	findingSvc := &comparisonFindingService{
		data: map[int64][]domain.RiskFinding{
			1: {
				{Topic: "hypertension", Level: domain.RiskLevelRed},
				{Topic: "diabetes", Level: domain.RiskLevelYellow},
			},
			2: {
				{Topic: "hypertension", Level: domain.RiskLevelYellow},
				{Topic: "fatty_liver", Level: domain.RiskLevelYellow},
			},
		},
	}
	svc := NewComparisonService(taskSvc, findingSvc)

	result, err := svc.CompareTasks(context.Background(), 1, 2, 1)
	if err != nil {
		t.Fatalf("CompareTasks error: %v", err)
	}

	if len(result.NewRisks) != 1 {
		t.Fatalf("expected 1 new risk, got %d", len(result.NewRisks))
	}
	if len(result.ResolvedRisks) != 1 {
		t.Fatalf("expected 1 resolved risk, got %d", len(result.ResolvedRisks))
	}
	if len(result.LevelChanges) != 1 {
		t.Fatalf("expected 1 level change, got %d", len(result.LevelChanges))
	}
}

type comparisonTaskService struct{}

func (c *comparisonTaskService) GetTask(context.Context, int64, int64) (*domain.AnalysisTask, error) {
	return &domain.AnalysisTask{ID: 1, UserID: 1}, nil
}

type comparisonFindingService struct {
	data map[int64][]domain.RiskFinding
}

func (c *comparisonFindingService) ListByTask(_ context.Context, taskID int64) ([]domain.RiskFinding, error) {
	return c.data[taskID], nil
}
