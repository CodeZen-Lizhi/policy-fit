package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/zhenglizhi/policy-fit/internal/domain"
)

// RuleConfigRepository 规则配置仓储
type RuleConfigRepository struct {
	db *sql.DB
}

// NewRuleConfigRepository 创建规则配置仓储
func NewRuleConfigRepository(db *sql.DB) *RuleConfigRepository {
	return &RuleConfigRepository{db: db}
}

// CreateVersion 创建规则版本
func (r *RuleConfigRepository) CreateVersion(ctx context.Context, cfg *domain.RuleConfigVersion) error {
	content, err := marshalJSON(cfg.Content)
	if err != nil {
		return err
	}
	const query = `
INSERT INTO rule_config_version(version, changelog, content, is_active, is_gray, created_by)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at`
	return r.db.QueryRowContext(
		ctx,
		query,
		cfg.Version,
		cfg.Changelog,
		content,
		cfg.IsActive,
		cfg.IsGray,
		nullInt64(cfg.CreatedBy),
	).Scan(&cfg.ID, &cfg.CreatedAt)
}

// SetActiveVersion 激活规则版本
func (r *RuleConfigRepository) SetActiveVersion(ctx context.Context, version string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `UPDATE rule_config_version SET is_active = FALSE`); err != nil {
		return err
	}
	result, err := tx.ExecContext(ctx, `UPDATE rule_config_version SET is_active = TRUE WHERE version = $1`, version)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return tx.Commit()
}

// SetGrayMode 设置灰度开关
func (r *RuleConfigRepository) SetGrayMode(ctx context.Context, version string, isGray bool) error {
	result, err := r.db.ExecContext(ctx, `UPDATE rule_config_version SET is_gray = $2 WHERE version = $1`, version, isGray)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// GetActiveVersion 获取当前生效版本
func (r *RuleConfigRepository) GetActiveVersion(ctx context.Context) (*domain.RuleConfigVersion, error) {
	const query = `
SELECT id, version, changelog, content, is_active, is_gray, created_by, created_at
FROM rule_config_version
WHERE is_active = TRUE
ORDER BY id DESC
LIMIT 1`

	var (
		cfg         domain.RuleConfigVersion
		contentRaw  []byte
		createdByID sql.NullInt64
	)
	if err := r.db.QueryRowContext(ctx, query).Scan(
		&cfg.ID,
		&cfg.Version,
		&cfg.Changelog,
		&contentRaw,
		&cfg.IsActive,
		&cfg.IsGray,
		&createdByID,
		&cfg.CreatedAt,
	); err != nil {
		return nil, err
	}
	if createdByID.Valid {
		cfg.CreatedBy = &createdByID.Int64
	}
	if err := json.Unmarshal(contentRaw, &cfg.Content); err != nil {
		return nil, fmt.Errorf("unmarshal rule config content: %w", err)
	}
	return &cfg, nil
}

// ListVersions 查询规则版本列表
func (r *RuleConfigRepository) ListVersions(ctx context.Context, limit int) ([]domain.RuleConfigVersion, error) {
	if limit <= 0 {
		limit = 20
	}
	const query = `
SELECT id, version, changelog, content, is_active, is_gray, created_by, created_at
FROM rule_config_version
ORDER BY id DESC
LIMIT $1`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.RuleConfigVersion, 0, limit)
	for rows.Next() {
		var (
			cfg         domain.RuleConfigVersion
			contentRaw  []byte
			createdByID sql.NullInt64
		)
		if err := rows.Scan(
			&cfg.ID,
			&cfg.Version,
			&cfg.Changelog,
			&contentRaw,
			&cfg.IsActive,
			&cfg.IsGray,
			&createdByID,
			&cfg.CreatedAt,
		); err != nil {
			return nil, err
		}
		if createdByID.Valid {
			cfg.CreatedBy = &createdByID.Int64
		}
		if err := json.Unmarshal(contentRaw, &cfg.Content); err != nil {
			return nil, err
		}
		items = append(items, cfg)
	}
	return items, rows.Err()
}
