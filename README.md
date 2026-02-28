# 保单避坑雷达（Policy × Health Fit）

> 用户上传体检报告与保险条款，系统输出"潜在拒赔/除外风险"与"条款-证据对照解释"，帮助用户在投保前后识别风险盲区。

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://go.dev/)

## ✨ 功能特性

- 🔴 **红黄绿风险分级**：智能识别高风险、待确认、暂无冲突三类风险
- 📋 **证据对照展示**：每条风险均附体检原文与条款原文定位
- ❓ **追问问题清单**：生成用户需要补充确认的关键信息
- 🎯 **行动建议**：提供下一步操作指引
- 🔒 **隐私保护**：数据加密存储，支持一键删除
- 🧠 **RAG 条款定位**：按条款结构切片并支持向量召回重排
- 🌍 **多语言支持**：前端中英文切换，导出报告支持多语言模板
- 🖼️ **OCR 扩展能力**：提供 Python `document-parser` 服务接口

## 🚀 快速开始

### 前置要求

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 14+
- Redis 7+

### 一键启动（推荐）

```bash
# 克隆仓库
git clone https://github.com/zhenglizhi/policy-fit.git
cd policy-fit

# 启动基础设施
docker-compose up -d

# 配置环境变量
cp .env.example .env
# 编辑 .env 填入 LLM API Key

# 校验配置
make env-check

# 运行数据库迁移
make migrate-up

# 启动 API 服务
make run-api

# 启动 Worker 服务
make run-worker
```

访问:
- API: http://localhost:8080
- 健康探针: http://localhost:8080/health
- 就绪探针: http://localhost:8080/ready
- 指标快照: http://localhost:8080/metrics

### 手动安装

```bash
# 安装依赖
go mod download

# 编译
make build

# 运行
./bin/api
./bin/worker
```

### 前端本地运行（Next.js）

```bash
cd web
npm install
npm run dev
```

访问 http://localhost:3001

### 鉴权说明（JWT）

- `/api/v1/*` 默认启用 Bearer Token 鉴权
- Token payload 需包含 `user_id`
- 请求会自动注入 `X-Request-ID` 并在响应体返回 `request_id`

## 📖 文档

- [产品需求文档 (PRD)](./PRD-保单避坑雷达-v1.md)
- [API 文档](./docs/api.md)
- [部署指南](./docs/deployment.md)
- [迁移回滚指南](./docs/migration-rollback.md)
- [错误码文档](./docs/error-codes.md)
- [架构文档](./docs/architecture.md)
- [运行手册](./docs/runbooks/operations.md)
- [用户协议草案](./docs/compliance-user-agreement.md)
- [贡献指南](./CONTRIBUTING.md)

## 🏗️ 架构

```
┌─────────────┐
│   前端 UI   │
└──────┬──────┘
       │ HTTP
┌──────▼──────────────────────────────────┐
│          API Gateway (Gin)              │
│  ┌────────────────────────────────────┐ │
│  │  Handler → Service → Repository    │ │
│  └────────────────────────────────────┘ │
└──────┬──────────────────────────────────┘
       │
┌──────▼──────────┐    ┌─────────────────┐
│   PostgreSQL    │    │  Redis (Queue)  │
└─────────────────┘    └────────┬────────┘
                                │
                       ┌────────▼────────┐
                       │  Analysis Worker│
                       │  ┌────────────┐ │
                       │  │ Parser     │ │
                       │  │ LLM Client │ │
                       │  │ Rule Engine│ │
                       │  └────────────┘ │
                       └─────────────────┘
```

## 🛠️ 技术栈

- **后端框架**: Gin
- **数据库**: PostgreSQL
- **缓存/队列**: Redis
- **文档解析**: pdftotext
- **日志**: Zap
- **配置管理**: Viper

## ⚙️ 配置策略

- 配置文件按 `APP_ENV` 分层加载，优先级如下：
1. `.env.<APP_ENV>.local`
2. `.env.<APP_ENV>`
3. `.env`
- 默认 `APP_ENV=dev`，可设置为 `test` 或 `prod`。
- 使用 `make env-check` 在启动前做必填项校验。
- MVP 阶段不支持配置热加载，修改配置后需要重启服务生效。

## ⚠️ 免责声明

本工具仅用于**条款辅助解读与风险提示**，不构成保险销售建议或理赔结论。

- 最终承保与理赔结果以保险公司核保、合同原文与调查结论为准
- 不提供法律意见，不替代医生诊断意见
- 用户上传的健康资料属于敏感信息，请妥善保管

## 📄 开源协议

本项目采用 [Apache License 2.0](LICENSE) 开源协议。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

详见 [贡献指南](./CONTRIBUTING.md)

## 📧 联系方式

- Issue: https://github.com/zhenglizhi/policy-fit/issues
- Email: your-email@example.com
