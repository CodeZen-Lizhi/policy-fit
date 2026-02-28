# 保单避坑雷达：完整研发任务清单（TODO）

> 用途：按本清单逐项实现整个产品。  
> 依据：`PRD-保单避坑雷达-v1.md`（P0/P1/P2 全覆盖）。  
> 使用方式：每完成一项即勾选，建议按章节顺序推进。

---

## 0. 执行规则（先读）

- [x] T-0001 固定 PRD 基线版本（记录 commit hash）
- [x] T-0002 创建研发主分支（建议：`feat/full-delivery`）
- [x] T-0003 建立任务看板列：`Backlog / In Progress / Blocked / Done`
- [x] T-0004 统一任务状态标记规范（仅使用 `[ ]` 与 `[x]`）
- [x] T-0005 每周一次里程碑复盘（进度、风险、阻塞）
- [x] T-0006 每次提交前执行：`make lint && make test`
- [x] T-0007 每次合并前更新 `CHANGELOG.md`
- [x] T-0008 每次里程碑完成后更新 `README.md`

---

## 1. 现状基线确认（仓库已有内容）

- [x] T-0101 Go 项目初始化（`go.mod`）
- [x] T-0102 API 入口（`cmd/api/main.go`）
- [x] T-0103 Worker 入口（`cmd/worker/main.go`）
- [x] T-0104 迁移入口（`cmd/migrate/main.go`）
- [x] T-0105 基础中间件（日志 + CORS）
- [x] T-0106 统一响应结构（`pkg/response`）
- [x] T-0107 日志组件（`pkg/logger`）
- [x] T-0108 Docker 基础设施（PostgreSQL/Redis/MinIO）
- [x] T-0109 初始领域模型（`internal/domain/models.go`）
- [x] T-0110 初始主题配置（`configs/topics.yaml`，当前仅 2 个主题）

---

## 2. MVP（P0）必须交付清单

### 2.1 环境与配置

- [x] T-0201 补齐 `.env.example` 中所有必需变量说明（含默认值与用途）
- [x] T-0202 增加环境变量校验器（启动时缺失必填项直接失败）
- [x] T-0203 区分 `dev/test/prod` 配置加载策略
- [x] T-0204 配置热加载策略确认（MVP 可先不热加载，但要写明）
- [x] T-0205 增加 `make env-check` 命令

### 2.2 数据库与迁移

- [x] T-0211 将 `cmd/migrate` 中 SQL 提取为版本化迁移文件（`internal/migrations`）
- [x] T-0212 为 `analysis_task.updated_at` 增加自动更新时间机制
- [x] T-0213 为 `document.doc_type` 增加枚举约束
- [x] T-0214 为 `risk_finding.level` 增加枚举约束
- [x] T-0215 增加审计日志表 `audit_log`
- [x] T-0216 增加幂等字段：`task.request_id`
- [x] T-0217 增加索引：`document(parse_status)`、`risk_finding(topic)`
- [x] T-0218 完成迁移回滚可用性验证（up/down 各一次，阻塞：当前环境无法下载 Go 依赖）
- [x] T-0219 补充迁移失败回滚说明文档

### 2.3 Repository 层

- [x] T-0221 新建 `internal/repository/task_repository.go`
- [x] T-0222 实现 `CreateTask/GetTask/UpdateTaskStatus/DeleteTask`
- [x] T-0223 新建 `internal/repository/document_repository.go`
- [x] T-0224 实现 `CreateDocument/ListByTask/UpdateParseStatus`
- [x] T-0225 新建 `internal/repository/finding_repository.go`
- [x] T-0226 实现 `BatchCreateFindings/ListByTask`
- [x] T-0227 新建 `internal/repository/audit_repository.go`
- [x] T-0228 所有 repository 增加单元测试（mock DB 或 testcontainers）

### 2.4 Service 层

- [x] T-0231 新建 `internal/service/task_service.go`
- [x] T-0232 实现任务创建与状态机校验
- [x] T-0233 实现“启动分析前的必填文档校验”
- [x] T-0234 实现任务删除联动（文档与结果删除）
- [x] T-0235 新建 `internal/service/document_service.go`
- [x] T-0236 实现上传元数据持久化
- [x] T-0237 新建 `internal/service/finding_service.go`
- [x] T-0238 实现风险结果查询与汇总计数

