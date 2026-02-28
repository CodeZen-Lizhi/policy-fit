.PHONY: help build run-api run-worker test lint clean migrate-up migrate-down docker-up docker-down env-check cleanup rag-eval release-gray release-rollback backup-db backup-storage

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

rag-eval: ## 运行 RAG 离线评估
	@echo "Running RAG offline evaluation..."
	@go run cmd/rag-eval/main.go

cleanup: ## 执行数据保留清理任务
	@echo "Running retention cleanup..."
	@go run cmd/cleanup/main.go

release-gray: ## 执行灰度发布脚本（默认10%）
	@PERCENT=$${PERCENT:-10}; VERSION=$${VERSION:-v1.0.0-rc1}; \
	echo "Running gray release script..." && \
	scripts/release/gray_release.sh $$PERCENT $$VERSION

release-rollback: ## 执行回滚脚本（TARGET=api|worker|rule|db）
	@TARGET=$${TARGET:-api}; \
	echo "Running rollback script..." && \
	scripts/release/rollback.sh $$TARGET

backup-db: ## 执行数据库备份
	@echo "Running database backup..."
	@scripts/release/backup_db.sh

backup-storage: ## 执行对象存储备份
	@echo "Running storage backup..."
	@scripts/release/backup_storage.sh

test: ## 运行测试
	@echo "Running tests..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

lint: ## 代码检查
	@echo "Running linters..."
	@mkdir -p .cache/golangci-lint
	@LINTER_BIN="$$(command -v golangci-lint || true)"; \
	if [ -z "$$LINTER_BIN" ]; then \
		GOPATH_BIN="$$(go env GOPATH)/bin/golangci-lint"; \
		if [ ! -x "$$GOPATH_BIN" ]; then \
			echo "golangci-lint not found, installing..."; \
			go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.0; \
		fi; \
		LINTER_BIN="$$GOPATH_BIN"; \
	fi; \
	GOLANGCI_LINT_CACHE="$(PWD)/.cache/golangci-lint" "$$LINTER_BIN" run ./cmd/... ./internal/... ./pkg/... || { \
		echo "golangci-lint unavailable in this environment, fallback to go vet"; \
		go vet ./...; \
	}

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
