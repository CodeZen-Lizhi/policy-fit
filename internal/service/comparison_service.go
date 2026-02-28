package service

import (
	"context"
	"fmt"

	"github.com/zhenglizhi/policy-fit/internal/domain"
)

// ComparisonItem 对比项
type ComparisonItem struct {
	Topic     string           `json:"topic"`
	FromLevel domain.RiskLevel `json:"from_level,omitempty"`
	ToLevel   domain.RiskLevel `json:"to_level,omitempty"`
	Change    string           `json:"change"`
}

// ComparisonResult 对比结果
type ComparisonResult struct {
	FromTaskID    int64            `json:"from_task_id"`
	ToTaskID      int64            `json:"to_task_id"`
	NewRisks      []ComparisonItem `json:"new_risks"`
	ResolvedRisks []ComparisonItem `json:"resolved_risks"`
	LevelChanges  []ComparisonItem `json:"level_changes"`
	Summary       string           `json:"summary"`
}

// ComparisonService 历史报告对比服务
type ComparisonService struct {
	taskSvc    exportTaskReader
	findingSvc exportFindingReader
}

// NewComparisonService 创建对比服务
func NewComparisonService(taskSvc exportTaskReader, findingSvc exportFindingReader) *ComparisonService {
	return &ComparisonService{
		taskSvc:    taskSvc,
		findingSvc: findingSvc,
	}
}

// CompareTasks 对比两次任务结果
func (s *ComparisonService) CompareTasks(
	ctx context.Context,
	fromTaskID int64,
	toTaskID int64,
	actorID int64,
) (*ComparisonResult, error) {
	if fromTaskID <= 0 || toTaskID <= 0 {
		return nil, fmt.Errorf("%w: invalid task ids", ErrInvalidArgument)
	}
	if _, err := s.taskSvc.GetTask(ctx, fromTaskID, actorID); err != nil {
		return nil, err
	}
	if _, err := s.taskSvc.GetTask(ctx, toTaskID, actorID); err != nil {
		return nil, err
	}

	fromFindings, err := s.findingSvc.ListByTask(ctx, fromTaskID)
	if err != nil {
		return nil, err
	}
	toFindings, err := s.findingSvc.ListByTask(ctx, toTaskID)
	if err != nil {
		return nil, err
	}

	fromMap := mapTopicLevel(fromFindings)
	toMap := mapTopicLevel(toFindings)

	result := &ComparisonResult{
		FromTaskID: fromTaskID,
		ToTaskID:   toTaskID,
	}

	for topic, toLevel := range toMap {
		fromLevel, exists := fromMap[topic]
		if !exists {
			result.NewRisks = append(result.NewRisks, ComparisonItem{
				Topic:   topic,
				ToLevel: toLevel,
				Change:  "new",
			})
			continue
		}
		if fromLevel != toLevel {
			result.LevelChanges = append(result.LevelChanges, ComparisonItem{
				Topic:     topic,
				FromLevel: fromLevel,
				ToLevel:   toLevel,
				Change:    "level_changed",
			})
		}
	}

	for topic, fromLevel := range fromMap {
		if _, exists := toMap[topic]; !exists {
			result.ResolvedRisks = append(result.ResolvedRisks, ComparisonItem{
				Topic:     topic,
				FromLevel: fromLevel,
				Change:    "resolved",
			})
		}
	}

	result.Summary = fmt.Sprintf(
		"新增风险 %d 项，等级变化 %d 项，消失风险 %d 项",
		len(result.NewRisks),
		len(result.LevelChanges),
		len(result.ResolvedRisks),
	)
	return result, nil
}

func mapTopicLevel(findings []domain.RiskFinding) map[string]domain.RiskLevel {
	m := make(map[string]domain.RiskLevel, len(findings))
	for _, finding := range findings {
		m[finding.Topic] = finding.Level
	}
	return m
}