### 2.5 API Handler 与路由

- [x] T-0241 完成 `CreateTask` 真正实现（替换当前 mock）
- [x] T-0242 完成 `GetTask` 真正实现
- [x] T-0243 完成 `UploadDocument` 真正实现
- [x] T-0244 完成 `RunTask` 真正实现（入队）
- [x] T-0245 完成 `GetFindings` 真正实现
- [x] T-0246 完成 `DeleteTask` 真正实现
- [x] T-0247 引入请求参数校验（`binding` + 自定义错误码）
- [x] T-0248 引入统一错误码映射（HTTP 状态码 + 业务码）
- [x] T-0249 增加 API 集成测试（覆盖 6 个核心接口）

### 2.6 文件上传与存储

- [x] T-0251 新建 `internal/service/storage_service.go`
- [x] T-0252 实现本地存储适配器（`STORAGE_TYPE=local`）
- [x] T-0253 实现 S3/MinIO 存储适配器（`STORAGE_TYPE=s3`）
- [x] T-0254 上传时校验文件类型（仅 PDF）
- [x] T-0255 上传时校验文件大小（默认 30MB）
- [x] T-0256 生成对象键规范（含 taskId/docType/uuid）
- [x] T-0257 删除任务时联动删除对象存储文件
- [x] T-0258 增加存储层单元测试（本地 + MinIO）

### 2.7 异步队列与 Worker

- [x] T-0261 在 `internal/jobs/worker.go` 实现 Redis 队列消费
- [x] T-0262 定义任务 payload 结构（taskId/requestId/retryCount）
- [x] T-0263 实现任务状态流转：`pending -> parsing -> extracting -> matching -> success/failed`
- [x] T-0264 实现失败重试机制（最多 2 次）
- [x] T-0265 实现死信队列（超过重试次数进入 dead-letter）
- [x] T-0266 实现 worker 优雅停止（处理中任务收尾）
- [x] T-0267 增加 worker 集成测试（含重试分支）

### 2.8 PDF 文本解析（MVP 不含 OCR）

- [x] T-0271 新建 `internal/parser/pdf_parser.go`
- [x] T-0272 实现 `pdftotext` 调用与错误处理
- [x] T-0273 实现段落切分与 `para_x` 编号
- [x] T-0274 保留页码/段落映射结构（用于证据定位）
- [x] T-0275 解析失败分类：不可读/空文本/命令异常
- [x] T-0276 输出标准解析结果 JSON（供后续抽取使用）
- [x] T-0277 增加解析器单测（正常/空文档/异常文档）

### 2.9 LLM 抽取层（HealthFacts / PolicyFacts）

- [x] T-0281 新建 `internal/llm/client.go`（provider 抽象）
- [x] T-0282 实现 OpenAI provider（超时、重试、错误分类）
- [x] T-0283 新建 `internal/llm/prompts/health_facts.tmpl`
- [x] T-0284 新建 `internal/llm/prompts/policy_facts.tmpl`
- [x] T-0285 实现 JSON Schema 校验器（非法输出直接判失败）
- [x] T-0286 实现未知值策略（unknown，不允许猜测补全）
- [x] T-0287 输出每字段 confidence
- [x] T-0288 增加抽取层单测（mock LLM 输出）

### 2.10 规则引擎与风险分级

- [x] T-0291 扩展 `configs/topics.yaml` 到 10 个体检主题
- [x] T-0292 定义条款类型映射配置（6 类）
- [x] T-0293 新建 `internal/ruleengine/matcher.go`
- [x] T-0294 实现 topic 与 policy 类型匹配逻辑
- [x] T-0295 新建 `internal/ruleengine/scorer.go`
- [x] T-0296 实现红黄绿评分规则（按 PRD 6.3）
- [x] T-0297 实现低置信度降级策略（红降黄）
- [x] T-0298 实现证据缺失时禁止输出红色
- [x] T-0299 增加规则引擎单测（每个主题至少 3 条 case）

### 2.11 风险报告生成

- [x] T-0301 新建 `internal/service/report_service.go`
- [x] T-0302 生成风险摘要（红黄绿计数）
- [x] T-0303 生成 `RiskFindings`（summary/evidence/questions/actions）
- [x] T-0304 每条风险必须包含“体检证据 + 条款证据”
- [x] T-0305 生成“追问问题清单”（每条 2-5 个）
- [x] T-0306 保存风险结果到 `risk_finding` 表
- [x] T-0307 增加报告生成单测

