package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/zhenglizhi/policy-fit/internal/domain"
)

func TestDocumentRepositoryCreateDocument(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewDocumentRepository(db)
	now := time.Now()

	doc := &domain.Document{
		TaskID:      1,
		DocType:     domain.DocTypeReport,
		FileName:    "report.pdf",
		StorageKey:  "task/1/report/file.pdf",
		ParseStatus: domain.ParseStatusPending,
	}

	mock.ExpectQuery(mustMatchSQL(`
INSERT INTO document (task_id, doc_type, file_name, storage_key, parse_status, parsed_text)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at`)).
		WithArgs(doc.TaskID, doc.DocType, doc.FileName, doc.StorageKey, doc.ParseStatus, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(11, now))

	if err := repo.CreateDocument(context.Background(), doc); err != nil {
		t.Fatalf("CreateDocument error: %v", err)
	}
	if doc.ID != 11 {
		t.Fatalf("unexpected id: %d", doc.ID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDocumentRepositoryListByTask(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewDocumentRepository(db)
	now := time.Now()

	mock.ExpectQuery(mustMatchSQL(`
SELECT id, task_id, doc_type, file_name, storage_key, parse_status, parsed_text, created_at
FROM document
WHERE task_id = $1
ORDER BY id ASC`)).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "task_id", "doc_type", "file_name", "storage_key", "parse_status", "parsed_text", "created_at"},
		).
			AddRow(11, 1, "report", "a.pdf", "k1", "pending", nil, now).
			AddRow(12, 1, "policy", "b.pdf", "k2", "success", "parsed", now))

	docs, err := repo.ListByTask(context.Background(), 1)
	if err != nil {
		t.Fatalf("ListByTask error: %v", err)
	}
	if len(docs) != 2 {
		t.Fatalf("unexpected docs length: %d", len(docs))
	}
	if docs[1].ParsedText != "parsed" {
		t.Fatalf("unexpected parsed text: %s", docs[1].ParsedText)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDocumentRepositoryUpdateParseStatus(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := NewDocumentRepository(db)

	mock.ExpectExec(mustMatchSQL(`
UPDATE document
SET parse_status = $2,
    parsed_text = $3
WHERE id = $1`)).
		WithArgs(int64(11), domain.ParseStatusSuccess, "parsed body").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.UpdateParseStatus(context.Background(), 11, domain.ParseStatusSuccess, "parsed body"); err != nil {
		t.Fatalf("UpdateParseStatus error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
