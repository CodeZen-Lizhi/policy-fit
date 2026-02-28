package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// DefaultQueueName 默认任务队列名
	DefaultQueueName = "analysis_tasks"
)

// Enqueuer 入队接口
type Enqueuer interface {
	EnqueueTask(ctx context.Context, payload TaskPayload) error
}

// RedisQueue Redis 队列
type RedisQueue struct {
	client    *redis.Client
	queueName string
}

// NewRedisQueue 创建 Redis 队列客户端
func NewRedisQueue(client *redis.Client, queueName string) *RedisQueue {
	if queueName == "" {
		queueName = DefaultQueueName
	}
	return &RedisQueue{
		client:    client,
		queueName: queueName,
	}
}

// EnqueueTask 入队任务
func (q *RedisQueue) EnqueueTask(ctx context.Context, payload TaskPayload) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	if err := q.Push(ctx, q.queueName, string(raw)); err != nil {
		return err
	}

	return nil
}

// Pop 从队列弹出任务
func (q *RedisQueue) Pop(ctx context.Context, queueName string, timeout time.Duration) (string, error) {
	result, err := q.client.BLPop(ctx, timeout, queueName).Result()
	if err != nil {
		return "", err
	}
	if len(result) < 2 {
		return "", fmt.Errorf("invalid blpop result")
	}
	return result[1], nil
}

// Push 推送任务到队列
func (q *RedisQueue) Push(ctx context.Context, queueName string, payload string) error {
	if err := q.client.RPush(ctx, queueName, payload).Err(); err != nil {
		return fmt.Errorf("push payload to queue: %w", err)
	}

	return nil
}