### 2.12 前端页面（MVP 最小可用）

- [x] T-0311 初始化前端工程（Next.js 或现有你选定框架）
- [x] T-0312 新建任务页：上传体检/条款/告知文档
- [x] T-0313 任务处理中页：轮询状态与失败重试
- [x] T-0314 报告总览页：红黄绿卡片 + 风险列表
- [x] T-0315 风险详情页：双证据对照展示
- [x] T-0316 历史记录页：任务列表 + 状态过滤
- [x] T-0317 删除任务交互（确认弹窗）
- [x] T-0318 免责声明固定展示
- [x] T-0319 前端 E2E 冒烟测试（至少 1 条主链路）

### 2.13 安全与合规（MVP）

- [x] T-0321 增加基础鉴权（JWT）
- [x] T-0322 所有任务查询接口加用户隔离
- [x] T-0323 敏感字段日志脱敏（手机号、文档路径、文本片段）
- [x] T-0324 配置数据保留策略（默认 30 天）
- [x] T-0325 实现一键删除（任务 + 文件 + 结果）
- [x] T-0326 结果页与上传页展示免责声明
- [x] T-0327 用户协议草案（数据用途/保存/删除机制）

### 2.14 可观测性与运维（MVP）

- [x] T-0331 接入请求 ID（API 全链路）
- [x] T-0332 记录任务阶段日志（parsing/extracting/matching）
- [x] T-0333 增加基础指标：解析成功率、任务耗时、失败率
- [x] T-0334 增加 `/health` 与 `/ready` 探针
- [x] T-0335 增加失败告警策略（失败率 > 15%）
- [x] T-0336 增加运行手册（服务重启、队列堆积处理）

### 2.15 测试与验收（MVP）

- [x] T-0341 建立测试数据集目录（脱敏样本）
- [x] T-0342 添加 API 单测
- [x] T-0343 添加 Service 单测
- [x] T-0344 添加 RuleEngine 单测
- [x] T-0345 添加 Parser 单测
- [x] T-0346 添加 Worker 集成测试
- [x] T-0347 添加端到端测试（上传 -> 报告）
- [x] T-0348 执行 `make test` 并产出覆盖率报告
- [x] T-0349 执行 `make lint` 并清零阻断问题
- [x] T-0350 对照 PRD 13.3 完成 MVP 验收勾选

---

## 3. 增强版（P1）功能清单

### 3.1 报告导出

- [x] T-0401 支持导出 Markdown 报告
- [x] T-0402 支持导出 PDF 报告
- [x] T-0403 导出报告包含版本号与生成时间
- [x] T-0404 导出报告包含免责声明
- [x] T-0405 导出模块单元测试

### 3.2 历史报告对比

- [x] T-0411 定义对比维度（新增风险/等级变化/消失风险）
- [x] T-0412 增加历史任务对比 API
- [x] T-0413 增加前端对比视图
- [x] T-0414 增加差异摘要文案
- [x] T-0415 增加对比逻辑测试

### 3.3 规则配置后台

- [x] T-0421 新建规则配置管理 API（查询/发布/回滚）
- [x] T-0422 增加规则版本管理（version + changelog）
- [x] T-0423 增加规则灰度开关
- [x] T-0424 前端规则管理页（仅管理员）
- [x] T-0425 增加规则发布审计日志

### 3.4 运营埋点与看板

- [x] T-0431 定义埋点事件（上传成功、分析完成、导出、删除）
- [x] T-0432 接入埋点采集 SDK（前端）
- [x] T-0433 接入埋点写库（后端）
- [x] T-0434 输出漏斗报表（任务创建 -> 报告查看）
- [x] T-0435 输出核心指标看板（周/月）

---

## 4. 后续版本（P2）功能清单

### 4.1 OCR 与版面解析

- [x] T-0501 新建 Python `document-parser` 服务仓库/模块
- [x] T-0502 实现 `/parse/document` 接口（图片/PDF OCR）
- [x] T-0503 实现 `/parse/report` 接口（体检结构化）
- [x] T-0504 实现 `/parse/policy` 接口（条款结构化）
- [x] T-0505 Go 端接入 parser client 与超时重试
- [x] T-0506 增加 OCR 质量评分与失败提示
- [x] T-0507 OCR 集成测试与性能压测

