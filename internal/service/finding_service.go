package service

import (
	"context"
	"fmt"

	"github.com/zhenglizhi/policy-fit/internal/domain"
	"github.com/zhenglizhi/policy-fit/internal/repository"
)

// FindingService 风险发现服务
type FindingService struct {
	findingRepo *repository.FindingRepository
}

// NewFindingService 创建风险发现服务
func NewFindingService(findingRepo *repository.FindingRepository) *FindingService {
	return &FindingService{
		findingRepo: findingRepo,
	}
}

// ListByTask 获取任务风险列表
func (s *FindingService) ListByTask(ctx context.Context, taskID int64) ([]domain.RiskFinding, error) {
	if taskID <= 0 {
		return nil, fmt.Errorf("%w: invalid task id", ErrInvalidArgument)
	}
	return s.findingRepo.ListByTask(ctx, taskID)
}

// Summarize 统计风险摘要
func (s *FindingService) Summarize(findings []domain.RiskFinding) map[string]int {
	summary := map[string]int{
		"red":    0,
		"yellow": 0,
		"green":  0,
	}

	for _, f := range findings {
		switch f.Level {
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
