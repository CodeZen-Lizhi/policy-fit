package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/zhenglizhi/policy-fit/internal/config"
	"github.com/zhenglizhi/policy-fit/internal/domain"
	"github.com/zhenglizhi/policy-fit/pkg/logger"
)

const (
	// DefaultDeadLetterQueue 默认死信队列名
	DefaultDeadLetterQueue = "analysis_tasks_dead_letter"
	maxRetries             = 2
)

// Queue 队列抽象
type Queue interface {
	Pop(ctx context.Context, queueName string, timeout time.Duration) (string, error)
	Push(ctx context.Context, queueName string, payload string) error
}

// TaskStatusUpdater 任务状态更新抽象
type TaskStatusUpdater interface {
	UpdateTaskStatus(ctx context.Context, id int64, status domain.TaskStatus, riskSummary map[string]int) error
}

// Processor 任务处理抽象
type Processor interface {
	Process(ctx context.Context, payload TaskPayload) error
}

// AnalyticsTracker 埋点写入抽象
type AnalyticsTracker interface {
	Track(
		ctx context.Context,
		eventName string,
		userID *int64,
		taskID *int64,
		properties map[string]interface{},
	) error
}

// ProcessorFunc 任务处理函数适配器
type ProcessorFunc func(ctx context.Context, payload TaskPayload) error

// Process 执行处理函数
func (f ProcessorFunc) Process(ctx context.Context, payload TaskPayload) error {
	return f(ctx, payload)
}

// Worker 任务处理器
type Worker struct {
	cfg             *config.Config
	queue           Queue
	statusUpdater   TaskStatusUpdater
	processor       Processor
	analytics       AnalyticsTracker
	queueName       string
	deadLetterQueue string
}

// NewWorker 创建 Worker
func NewWorker(cfg *config.Config, queue Queue, statusUpdater TaskStatusUpdater) *Worker {
	return &Worker{
		cfg:             cfg,
		queue:           queue,
		statusUpdater:   statusUpdater,
		processor:       ProcessorFunc(func(context.Context, TaskPayload) error { return nil }),
		queueName:       DefaultQueueName,
		deadLetterQueue: DefaultDeadLetterQueue,
	}
}

// SetAnalyticsTracker 设置埋点追踪器
func (w *Worker) SetAnalyticsTracker(tracker AnalyticsTracker) {
	w.analytics = tracker
}

// Start 启动 Worker
func (w *Worker) Start(ctx context.Context) error {
	logger.Info("Worker started", "concurrency", w.cfg.Worker.Concurrency)

	concurrency := w.cfg.Worker.Concurrency
	if concurrency <= 0 {
		concurrency = 1
	}

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		workerID := i + 1
		wg.Add(1)
		go func() {
			defer wg.Done()
			w.consumeLoop(ctx, workerID)
		}()
	}

	<-ctx.Done()
	logger.Info("Worker shutting down, waiting for running tasks to finish")
	wg.Wait()
	return nil
}

func (w *Worker) consumeLoop(ctx context.Context, workerID int) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("Worker loop stopped", "worker_id", workerID)
			return
		default:
		}

		rawPayload, err := w.queue.Pop(ctx, w.queueName, 5*time.Second)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				continue
			}
			logger.Error("Failed to pop task from queue", "worker_id", workerID, "error", err)
			continue
		}

		var payload TaskPayload
		if err := json.Unmarshal([]byte(rawPayload), &payload); err != nil {
			logger.Error("Invalid queue payload", "worker_id", workerID, "error", err)
			continue
		}

		if err := w.processTask(ctx, payload); err != nil {
			logger.Error("Task processing failed", "worker_id", workerID, "task_id", payload.TaskID, "error", err)
			w.handleFailure(ctx, payload, err)
		}
	}
}

func (w *Worker) processTask(ctx context.Context, payload TaskPayload) error {
	if payload.TaskID <= 0 {
		return fmt.Errorf("invalid task id: %d", payload.TaskID)
	}

	phases := []domain.TaskStatus{
		domain.TaskStatusParsing,
		domain.TaskStatusExtracting,
		domain.TaskStatusMatching,
	}
	for _, phase := range phases {
		logger.Info("Task phase start", "task_id", payload.TaskID, "phase", phase)
		if err := w.statusUpdater.UpdateTaskStatus(ctx, payload.TaskID, phase, nil); err != nil {
			return fmt.Errorf("update status to %s: %w", phase, err)
		}
	}

	if err := w.processor.Process(ctx, payload); err != nil {
		return err
	}

	if err := w.statusUpdater.UpdateTaskStatus(ctx, payload.TaskID, domain.TaskStatusSuccess, map[string]int{
		"red":    0,
		"yellow": 0,
		"green":  0,
	}); err != nil {
		return fmt.Errorf("update status to success: %w", err)
	}
	logger.Info("Task phase completed", "task_id", payload.TaskID, "phase", domain.TaskStatusSuccess)
	w.trackEvent(ctx, "task_completed", payload.TaskID, map[string]interface{}{
		"retry_count": payload.RetryCount,
	})

	return nil
}

func (w *Worker) handleFailure(ctx context.Context, payload TaskPayload, processErr error) {
	_ = w.statusUpdater.UpdateTaskStatus(ctx, payload.TaskID, domain.TaskStatusFailed, nil)
	w.trackEvent(ctx, "task_failed", payload.TaskID, map[string]interface{}{
		"retry_count": payload.RetryCount,
		"error":       processErr.Error(),
	})

	raw, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal failed payload", "task_id", payload.TaskID, "error", err)
		return
	}

	if payload.RetryCount < maxRetries {
		payload.RetryCount++
		retryRaw, _ := json.Marshal(payload)
		if err := w.queue.Push(ctx, w.queueName, string(retryRaw)); err != nil {
			logger.Error("Failed to requeue task", "task_id", payload.TaskID, "retry_count", payload.RetryCount, "error", err)
		}
		return
	}

	if err := w.queue.Push(ctx, w.deadLetterQueue, string(raw)); err != nil {
		logger.Error("Failed to push task to dead letter queue", "task_id", payload.TaskID, "error", err)
	}

	logger.Warn("Task moved to dead letter queue", "task_id", payload.TaskID, "error", processErr)
}

func (w *Worker) trackEvent(ctx context.Context, eventName string, taskID int64, properties map[string]interface{}) {
	if w.analytics == nil || taskID <= 0 {
		return
	}
	tid := taskID
	_ = w.analytics.Track(ctx, eventName, nil, &tid, properties)
}
