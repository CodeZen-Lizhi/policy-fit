package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zhenglizhi/policy-fit/internal/domain"
)

// FindingRepository 风险发现仓储
type FindingRepository struct {
	db *sql.DB
}

// NewFindingRepository 创建风险发现仓储
func NewFindingRepository(db *sql.DB) *FindingRepository {
	return &FindingRepository{db: db}
}

// BatchCreateFindings 批量创建风险发现
func (r *FindingRepository) BatchCreateFindings(ctx context.Context, findings []domain.RiskFinding) error {
	if len(findings) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	const query = `
INSERT INTO risk_finding (
    task_id, level, topic, summary, health_evidence, policy_evidence, questions, actions, confidence
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, created_at`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare insert finding: %w", err)
	}
	defer stmt.Close()

	for i := range findings {
		finding := &findings[i]
		if finding.TaskID == 0 {
			return errors.New("task_id is required")
		}
		if finding.Level == "" {
			return errors.New("level is required")
		}
		if finding.Topic == "" {
			return errors.New("topic is required")
		}
		if finding.Summary == "" {
			return errors.New("summary is required")
		}

		healthEvidenceJSON, err := json.Marshal(finding.HealthEvidence)
		if err != nil {
			return fmt.Errorf("marshal health evidence: %w", err)
		}
		policyEvidenceJSON, err := json.Marshal(finding.PolicyEvidence)
		if err != nil {
			return fmt.Errorf("marshal policy evidence: %w", err)
		}
		questionsJSON, err := json.Marshal(finding.Questions)
		if err != nil {
			return fmt.Errorf("marshal questions: %w", err)
		}
		actionsJSON, err := marshalJSON(finding.Actions)
		if err != nil {
			return fmt.Errorf("marshal actions: %w", err)
		}

		if err := stmt.QueryRowContext(
			ctx,
			finding.TaskID,
			finding.Level,
			finding.Topic,
			finding.Summary,
			healthEvidenceJSON,
			policyEvidenceJSON,
			questionsJSON,
			actionsJSON,
			finding.Confidence,
		).Scan(&finding.ID, &finding.CreatedAt); err != nil {
			return fmt.Errorf("insert finding: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

// ListByTask 按任务查询风险发现
func (r *FindingRepository) ListByTask(ctx context.Context, taskID int64) ([]domain.RiskFinding, error) {
	const query = `
SELECT id, task_id, level, topic, summary, health_evidence, policy_evidence, questions, actions, confidence, created_at
FROM risk_finding
WHERE task_id = $1
ORDER BY id ASC`

	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("query findings by task: %w", err)
	}
	defer rows.Close()

	findings := make([]domain.RiskFinding, 0)
	for rows.Next() {
		var (
			finding           domain.RiskFinding
			healthEvidenceRaw []byte
			policyEvidenceRaw []byte
			questionsRaw      []byte
			actionsRaw        []byte
			confidence        sql.NullFloat64
		)

		if err := rows.Scan(
			&finding.ID,
			&finding.TaskID,
			&finding.Level,
			&finding.Topic,
			&finding.Summary,
			&healthEvidenceRaw,
			&policyEvidenceRaw,
			&questionsRaw,
			&actionsRaw,
			&confidence,
			&finding.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan finding: %w", err)
		}

		if err := json.Unmarshal(healthEvidenceRaw, &finding.HealthEvidence); err != nil {
			return nil, fmt.Errorf("unmarshal health evidence: %w", err)
		}
		if err := json.Unmarshal(policyEvidenceRaw, &finding.PolicyEvidence); err != nil {
			return nil, fmt.Errorf("unmarshal policy evidence: %w", err)
		}
		if err := json.Unmarshal(questionsRaw, &finding.Questions); err != nil {
			return nil, fmt.Errorf("unmarshal questions: %w", err)
		}
		if len(actionsRaw) > 0 {
			if err := json.Unmarshal(actionsRaw, &finding.Actions); err != nil {
				return nil, fmt.Errorf("unmarshal actions: %w", err)
			}
		}
		if confidence.Valid {
			finding.Confidence = confidence.Float64
		}

		findings = append(findings, finding)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate findings: %w", err)
	}

	return findings, nil
}
