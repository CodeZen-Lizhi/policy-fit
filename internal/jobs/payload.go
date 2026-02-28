package jobs

import "time"

// TaskPayload 队列任务载荷
type TaskPayload struct {
	TaskID     int64     `json:"task_id"`
	RequestID  string    `json:"request_id,omitempty"`
	RetryCount int       `json:"retry_count"`
	EnqueuedAt time.Time `json:"enqueued_at"`
}
