# 贡献指南

感谢你对 Policy Fit 项目的关注！

## 贡献流程

### 1. Issue 先行

- 新功能/大改动：先开 Issue 讨论设计方案
- 小 Bug 修复：可直接提交 PR

### 2. 分支规范

从 `main` 分支创建特性分支：

- `feat/xxx` - 新功能
- `fix/xxx` - Bug 修复
- `docs/xxx` - 文档更新
- `refactor/xxx` - 代码重构
- `test/xxx` - 测试相关

### 3. Commit 规范

遵循 [Conventional Commits](https://www.conventionalcommits.org/)：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type 类型：**

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具链相关

**示例：**

```
feat(parser): 添加 OCR 扫描件支持

- 集成 Tesseract OCR 引擎
- 添加图片预处理流程
- 更新文档解析接口

Closes #123
```

### 4. 代码规范

**Go 代码：**

```bash
# 格式化
gofmt -w .

# Lint 检查
golangci-lint run ./...
```

**必须通过：**

- `go fmt`
- `golangci-lint`
- 单元测试覆盖率不低于当前基线

### 5. 测试要求

- 新功能必须附带单元测试
- 修复 Bug 需添加回归测试
- 运行 `make test` 确保所有测试通过

### 6. Pull Request

**PR 标题格式：**

```
<type>: <简短描述>
```

**PR 描述模板：**

```markdown
## 变更内容

简要描述本次变更的内容

## 关联 Issue

Closes #123

## 测试方式

描述如何测试本次变更

## 截图（如有 UI 变更）

粘贴截图

## Checklist

- [ ] 代码已通过 `make lint`
- [ ] 测试已通过 `make test`
- [ ] 文档已更新（如需要）
- [ ] CHANGELOG 已更新（如需要）
```

### 7. Code Review

- 至少需要 1 位 Maintainer 审核通过
- 解决所有 Review 意见后方可合并
- 保持 PR 小而聚焦，便于审核

## 开发环境搭建

```bash
# 克隆仓库
git clone https://github.com/zhenglizhi/policy-fit.git
cd policy-fit

# 安装依赖
make deps

# 启动基础设施
make docker-up

# 配置环境变量
cp .env.example .env

# 运行迁移
make migrate-up

# 启动服务
make run-api
```

## 项目结构

```
policy-fit/
├── cmd/              # 入口程序
├── internal/         # 内部代码（不对外暴露）
│   ├── handler/      # HTTP 处理器
│   ├── service/      # 业务逻辑
│   ├── domain/       # 领域模型
│   ├── repository/   # 数据访问
│   ├── parser/       # 文档解析
│   ├── llm/          # LLM 客户端
│   ├── ruleengine/   # 规则引擎
│   └── config/       # 配置管理
├── pkg/              # 可对外暴露的库
├── docs/             # 文档
└── deploy/           # 部署配置
```

## 常见问题

### Q: 如何添加新的体检异常主题？

1. 在 `configs/topics.yaml` 添加主题配置
2. 在 `internal/ruleengine/matcher.go` 添加匹配逻辑
3. 添加单元测试
4. 更新文档

### Q: 如何切换 LLM 提供商？

修改 `.env` 中的 `LLM_PROVIDER` 和相关配置即可。

### Q: 如何本地调试 Worker？

```bash
# 启动 Redis
make docker-up

# 运行 Worker（带调试日志）
LOG_LEVEL=debug make run-worker
```

## 行为准则

- 尊重所有贡献者
- 保持友好、专业的沟通
- 接受建设性的批评
- 关注项目整体利益

## 联系方式

- GitHub Issues: https://github.com/zhenglizhi/policy-fit/issues
- Email: your-email@example.com

再次感谢你的贡献！🎉
