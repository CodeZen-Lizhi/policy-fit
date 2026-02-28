package service

import (
	"context"
	"fmt"
	"time"

	"github.com/zhenglizhi/policy-fit/internal/domain"
	"github.com/zhenglizhi/policy-fit/internal/repository"
)

// RuleConfigService 规则配置服务
type RuleConfigService struct {
	repo      *repository.RuleConfigRepository
	auditRepo *repository.AuditRepository
}

// NewRuleConfigService 创建规则配置服务
func NewRuleConfigService(repo *repository.RuleConfigRepository, auditRepo *repository.AuditRepository) *RuleConfigService {
	return &RuleConfigService{
		repo:      repo,
		auditRepo: auditRepo,
	}
}

// Publish 发布新规则版本
func (s *RuleConfigService) Publish(
	ctx context.Context,
	content map[string]interface{},
	changelog string,
	actorID int64,
) (*domain.RuleConfigVersion, error) {
	version := fmt.Sprintf("v%s", time.Now().Format("20060102150405"))
	cfg := &domain.RuleConfigVersion{
		Version:   version,
		Changelog: changelog,
		Content:   content,
		IsActive:  true,
		IsGray:    false,
		CreatedBy: optionalActorID(actorID),
	}
	if err := s.repo.CreateVersion(ctx, cfg); err != nil {
		return nil, err
	}
	if err := s.repo.SetActiveVersion(ctx, cfg.Version); err != nil {
		return nil, err
	}
	_ = s.auditRepo.CreateLog(ctx, &domain.AuditLog{
		ActorID:    optionalActorID(actorID),
		Action:     "publish_rule_version",
		TargetType: "rule_config",
		TargetID:   cfg.Version,
		Detail: map[string]interface{}{
			"changelog": changelog,
		},
	})
	return cfg, nil
}

// Rollback 回滚至指定版本
func (s *RuleConfigService) Rollback(ctx context.Context, version string, actorID int64) error {
	if err := s.repo.SetActiveVersion(ctx, version); err != nil {
		return err
	}
	_ = s.auditRepo.CreateLog(ctx, &domain.AuditLog{
		ActorID:    optionalActorID(actorID),
		Action:     "rollback_rule_version",
		TargetType: "rule_config",
		TargetID:   version,
		Detail:     map[string]interface{}{},
	})
	return nil
}

// SetGrayMode 设置灰度开关
func (s *RuleConfigService) SetGrayMode(ctx context.Context, version string, enabled bool, actorID int64) error {
	if err := s.repo.SetGrayMode(ctx, version, enabled); err != nil {
		return err
	}
	_ = s.auditRepo.CreateLog(ctx, &domain.AuditLog{
		ActorID:    optionalActorID(actorID),
		Action:     "set_rule_gray_mode",
		TargetType: "rule_config",
		TargetID:   version,
		Detail: map[string]interface{}{
			"enabled": enabled,
		},
	})
	return nil
}

// GetActive 获取生效版本
func (s *RuleConfigService) GetActive(ctx context.Context) (*domain.RuleConfigVersion, error) {
	return s.repo.GetActiveVersion(ctx)
}

// ListVersions 查询版本列表
func (s *RuleConfigService) ListVersions(ctx context.Context, limit int) ([]domain.RuleConfigVersion, error) {
	return s.repo.ListVersions(ctx, limit)
}
