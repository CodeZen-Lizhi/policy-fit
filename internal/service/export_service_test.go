package service

import (
	"context"
	"testing"
	"time"

	"github.com/zhenglizhi/policy-fit/internal/domain"
)

func TestExportServiceMarkdownAndPDF(t *testing.T) {
	taskSvc := &fakeExportTaskService{
		task: &domain.AnalysisTask{
			ID:        1,
			UserID:    1,
			Status:    domain.TaskStatusSuccess,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	findingSvc := &fakeExportFindingService{
		findings: []domain.RiskFinding{
			{
				Topic:   "hypertension",
				Level:   domain.RiskLevelRed,
				Summary: "risk",
				HealthEvidence: []domain.Evidence{
					{Loc: "para_1", Text: "bp high"},
				},
				PolicyEvidence: []domain.Evidence{
					{Loc: "para_2", Text: "exclusion"},
				},
				Questions: []string{"q1", "q2"},
			},
		},
	}

	svc := NewExportService(taskSvc, findingSvc)
	md, err := svc.ExportMarkdown(context.Background(), 1, 1, "zh-CN")
	if err != nil {
		t.Fatalf("ExportMarkdown error: %v", err)
	}
	if len(md) == 0 {
		t.Fatalf("expected markdown output")
	}

	enMD, err := svc.ExportMarkdown(context.Background(), 1, 1, "en-US")
	if err != nil {
		t.Fatalf("ExportMarkdown en error: %v", err)
	}
	if string(enMD) == string(md) {
		t.Fatalf("expected different language output")
	}

	pdf, err := svc.ExportPDF(context.Background(), 1, 1, "en-US")
	if err != nil {
		t.Fatalf("ExportPDF error: %v", err)
	}
	if len(pdf) == 0 {
		t.Fatalf("expected pdf output")
	}
}

type fakeExportTaskService struct {
	task *domain.AnalysisTask
}

func (f *fakeExportTaskService) CreateTask(context.Context, int64, string) (*domain.AnalysisTask, error) {
	return f.task, nil
}

func (f *fakeExportTaskService) GetTask(context.Context, int64, int64) (*domain.AnalysisTask, error) {
	return f.task, nil
}

func (f *fakeExportTaskService) RunTask(context.Context, int64, int64) error {
	return nil
}

func (f *fakeExportTaskService) DeleteTask(context.Context, int64, int64) error {
	return nil
}

type fakeExportFindingService struct {
	findings []domain.RiskFinding
}

func (f *fakeExportFindingService) ListByTask(context.Context, int64) ([]domain.RiskFinding, error) {
	return f.findings, nil
}

func (f *fakeExportFindingService) Summarize(findings []domain.RiskFinding) map[string]int {
	summary := map[string]int{"red": 0, "yellow": 0, "green": 0}
	for _, finding := range findings {
		switch finding.Level {
		case domain.RiskLevelRed:
			summary["red"]++
		case domain.RiskLevelYellow:
			summary["yellow"]++
		case domain.RiskLevelGreen:
			summary["green"]++
		}
	}
	return summary
}
