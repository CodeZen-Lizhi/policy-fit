package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zhenglizhi/policy-fit/internal/config"
	"github.com/zhenglizhi/policy-fit/internal/jobs"
	"github.com/zhenglizhi/policy-fit/pkg/logger"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger.Init(cfg.Log.Level, cfg.Log.Format)
	defer logger.Sync()

	// 创建 Worker
	worker := jobs.NewWorker(cfg)

	// 启动 Worker
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		logger.Info("Starting worker", "concurrency", cfg.Worker.Concurrency)
		if err := worker.Start(ctx); err != nil {
			logger.Fatal("Worker failed", "error", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down worker...")
	cancel()

	logger.Info("Worker exited")
}
