package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/zhenglizhi/policy-fit/internal/domain"
)

func TestTaskRepositoryCreateTask(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewTaskRepository(db)
	now := time.Now()

	task := &domain.AnalysisTask{
		UserID:      101,
		RequestID:   "req-1",
		Status:      domain.TaskStatusPending,
		RiskSummary: map[string]int{"red": 1},
	}

	mock.ExpectQuery(mustMatchSQL(`
INSERT INTO analysis_task (user_id, request_id, status, risk_summary)
VALUES ($1, NULLIF($2, ''), $3, $4)
RETURNING id, created_at, updated_at`)).
		WithArgs(task.UserID, task.RequestID, task.Status, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(1, now, now))

	if err := repo.CreateTask(context.Background(), task); err != nil {
		t.Fatalf("CreateTask error: %v", err)
	}

	if task.ID != 1 {
		t.Fatalf("unexpected task id: %d", task.ID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestTaskRepositoryGetTask(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewTaskRepository(db)
	now := time.Now()

	mock.ExpectQuery(mustMatchSQL(`
SELECT id, user_id, request_id, status, risk_summary, created_at, updated_at
FROM analysis_task
WHERE id = $1`)).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "user_id", "request_id", "status", "risk_summary", "created_at", "updated_at"},
		).AddRow(1, 1001, "req-1", "pending", `{"red":1,"yellow":2}`, now, now))

	task, err := repo.GetTask(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetTask error: %v", err)
	}
	if task.UserID != 1001 {
		t.Fatalf("unexpected user id: %d", task.UserID)
	}
	if task.RiskSummary["yellow"] != 2 {
		t.Fatalf("unexpected risk summary: %#v", task.RiskSummary)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestTaskRepositoryUpdateTaskStatus(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewTaskRepository(db)

	mock.ExpectExec(mustMatchSQL(`
UPDATE analysis_task
SET status = $2,
    risk_summary = COALESCE($3::jsonb, risk_summary)
WHERE id = $1`)).
		WithArgs(int64(1), domain.TaskStatusSuccess, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.UpdateTaskStatus(
		context.Background(),
		1,
		domain.TaskStatusSuccess,
		map[string]int{"green": 1},
	); err != nil {
		t.Fatalf("UpdateTaskStatus error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestTaskRepositoryDeleteTaskNotFound(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewTaskRepository(db)

	mock.ExpectExec(mustMatchSQL(`DELETE FROM analysis_task WHERE id = $1`)).
		WithArgs(int64(999)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeleteTask(context.Background(), 999)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
