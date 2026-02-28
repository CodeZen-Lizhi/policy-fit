package service

import (
	"context"
	"testing"

	"github.com/zhenglizhi/policy-fit/internal/domain"
	"github.com/zhenglizhi/policy-fit/internal/ruleengine"
)

func TestReportServiceGenerateAndSave(t *testing.T) {
	cfg := &ruleengine.TopicsConfig{
		Topics: map[string]ruleengine.TopicRule{
			"hypertension": {
				HitKeywords: []string{"高血压"},
				PolicyTypes: []string{"exclusion"},
			},
		},
	}
	service := NewReportService(
		ruleengine.NewMatcher(cfg),
		ruleengine.NewScorer(),
		&fakeFindingRepo{},
		&fakeTaskRepo{},
	)

	diagnosed := true
	findings, summary, err := service.GenerateAndSave(
		context.Background(),
		1,
		[]domain.HealthFact{
			{
				Category:   "hypertension",
				Label:      "高血压",
				Evidence:   domain.EvidenceDetail{Loc: "para_1", Text: "血压偏高", Date: "2026-01-01"},
				Diagnosed:  &diagnosed,
				Confidence: 0.9,
			},
		},
		[]domain.PolicyFact{
			{
				Type:       "exclusion",
				Title:      "责任免除",
				Content:    "既往症免责",
				Loc:        "para_100",
				Confidence: 0.9,
			},
		},
	)
	if err != nil {
		t.Fatalf("GenerateAndSave error: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("unexpected findings len: %d", len(findings))
	}
	if summary["red"] != 1 {
		t.Fatalf("unexpected summary: %#v", summary)
	}
}

type fakeFindingRepo struct{}

func (f *fakeFindingRepo) BatchCreateFindings(context.Context, []domain.RiskFinding) error {
	return nil
}

type fakeTaskRepo struct{}

func (f *fakeTaskRepo) UpdateTaskStatus(context.Context, int64, domain.TaskStatus, map[string]int) error {
	return nil
}
