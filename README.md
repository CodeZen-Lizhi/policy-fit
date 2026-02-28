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

# 运行数据库迁移
make migrate-up

# 启动 API 服务
make run-api

# 启动 Worker 服务
make run-worker
```

访问 http://localhost:8080

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

## 📖 文档

- [产品需求文档 (PRD)](./PRD-保单避坑雷达-v1.md)
- [API 文档](./docs/api.md)
- [部署指南](./docs/deployment.md)
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
