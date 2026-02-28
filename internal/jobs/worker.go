package jobs

import (
	"context"

	"github.com/zhenglizhi/policy-fit/internal/config"
	"github.com/zhenglizhi/policy-fit/pkg/logger"
)

// Worker 任务处理器
type Worker struct {
	cfg *config.Config
}

// NewWorker 创建 Worker
func NewWorker(cfg *config.Config) *Worker {
	return &Worker{
		cfg: cfg,
	}
}

// Start 启动 Worker
func (w *Worker) Start(ctx context.Context) error {
	logger.Info("Worker started")

	// TODO: 实现任务队列消费逻辑
	// 1. 连接 Redis
	// 2. 监听任务队列
	// 3. 处理任务（解析、抽取、匹配）

	<-ctx.Done()
	return nil
}
