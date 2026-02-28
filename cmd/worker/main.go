package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/zhenglizhi/policy-fit/internal/config"
	"github.com/zhenglizhi/policy-fit/internal/jobs"
	"github.com/zhenglizhi/policy-fit/internal/repository"
	"github.com/zhenglizhi/policy-fit/internal/service"
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

	db, err := openPostgres(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer db.Close()

	redisClient, err := openRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to connect redis: %v", err)
	}
	defer redisClient.Close()

	taskRepo := repository.NewTaskRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)
	queue := jobs.NewRedisQueue(redisClient, jobs.DefaultQueueName)
	analyticsSvc := service.NewAnalyticsService(analyticsRepo)

	// 创建 Worker
	worker := jobs.NewWorker(cfg, queue, taskRepo)
	worker.SetAnalyticsTracker(analyticsSvc)

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

func openPostgres(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func openRedis(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}
	return client, nil
}
