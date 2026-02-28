package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/zhenglizhi/policy-fit/internal/domain"
)

func TestFindingRepositoryBatchCreateFindings(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewFindingRepository(db)
	now := time.Now()

	findings := []domain.RiskFinding{
		{
			TaskID:         1,
			Level:          domain.RiskLevelRed,
			Topic:          "hypertension",
			Summary:        "possible conflict",
			HealthEvidence: []domain.Evidence{{Loc: "para_1", Text: "bp 150/95"}},
			PolicyEvidence: []domain.Evidence{{Loc: "para_2", Text: "preexisting"}},
			Questions:      []string{"q1", "q2"},
			Actions:        []string{"a1"},
			Confidence:     0.92,
		},
	}

	mock.ExpectBegin()
	prep := mock.ExpectPrepare(mustMatchSQL(`
INSERT INTO risk_finding (
    task_id, level, topic, summary, health_evidence, policy_evidence, questions, actions, confidence
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, created_at`))
	prep.ExpectQuery().
		WithArgs(
			int64(1),
			domain.RiskLevelRed,
			"hypertension",
			"possible conflict",
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			0.92,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(21, now))
	mock.ExpectCommit()

	if err := repo.BatchCreateFindings(context.Background(), findings); err != nil {
		t.Fatalf("BatchCreateFindings error: %v", err)
	}
	if findings[0].ID != 21 {
		t.Fatalf("unexpected finding id: %d", findings[0].ID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestFindingRepositoryListByTask(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewFindingRepository(db)
	now := time.Now()

	mock.ExpectQuery(mustMatchSQL(`
SELECT id, task_id, level, topic, summary, health_evidence, policy_evidence, questions, actions, confidence, created_at
FROM risk_finding
WHERE task_id = $1
ORDER BY id ASC`)).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "task_id", "level", "topic", "summary", "health_evidence", "policy_evidence", "questions", "actions", "confidence", "created_at"},
		).AddRow(
			21,
			1,
			"yellow",
			"dyslipidemia",
			"need confirm",
			`[{"loc":"para_1","text":"tc high"}]`,
			`[{"loc":"para_2","text":"clause"}]`,
			`["q1","q2"]`,
			`["a1"]`,
			0.85,
			now,
		))

	items, err := repo.ListByTask(context.Background(), 1)
	if err != nil {
		t.Fatalf("ListByTask error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("unexpected findings length: %d", len(items))
	}
	if items[0].Topic != "dyslipidemia" {
		t.Fatalf("unexpected topic: %s", items[0].Topic)
	}
	if len(items[0].Actions) != 1 {
		t.Fatalf("unexpected actions: %#v", items[0].Actions)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