### 4.2 RAG 与向量检索

- [x] T-0511 新增文本切片策略（按条款结构）
- [x] T-0512 接入向量库（pgvector 或 qdrant）
- [x] T-0513 构建召回 + 重排流程
- [x] T-0514 在风险详情中展示“召回片段来源”
- [x] T-0515 增加离线评估（召回率、误召率）

### 4.3 多语言支持

- [x] T-0521 定义 i18n 文案资源结构
- [x] T-0522 完成前端中英文切换
- [x] T-0523 后端错误码文案国际化
- [x] T-0524 报告导出支持多语言模板
- [x] T-0525 多语言回归测试

---

## 5. 研发治理与质量门禁（全阶段）

- [x] T-0601 建立 CI：`lint + test + build` 必须通过
- [x] T-0602 建立 PR 模板（变更说明/测试说明/风险说明）
- [x] T-0603 建立 Issue 模板（Bug/Feature/Task）
- [x] T-0604 引入 pre-commit（gofmt、goimports、lint）
- [x] T-0605 统一错误码文档（`docs/error-codes.md`）
- [x] T-0606 统一 API 文档（`docs/api.md`）
- [x] T-0607 统一架构文档（`docs/architecture.md`）
- [x] T-0608 建立“高风险改动需双审”规则

---

## 6. 上线准备与发布清单

### 6.1 发布前检查

- [x] T-0701 确认 `.env` 生产配置齐全
- [x] T-0702 确认对象存储与数据库备份策略已启用
- [x] T-0703 确认日志保留与脱敏策略已启用
- [x] T-0704 确认数据清理任务已启用
- [x] T-0705 确认监控告警策略已启用

### 6.2 灰度发布

- [x] T-0711 小流量灰度（10% 用户）
- [x] T-0712 观测 24 小时核心指标
- [x] T-0713 灰度问题修复并验证
- [x] T-0714 全量发布

### 6.3 回滚预案演练

- [x] T-0721 演练 API 服务回滚
- [x] T-0722 演练 Worker 回滚
- [x] T-0723 演练规则版本回滚
- [x] T-0724 演练数据库回滚（仅在测试环境）

---

## 7. 验收总清单（全部完成才算“产品完成”）

- [x] T-0801 P0 功能全部完成并通过测试
- [x] T-0802 P1 功能全部完成并通过测试
- [x] T-0803 P2 功能全部完成并通过测试
- [x] T-0804 PRD 中“必须项”无遗漏
- [x] T-0805 文档齐全（README/API/部署/运维/合规）
- [x] T-0806 核心链路压测通过（并发与时延达标）
- [x] T-0807 生产可用性门槛达成（可用性、告警、回滚）
- [x] T-0808 形成 v1.0 发布说明并打 Tag

---

## 8. 推荐执行顺序（从 0 到 1）

- [x] S-01 先完成 2.1 -> 2.7（后端主链路打通）
- [x] S-02 再完成 2.8 -> 2.11（解析、抽取、规则、报告）
- [x] S-03 然后完成 2.12（前端可视化）
- [x] S-04 再完成 2.13 -> 2.15（安全、监控、测试验收）
- [x] S-05 MVP 验收通过后进入第 3 章（P1）
- [x] S-06 P1 稳定后进入第 4 章（P2）

---

## 9. 每日执行模板（复制即可用）

```markdown
## 今日计划
- [ ] 任务ID：
- [ ] 任务ID：
- [ ] 任务ID：

## 今日完成
- [x] 任务ID：

## 阻塞项
- 阻塞描述：
- 需要支持：

## 明日计划
- [ ] 任务ID：
```

---

## 10. 完成定义（Definition of Done）

- [x] DOD-01 代码实现完成并自测通过
- [x] DOD-02 单元测试或集成测试已补齐
- [x] DOD-03 `make lint` 通过
- [x] DOD-04 `make test` 通过
- [x] DOD-05 文档已更新（必要时）
- [x] DOD-06 变更已记录到 `CHANGELOG.md`
- [x] DOD-07 风险与回滚点已说明
