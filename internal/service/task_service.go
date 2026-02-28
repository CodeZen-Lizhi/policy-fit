package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/zhenglizhi/policy-fit/internal/domain"
	"github.com/zhenglizhi/policy-fit/internal/jobs"
	"github.com/zhenglizhi/policy-fit/internal/repository"
)

// TaskService 任务服务
type TaskService struct {
	taskRepo     *repository.TaskRepository
	documentRepo *repository.DocumentRepository
	auditRepo    *repository.AuditRepository
	storage      objectStorage
	queue        jobs.Enqueuer
}

type objectStorage interface {
	Delete(ctx context.Context, storageKey string) error
}

// NewTaskService 创建任务服务
func NewTaskService(
	taskRepo *repository.TaskRepository,
	documentRepo *repository.DocumentRepository,
	auditRepo *repository.AuditRepository,
	storage objectStorage,
	queue jobs.Enqueuer,
) *TaskService {
	return &TaskService{
		taskRepo:     taskRepo,
		documentRepo: documentRepo,
		auditRepo:    auditRepo,
		storage:      storage,
		queue:        queue,
	}
}

// CreateTask 创建任务
func (s *TaskService) CreateTask(ctx context.Context, userID int64, requestID string) (*domain.AnalysisTask, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("%w: user_id must be positive", ErrInvalidArgument)
	}

	task := &domain.AnalysisTask{
		UserID:      userID,
		RequestID:   requestID,
		Status:      domain.TaskStatusPending,
		RiskSummary: map[string]int{"red": 0, "yellow": 0, "green": 0},
	}
	if err := s.taskRepo.CreateTask(ctx, task); err != nil {
		return nil, err
	}

	_ = s.auditRepo.CreateLog(ctx, &domain.AuditLog{
		TaskID:     &task.ID,
		ActorID:    &userID,
		Action:     "create_task",
		TargetType: "analysis_task",
		TargetID:   fmt.Sprintf("%d", task.ID),
		Detail: map[string]interface{}{
			"request_id": requestID,
		},
	})

	return task, nil
}

// GetTask 获取任务
func (s *TaskService) GetTask(ctx context.Context, taskID int64, actorID int64) (*domain.AnalysisTask, error) {
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
	return task, nil
}

// RunTask 运行任务（入队）
func (s *TaskService) RunTask(ctx context.Context, taskID int64, actorID int64) error {
	task, err := s.GetTask(ctx, taskID, actorID)
	if err != nil {
		return err
	}

	if task.Status != domain.TaskStatusPending && task.Status != domain.TaskStatusFailed {
		return fmt.Errorf("%w: current status is %s", ErrTaskStatusConflict, task.Status)
	}

	docs, err := s.documentRepo.ListByTask(ctx, taskID)
	if err != nil {
		return err
	}
	if !hasRequiredDocuments(docs) {
		return ErrRequiredDocumentsMissing
	}

	if err := s.queue.EnqueueTask(ctx, jobs.TaskPayload{
		TaskID:     taskID,
		RequestID:  task.RequestID,
		RetryCount: 0,
		EnqueuedAt: time.Now(),
	}); err != nil {
		return err
	}

	if err := s.taskRepo.UpdateTaskStatus(ctx, taskID, domain.TaskStatusPending, nil); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTaskNotFound
		}
		return err
	}

	_ = s.auditRepo.CreateLog(ctx, &domain.AuditLog{
		TaskID:     &taskID,
		ActorID:    optionalActorID(actorID),
		Action:     "run_task",
		TargetType: "analysis_task",
		TargetID:   fmt.Sprintf("%d", taskID),
		Detail: map[string]interface{}{
			"status": string(task.Status),
		},
	})

	return nil
}

// DeleteTask 删除任务（文档和 findings 由数据库级联删除）
func (s *TaskService) DeleteTask(ctx context.Context, taskID int64, actorID int64) error {
	task, err := s.GetTask(ctx, taskID, actorID)
	if err != nil {
		return err
	}

	docs, err := s.documentRepo.ListByTask(ctx, taskID)
	if err != nil {
		return err
	}
	for _, doc := range docs {
		if doc.StorageKey == "" {
			continue
		}
		if err := s.storage.Delete(ctx, doc.StorageKey); err != nil {
			return fmt.Errorf("delete object %s: %w", doc.StorageKey, err)
		}
	}

	if err := s.taskRepo.DeleteTask(ctx, task.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTaskNotFound
		}
		return err
	}

	_ = s.auditRepo.CreateLog(ctx, &domain.AuditLog{
		TaskID:     &taskID,
		ActorID:    optionalActorID(actorID),
		Action:     "delete_task",
		TargetType: "analysis_task",
		TargetID:   fmt.Sprintf("%d", taskID),
		Detail:     map[string]interface{}{},
	})

	return nil
}

func hasRequiredDocuments(docs []domain.Document) bool {
	hasReport := false
	hasPolicy := false

	for _, doc := range docs {
		if doc.DocType == domain.DocTypeReport {
			hasReport = true
		}
		if doc.DocType == domain.DocTypePolicy {
			hasPolicy = true
		}
	}
	return hasReport && hasPolicy
}

func optionalActorID(actorID int64) *int64 {
	if actorID <= 0 {
		return nil
	}
	return &actorID
}
