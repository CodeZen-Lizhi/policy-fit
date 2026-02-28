package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/zhenglizhi/policy-fit/internal/config"
	"github.com/zhenglizhi/policy-fit/internal/domain"
	"github.com/zhenglizhi/policy-fit/pkg/logger"
)

func TestWorkerRequeueOnFailure(t *testing.T) {
	logger.Init("error", "json")
	defer logger.Sync()

	queue := newFakeQueue()
	updater := &fakeTaskStatusUpdater{}
	cfg := &config.Config{Worker: config.WorkerConfig{Concurrency: 1}}

	worker := NewWorker(cfg, queue, updater)
	worker.processor = ProcessorFunc(func(context.Context, TaskPayload) error {
		return errors.New("forced error")
	})

	payload := TaskPayload{
		TaskID:     101,
		RetryCount: 0,
	}
	raw, _ := json.Marshal(payload)
	queue.pushInput(string(raw))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go worker.consumeLoop(ctx, 1)
	time.Sleep(100 * time.Millisecond)
	cancel()
	time.Sleep(50 * time.Millisecond)

	requeued := queue.lastPush(DefaultQueueName)
	if requeued == "" {
		t.Fatalf("expected payload requeued")
	}

	var retryPayload TaskPayload
	if err := json.Unmarshal([]byte(requeued), &retryPayload); err != nil {
		t.Fatalf("unmarshal requeued payload: %v", err)
	}
	if retryPayload.RetryCount != 1 {
		t.Fatalf("expected retry_count=1, got %d", retryPayload.RetryCount)
	}

	if !updater.hasStatus(101, domain.TaskStatusFailed) {
		t.Fatalf("expected task status updated to failed")
	}
}

func TestWorkerMovesToDeadLetterAfterMaxRetries(t *testing.T) {
	logger.Init("error", "json")
	defer logger.Sync()

	queue := newFakeQueue()
	updater := &fakeTaskStatusUpdater{}
	cfg := &config.Config{Worker: config.WorkerConfig{Concurrency: 1}}

	worker := NewWorker(cfg, queue, updater)
	worker.processor = ProcessorFunc(func(context.Context, TaskPayload) error {
		return errors.New("forced error")
	})

	payload := TaskPayload{
		TaskID:     102,
		RetryCount: 2,
	}
	raw, _ := json.Marshal(payload)
	queue.pushInput(string(raw))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go worker.consumeLoop(ctx, 1)
	time.Sleep(100 * time.Millisecond)
	cancel()
	time.Sleep(50 * time.Millisecond)

	dead := queue.lastPush(DefaultDeadLetterQueue)
	if dead == "" {
		t.Fatalf("expected payload moved to dead letter queue")
	}
}

func TestWorkerTrackTaskCompletedEvent(t *testing.T) {
	logger.Init("error", "json")
	defer logger.Sync()

	queue := newFakeQueue()
	updater := &fakeTaskStatusUpdater{}
	tracker := &fakeAnalyticsTracker{}
	cfg := &config.Config{Worker: config.WorkerConfig{Concurrency: 1}}

	worker := NewWorker(cfg, queue, updater)
	worker.SetAnalyticsTracker(tracker)
	worker.processor = ProcessorFunc(func(context.Context, TaskPayload) error {
		return nil
	})

	payload := TaskPayload{
		TaskID:     103,
		RetryCount: 1,
	}
	if err := worker.processTask(context.Background(), payload); err != nil {
		t.Fatalf("process task: %v", err)
	}
	if tracker.lastEvent != "task_completed" {
		t.Fatalf("expected task_completed event, got %s", tracker.lastEvent)
	}
}

type fakeQueue struct {
	mu      sync.Mutex
	input   chan string
	outputs map[string][]string
}

func newFakeQueue() *fakeQueue {
	return &fakeQueue{
		input:   make(chan string, 16),
		outputs: make(map[string][]string),
	}
}

func (q *fakeQueue) pushInput(payload string) {
	q.input <- payload
}

func (q *fakeQueue) Pop(ctx context.Context, _ string, _ time.Duration) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case payload := <-q.input:
		return payload, nil
	}
}

func (q *fakeQueue) Push(_ context.Context, queueName string, payload string) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.outputs[queueName] = append(q.outputs[queueName], payload)
	return nil
}

func (q *fakeQueue) lastPush(queueName string) string {
	q.mu.Lock()
	defer q.mu.Unlock()
	items := q.outputs[queueName]
	if len(items) == 0 {
		return ""
	}
	return items[len(items)-1]
}

type fakeTaskStatusUpdater struct {
	mu       sync.Mutex
	statuses map[int64][]domain.TaskStatus
}

type fakeAnalyticsTracker struct {
	lastEvent string
}

func (f *fakeAnalyticsTracker) Track(
	_ context.Context,
	eventName string,
	_ *int64,
	_ *int64,
	_ map[string]interface{},
) error {
	f.lastEvent = eventName
	return nil
}

func (u *fakeTaskStatusUpdater) UpdateTaskStatus(_ context.Context, id int64, status domain.TaskStatus, _ map[string]int) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.statuses == nil {
		u.statuses = make(map[int64][]domain.TaskStatus)
	}
	u.statuses[id] = append(u.statuses[id], status)
	return nil
}

func (u *fakeTaskStatusUpdater) hasStatus(taskID int64, target domain.TaskStatus) bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, s := range u.statuses[taskID] {
		if s == target {
			return true
		}
	}
	return false
}
