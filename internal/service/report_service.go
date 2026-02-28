package service

import (
	"context"
	"fmt"

	"github.com/zhenglizhi/policy-fit/internal/domain"
	"github.com/zhenglizhi/policy-fit/internal/ruleengine"
)

type findingWriter interface {
	BatchCreateFindings(ctx context.Context, findings []domain.RiskFinding) error
}

type taskStatusWriter interface {
	UpdateTaskStatus(ctx context.Context, id int64, status domain.TaskStatus, riskSummary map[string]int) error
}

// ReportService 报告服务
type ReportService struct {
	matcher     *ruleengine.Matcher
	scorer      *ruleengine.Scorer
	findingRepo findingWriter
	taskRepo    taskStatusWriter
}

// NewReportService 创建报告服务
func NewReportService(
	matcher *ruleengine.Matcher,
	scorer *ruleengine.Scorer,
	findingRepo findingWriter,
	taskRepo taskStatusWriter,
) *ReportService {
	return &ReportService{
		matcher:     matcher,
		scorer:      scorer,
		findingRepo: findingRepo,
		taskRepo:    taskRepo,
	}
}

// GenerateAndSave 生成并保存风险报告
func (s *ReportService) GenerateAndSave(
	ctx context.Context,
	taskID int64,
	healthFacts []domain.HealthFact,
	policyFacts []domain.PolicyFact,
) ([]domain.RiskFinding, map[string]int, error) {
	if taskID <= 0 {
		return nil, nil, fmt.Errorf("%w: invalid task id", ErrInvalidArgument)
	}

	matches := s.matcher.Match(healthFacts, policyFacts)
	findings := s.scorer.Score(matches)
	for i := range findings {
		findings[i].TaskID = taskID
		if len(findings[i].HealthEvidence) == 0 || len(findings[i].PolicyEvidence) == 0 {
			return nil, nil, fmt.Errorf("%w: evidence missing", ErrInvalidArgument)
		}
	}

	if err := s.findingRepo.BatchCreateFindings(ctx, findings); err != nil {
		return nil, nil, err
	}

	summary := SummarizeFindings(findings)
	if err := s.taskRepo.UpdateTaskStatus(ctx, taskID, domain.TaskStatusSuccess, summary); err != nil {
		return nil, nil, err
	}

	return findings, summary, nil
}

// SummarizeFindings 汇总风险结果
func SummarizeFindings(findings []domain.RiskFinding) map[string]int {
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
