package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/zhenglizhi/policy-fit/internal/config"
	"github.com/zhenglizhi/policy-fit/internal/handler"
	"github.com/zhenglizhi/policy-fit/internal/jobs"
	"github.com/zhenglizhi/policy-fit/internal/middleware"
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

	// 初始化数据库
	db, err := openPostgres(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer db.Close()

	// 初始化 Redis
	redisClient, err := openRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to connect redis: %v", err)
	}
	defer redisClient.Close()

	// 初始化仓储与服务
	taskRepo := repository.NewTaskRepository(db)
	documentRepo := repository.NewDocumentRepository(db)
	findingRepo := repository.NewFindingRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)
	ruleConfigRepo := repository.NewRuleConfigRepository(db)

	queue := jobs.NewRedisQueue(redisClient, jobs.DefaultQueueName)

	storageSvc, err := service.NewStorageService(cfg.Storage)
	if err != nil {
		log.Fatalf("Failed to init storage service: %v", err)
	}
	taskSvc := service.NewTaskService(taskRepo, documentRepo, auditRepo, storageSvc, queue)
	documentSvc := service.NewDocumentService(taskRepo, documentRepo)
	findingSvc := service.NewFindingService(findingRepo)
	exportSvc := service.NewExportService(taskSvc, findingSvc)
	compareSvc := service.NewComparisonService(taskSvc, findingSvc)
	analyticsSvc := service.NewAnalyticsService(analyticsRepo)
	ruleConfigSvc := service.NewRuleConfigService(ruleConfigRepo, auditRepo)

	taskHandler := handler.NewTaskHandler(taskSvc, documentSvc, findingSvc, storageSvc, exportSvc, compareSvc, analyticsSvc)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsSvc, cfg.Security.AdminUserIDs)
	ruleConfigHandler := handler.NewRuleConfigHandler(ruleConfigSvc, cfg.Security.AdminUserIDs)

	// 初始化 Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.Metrics())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "not_ready",
				"message": "database unavailable",
			})
			return
		}
		if err := redisClient.Ping(ctx).Err(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "not_ready",
				"message": "redis unavailable",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})
	router.GET("/metrics", middleware.MetricsSnapshotHandler())

	// API 路由
	v1 := router.Group("/api/v1")
	v1.Use(middleware.Auth(cfg.Security.JWTSecret))
	{
		handler.RegisterTaskRoutes(v1, taskHandler)
		handler.RegisterAnalyticsRoutes(v1, analyticsHandler)
		handler.RegisterRuleConfigRoutes(v1, ruleConfigHandler)
	}

	// 启动服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// 优雅关闭
	go func() {
		logger.Info("Starting API server", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited")
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
