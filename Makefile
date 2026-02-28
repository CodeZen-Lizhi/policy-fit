.PHONY: help build run-api run-worker test lint clean migrate-up migrate-down docker-up docker-down env-check

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## 编译项目
	@echo "Building..."
	@go build -o bin/api cmd/api/main.go
	@go build -o bin/worker cmd/worker/main.go
	@echo "Build complete!"

run-api: ## 运行 API 服务
	@echo "Starting API server..."
	@go run cmd/api/main.go

run-worker: ## 运行 Worker 服务
	@echo "Starting Worker..."
	@go run cmd/worker/main.go

test: ## 运行测试
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

lint: ## 代码检查
	@echo "Running linters..."
	@golangci-lint run ./...

clean: ## 清理构建产物
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

migrate-up: ## 执行数据库迁移（升级）
	@echo "Running migrations..."
	@go run cmd/migrate/main.go up

migrate-down: ## 回滚数据库迁移
	@echo "Rolling back migrations..."
	@go run cmd/migrate/main.go down

docker-up: ## 启动 Docker 容器
	@echo "Starting Docker containers..."
	@docker-compose up -d

docker-down: ## 停止 Docker 容器
	@echo "Stopping Docker containers..."
	@docker-compose down

docker-logs: ## 查看 Docker 日志
	@docker-compose logs -f

deps: ## 安装依赖
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

env-check: ## 校验当前环境配置
	@echo "Validating environment configuration..."
	@go run cmd/envcheck/main.go
	@echo "Environment configuration is valid."

.DEFAULT_GOAL := help
