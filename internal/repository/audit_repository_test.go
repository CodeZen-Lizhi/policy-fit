package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/zhenglizhi/policy-fit/internal/domain"
)

func TestAuditRepositoryCreateLog(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewAuditRepository(db)
	now := time.Now()

	taskID := int64(1)
	actorID := int64(2)
	logEntry := &domain.AuditLog{
		TaskID:     &taskID,
		ActorID:    &actorID,
		Action:     "delete_task",
		TargetType: "task",
		TargetID:   "1",
		Detail: map[string]interface{}{
			"reason": "user requested",
		},
	}

	mock.ExpectQuery(mustMatchSQL(`
INSERT INTO audit_log (task_id, actor_id, action, target_type, target_id, detail)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at`)).
		WithArgs(
			taskID,
			actorID,
			"delete_task",
			"task",
			"1",
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(31, now))

	if err := repo.CreateLog(context.Background(), logEntry); err != nil {
		t.Fatalf("CreateLog error: %v", err)
	}
	if logEntry.ID != 31 {
		t.Fatalf("unexpected audit id: %d", logEntry.ID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
