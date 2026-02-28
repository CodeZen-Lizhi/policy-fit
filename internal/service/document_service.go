package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/zhenglizhi/policy-fit/internal/domain"
	"github.com/zhenglizhi/policy-fit/internal/repository"
)

// DocumentService 文档服务
type DocumentService struct {
	taskRepo     *repository.TaskRepository
	documentRepo *repository.DocumentRepository
}

// NewDocumentService 创建文档服务
func NewDocumentService(taskRepo *repository.TaskRepository, documentRepo *repository.DocumentRepository) *DocumentService {
	return &DocumentService{
		taskRepo:     taskRepo,
		documentRepo: documentRepo,
	}
}

// CreateDocument 创建文档元数据
func (s *DocumentService) CreateDocument(
	ctx context.Context,
	taskID int64,
	actorID int64,
	docType domain.DocumentType,
	fileName string,
	storageKey string,
) (*domain.Document, error) {
	if taskID <= 0 || docType == "" || fileName == "" || storageKey == "" {
		return nil, fmt.Errorf("%w: invalid document payload", ErrInvalidArgument)
	}

	task, err := s.taskRepo.GetTask(ctx, taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	if actorID > 0 && task.UserID != actorID {
		return nil, ErrForbidden
	}

	doc := &domain.Document{
		TaskID:      taskID,
		DocType:     docType,
		FileName:    fileName,
		StorageKey:  storageKey,
		ParseStatus: domain.ParseStatusPending,
	}
	if err := s.documentRepo.CreateDocument(ctx, doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// ListByTask 按任务查询文档
func (s *DocumentService) ListByTask(ctx context.Context, taskID int64) ([]domain.Document, error) {
	if taskID <= 0 {
		return nil, fmt.Errorf("%w: invalid task id", ErrInvalidArgument)
	}
	return s.documentRepo.ListByTask(ctx, taskID)
}
