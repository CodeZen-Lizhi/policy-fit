package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/zhenglizhi/policy-fit/internal/domain"
)

// DocumentRepository 文档仓储
type DocumentRepository struct {
	db *sql.DB
}

// NewDocumentRepository 创建文档仓储
func NewDocumentRepository(db *sql.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

// CreateDocument 创建文档
func (r *DocumentRepository) CreateDocument(ctx context.Context, doc *domain.Document) error {
	if doc == nil {
		return errors.New("document is nil")
	}
	if doc.ParseStatus == "" {
		doc.ParseStatus = domain.ParseStatusPending
	}

	const query = `
INSERT INTO document (task_id, doc_type, file_name, storage_key, parse_status, parsed_text)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at`

	if err := r.db.QueryRowContext(
		ctx,
		query,
		doc.TaskID,
		doc.DocType,
		doc.FileName,
		doc.StorageKey,
		doc.ParseStatus,
		nullString(doc.ParsedText),
	).Scan(&doc.ID, &doc.CreatedAt); err != nil {
		return fmt.Errorf("insert document: %w", err)
	}

	return nil
}

// ListByTask 按任务查询文档
func (r *DocumentRepository) ListByTask(ctx context.Context, taskID int64) ([]domain.Document, error) {
	const query = `
SELECT id, task_id, doc_type, file_name, storage_key, parse_status, parsed_text, created_at
FROM document
WHERE task_id = $1
ORDER BY id ASC`

	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("query documents by task: %w", err)
	}
	defer rows.Close()

	docs := make([]domain.Document, 0)
	for rows.Next() {
		var (
			doc        domain.Document
			parsedText sql.NullString
		)

		if err := rows.Scan(
			&doc.ID,
			&doc.TaskID,
			&doc.DocType,
			&doc.FileName,
			&doc.StorageKey,
			&doc.ParseStatus,
			&parsedText,
			&doc.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan document: %w", err)
		}

		if parsedText.Valid {
			doc.ParsedText = parsedText.String
		}
		docs = append(docs, doc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate documents: %w", err)
	}

	return docs, nil
}

// UpdateParseStatus 更新文档解析状态
func (r *DocumentRepository) UpdateParseStatus(
	ctx context.Context,
	id int64,
	status domain.ParseStatus,
	parsedText string,
) error {
	if status == "" {
		return errors.New("parse status is required")
	}

	const query = `
UPDATE document
SET parse_status = $2,
    parsed_text = $3
WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id, status, nullString(parsedText))
	if err != nil {
		return fmt.Errorf("update document parse status: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read affected rows: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func nullString(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}
