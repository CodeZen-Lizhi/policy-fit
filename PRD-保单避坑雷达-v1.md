# PRD：保单避坑雷达（Policy × Health Fit）

## 1. 文档信息

- 版本：v1.0（可开发基线版）
- 状态：可进入研发排期
- 更新时间：2026-02-28
- 适用阶段：MVP（0 -> 1）
- 目标读者：产品、后端、前端、算法、测试、运维、合规

---

## 2. 产品定义

### 2.1 一句话定义

用户上传体检/检查报告与保险条款，系统输出“潜在拒赔/除外风险”与“条款-证据对照解释”，帮助用户在投保前后识别风险盲区。

### 2.2 目标用户

1. 准备投保的个人用户：担心体检异常导致核保限制或后续理赔争议。
2. 刚投保用户：希望复查保单条款与既往体检是否存在潜在冲突。
3. 家庭决策者：需要对配偶/父母保单做风险复核。

### 2.3 核心价值

1. 风险预警：红黄绿分级提示“可能不赔点”。
2. 可解释性：每条风险均附体检证据与条款原文定位。
3. 行动建议：给出补充确认问题与下一步操作建议。

### 2.4 产品边界（必须展示）

1. 本产品仅用于条款辅助解读与风险提示，不构成承保/理赔结论。
2. 最终承保与理赔结果以保险公司核保、合同原文与调查结论为准。
3. 不提供法律意见，不替代医生诊断意见。

---

## 3. 目标与非目标

### 3.1 业务目标（MVP，8周）

1. 完成上传-解析-风险报告全链路，首版可稳定处理“文字型 PDF”。
2. 覆盖 10 类常见体检异常 + 6 类核心条款要素。
3. 风险报告可追溯：每条结论均可点击查看原文证据。
4. 单份文档处理时长：P50 < 60s，P95 < 120s。

### 3.2 成功指标

1. 解析成功率（可生成结构化结果）：>= 90%（文字型 PDF）。
2. 风险结果可解释率（有证据引用）：>= 95%。
3. 用户完成一次报告生成的转化率：>= 70%。
4. 首次报告后 7 日内复访率：>= 20%。

### 3.3 非目标（MVP 不做）

1. 不支持扫描件/照片 OCR（仅预留接口）。
2. 不支持企业团险、复杂法务争议案件分析。
3. 不做自动承保建议报价与保险产品推荐。
4. 不做实时在线问诊与医生咨询。

---

## 4. 用户流程与核心场景

### 4.1 主流程

1. 用户创建分析任务。
2. 上传体检报告（必选）与保单条款（必选），可选上传投保告知。
3. 系统解析文档并抽取结构化事实。
4. 系统执行匹配与评分，生成红黄绿风险清单。
5. 用户查看风险详情、证据对照、追问清单。
6. 用户下载报告或继续补充资料后二次分析。

### 4.2 关键场景

1. 投保前体检异常复核：避免遗漏告知项。
2. 投保后保单风险复盘：定位免责/等待期影响。
3. 家庭保单年检：批量查看多个家庭成员风险。

### 4.3 异常流程

1. 文档无可提取文本：提示“当前文件可能为扫描件”，建议换文字版或等待 OCR 版本。
2. 条款抽取不完整：标记低置信度并展示“需要用户补充页码/附件”。
3. 日期缺失：风险等级上限降为黄，并提示补充检查日期。

---

## 5. 功能需求（按优先级）

### 5.1 P0（MVP 必须）

1. 任务管理
   - 新建任务、查看任务状态（待处理/处理中/已完成/失败）。
2. 文档上传
   - 支持 PDF 文件上传，单文件大小限制 30MB。
   - 支持文本直贴输入（兜底）。
3. 文档解析
   - 体检文本抽取为 `HealthFacts`。
   - 条款文本抽取为 `PolicyFacts`。
4. 风险识别
   - 输出红黄绿分级风险列表。
   - 每条风险包含摘要、证据、追问问题。
5. 报告展示
   - 列表页：风险总览与等级计数。
   - 详情页：体检证据段落 + 条款证据段落 + 解释。
6. 数据安全基础
   - 文件与结果加密存储。
   - 支持用户删除任务与关联文件。
7. 合规提示
   - 报告页固定展示免责声明。

### 5.2 P1（MVP+）

1. 报告导出（PDF/Markdown）。
2. 历史报告对比（同一用户多次分析差异）。
3. 规则配置后台（风险阈值、提示文案可配置）。
4. 运营埋点与漏斗分析看板。

### 5.3 P2（后续）

1. OCR 解析（扫描件、拍照件）。
2. RAG + 向量检索提升长条款定位精度。
3. 多语言支持（简中 -> 繁中/英文）。

---

## 6. 规则与主题覆盖

### 6.1 体检异常主题（首版 10 个）

1. 高血压/血压偏高
2. 血糖异常/糖尿病前期/糖尿病
3. 血脂异常
4. BMI 异常/肥胖
5. 脂肪肝/肝功能异常（ALT/AST）
6. 甲状腺结节（TI-RADS）
7. 肺结节
8. 心电图异常/心律失常提示
9. 尿酸高/痛风
10. 肾功能异常/蛋白尿

### 6.2 条款主题（首版 6 类）

1. 既往症定义
2. 责任免除
3. 等待期
4. 投保告知
5. 特定疾病定义
6. 续保与不可抗辩条款

### 6.3 红黄绿评分规则（首版）

1. 红色
   - 存在明确诊断/长期用药/复查就诊建议
   - 且条款存在强关联既往症/免责/告知命中
2. 黄色
   - 指标异常但诊断不明确，或时间信息缺失
   - 或条款描述模糊、证据不充分
3. 绿色
   - 当前证据未见明显冲突
   - 仅可表述为“暂未发现高风险冲突”，禁止“保证理赔”

---

## 7. 详细交互与页面需求

### 7.1 页面清单

1. 登录/注册页（可简化为验证码登录）
2. 新建任务页（上传区 + 文档要求提示）
3. 任务处理中页（进度与失败重试）
4. 报告总览页（红黄绿卡片 + 风险列表）
5. 风险详情页（证据对照 + 追问问题 + 建议行动）
6. 历史记录页（任务列表 + 状态 + 创建时间）
7. 隐私与设置页（数据删除、导出申请）

### 7.2 报告页信息架构（必须）

1. 顶部：风险总评（红/黄/绿数量，非结论声明）
2. 中部：逐条风险卡片
3. 底部：免责声明 + 数据来源说明 + 生成时间

### 7.3 风险卡片字段

1. `topic`：风险主题
2. `level`：红/黄/绿
3. `summary`：一句话风险说明
4. `healthEvidence`：体检原文 + 定位
5. `policyEvidence`：条款原文 + 定位
6. `questions`：追问清单（2-5 条）
7. `actions`：建议动作（补材料/咨询客服/补充告知）

---

## 8. 数据模型设计

### 8.1 核心实体

1. `user`
2. `analysis_task`
3. `document`
4. `health_fact`
5. `policy_fact`
6. `risk_finding`
7. `audit_log`

### 8.2 建议表结构（PostgreSQL）

```sql
-- 用户
CREATE TABLE user_account (
  id BIGSERIAL PRIMARY KEY,
  phone VARCHAR(32) UNIQUE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 任务
CREATE TABLE analysis_task (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES user_account(id),
  status VARCHAR(32) NOT NULL, -- pending/running/success/failed
  risk_summary JSONB,          -- {"red":2,"yellow":3,"green":5}
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 文档
CREATE TABLE document (
  id BIGSERIAL PRIMARY KEY,
  task_id BIGINT NOT NULL REFERENCES analysis_task(id),
  doc_type VARCHAR(32) NOT NULL, -- report/policy/disclosure
  file_name VARCHAR(256) NOT NULL,
  storage_key VARCHAR(512) NOT NULL,
  parse_status VARCHAR(32) NOT NULL, -- pending/success/failed
  parsed_text TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 风险发现
CREATE TABLE risk_finding (
  id BIGSERIAL PRIMARY KEY,
  task_id BIGINT NOT NULL REFERENCES analysis_task(id),
  level VARCHAR(16) NOT NULL, -- red/yellow/green
  topic VARCHAR(64) NOT NULL,
  summary TEXT NOT NULL,
  health_evidence JSONB NOT NULL,
  policy_evidence JSONB NOT NULL,
  questions JSONB NOT NULL,
  confidence NUMERIC(4,3),
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 8.3 结构化 JSON 定义

#### HealthFacts

```json
{
  "facts": [
    {
      "category": "hypertension",
      "label": "血压偏高/高血压风险",
      "evidence": {
        "text": "血压 155/95 mmHg，建议复查",
        "date": "2026-01-10",
        "loc": "para_12",
        "source": "report"
      },
      "values": { "sbp": 155, "dbp": 95 },
      "confidence": 0.92
    }
  ]
}
```

#### PolicyFacts

```json
{
  "sections": [
    {
      "type": "preexisting_definition",
      "title": "既往症定义",
      "content": "......",
      "loc": "para_120",
      "confidence": 0.89
    }
  ]
}
```

#### RiskFindings

```json
{
  "findings": [
    {
      "level": "red",
      "topic": "hypertension",
      "summary": "可能触发既往症/告知冲突",
      "health_evidence": [{ "loc": "para_12", "text": "血压 155/95..." }],
      "policy_evidence": [{ "loc": "para_120", "text": "既往症包括..." }],
      "questions": ["是否已确诊并长期用药？", "体检是否发生在投保前？"],
      "actions": ["补充投保告知截图", "联系保险公司核保确认"]
    }
  ]
}
```

---

## 9. 技术架构设计（可直接开工）

### 9.1 总体架构

1. 前端（Next.js）
2. Go API 服务（Gin/Fiber）
3. 任务队列（Redis + Worker）
4. 文本解析层（`pdftotext`，后续可切换 Python OCR）
5. LLM 抽取层（HealthFacts/PolicyFacts）
6. 规则引擎层（映射 + 打分 + 文案模板）
7. 数据层（PostgreSQL + MinIO/S3）

### 9.2 服务边界

1. `api-gateway`：鉴权、任务入口、报告查询。
2. `analysis-worker`：解析、抽取、评分、持久化。
3. `document-parser`（预留）：OCR 与结构化增强服务。

### 9.3 目录结构建议（Go）

```text
policy-fit/
  cmd/
    api/
    worker/
  internal/
    handler/
    service/
    domain/
    repository/
    parser/
    llm/
    ruleengine/
    jobs/
    config/
  web/
    app/
    components/
  deploy/
    docker-compose.yml
  docs/
```

### 9.4 API 设计（MVP）

1. `POST /api/v1/tasks`
   - 入参：任务名、可选备注
   - 出参：`taskId`
2. `POST /api/v1/tasks/{taskId}/documents`
   - 入参：`docType` + file
   - 出参：`documentId`
3. `POST /api/v1/tasks/{taskId}/run`
   - 入参：空
   - 出参：任务状态
4. `GET /api/v1/tasks/{taskId}`
   - 出参：任务详情 + 风险摘要
5. `GET /api/v1/tasks/{taskId}/findings`
   - 出参：风险列表
6. `DELETE /api/v1/tasks/{taskId}`
   - 行为：删除任务、文档、解析结果

### 9.5 异步任务状态机

1. `pending` -> `parsing`
2. `parsing` -> `extracting`
3. `extracting` -> `matching`
4. `matching` -> `success`
5. 任一阶段失败 -> `failed`（记录失败码与重试次数）

---

## 10. 算法与提示词策略

### 10.1 执行链路

1. 文本切分：按段落编号（`para_x`）并保留页码映射。
2. LLM 抽取：生成 HealthFacts / PolicyFacts。
3. 规则匹配：按 `topic -> clause_type` 映射计算风险。
4. 解释生成：基于命中证据生成人话说明与追问清单。
5. 置信度校验：低于阈值的结论降级并标黄。

### 10.2 规则配置样例（YAML）

```yaml
topics:
  hypertension:
    hit_keywords: ["高血压", "血压偏高", "收缩压", "舒张压"]
    policy_types: ["preexisting_definition", "exclusion", "underwriting_disclosure"]
    red_conditions:
      - "diagnosed == true"
      - "long_term_medication == true"
    yellow_conditions:
      - "date_missing == true"
      - "diagnosed == false and abnormal_index == true"
```

### 10.3 提示词要求（工程约束）

1. 抽取提示词输出必须是 JSON，不允许自然语言混杂。
2. 每个字段必须返回 `confidence`，便于后续降级策略。
3. 对无法确认信息返回 `unknown`，禁止臆测补全。

---

## 11. 非功能需求

### 11.1 性能

1. 单任务文档总页数 <= 150 页。
2. 并发目标：首版支持 20 并发任务。
3. API P95 响应时间 < 500ms（非分析接口）。

### 11.2 可用性

1. 核心 API 可用性 >= 99.5%（月度）。
2. 任务失败可重试（最多 2 次）。

### 11.3 可观测性

1. 结构化日志：请求 ID、任务 ID、阶段、错误码。
2. 指标：解析成功率、平均耗时、模型调用失败率。
3. 告警：任务失败率 > 15% 触发告警。

### 11.4 安全与隐私

1. 数据传输 HTTPS，全链路 TLS。
2. 对象存储加密（服务端加密或应用层加密）。
3. 默认 30 天自动清理，支持用户即时删除。
4. 最小权限访问，后台操作写入审计日志。

---

## 12. 合规与法务需求

1. 页面固定展示风险提示与免责声明。
2. 禁止输出“必赔/拒赔确定性结论”。
3. 用户授权协议需明确：
   - 上传数据用途（仅分析与改进）
   - 保存时长
   - 删除机制
4. 保留报告生成版本号，便于争议时回溯。

---

## 13. 测试与验收标准

### 13.1 测试数据集

1. 体检报告样本 >= 100 份（脱敏）。
2. 保单条款样本 >= 60 份（覆盖 3 类险种）。
3. 人工标注集 >= 200 条“异常-条款命中”样本。

### 13.2 测试维度

1. 功能测试：流程可达、接口正确、状态转换正确。
2. 准确性测试：抽取准确率、风险命中率、误报率。
3. 异常测试：空文档、乱码、极端长文档、缺页。
4. 安全测试：越权访问、删除一致性、敏感日志泄漏。

### 13.3 MVP 验收门槛

1. P0 功能全部完成且通过测试。
2. 解析成功率 >= 90%（文字版 PDF）。
3. 每条风险均有双证据引用（体检+条款）。
4. 关键流程无阻断级缺陷（P0/P1 bug=0）。

---

## 14. 项目计划（8 周）

### 第 1-2 周：基础设施与主链路

1. 完成项目脚手架、数据库、对象存储、任务状态机。
2. 打通上传 -> 解析 -> 持久化最小链路。

### 第 3-4 周：结构化抽取

1. 落地 HealthFacts / PolicyFacts 抽取。
2. 完成 10 异常主题与 6 条款类型映射配置。

### 第 5-6 周：风险引擎与前端报告

1. 实现红黄绿评分规则引擎。
2. 完成报告总览页 + 详情页 + 追问清单展示。

### 第 7 周：稳定性与安全

1. 审计日志、错误重试、监控告警。
2. 数据删除、生命周期清理能力。

### 第 8 周：联调与灰度

1. 完成验收测试与回归测试。
2. 灰度上线并跟踪核心指标。

---

## 15. 风险清单与应对

1. 文档质量差导致抽取失败
   - 应对：增加失败分类与用户引导；预留 OCR 服务。
2. LLM 输出不稳定
   - 应对：强制 JSON Schema 校验 + 重试 + 降级。
3. 规则误报影响用户信任
   - 应对：每条结论强制证据展示，低置信度降黄。
4. 合规风险
   - 应对：前后台统一免责声明模板，审计可追溯。
5. 成本失控
   - 应对：分层模型策略（抽取模型与解释模型分离），缓存中间结果。

---

## 16. 回滚与应急预案

1. 功能回滚：开关控制“风险评分模块”与“解释生成模块”。
2. 模型回滚：保留上一个稳定提示词版本与模型版本。
3. 数据回滚：任务级别软删除 + 日志可追溯恢复。
4. 应急策略：当抽取服务异常时，仅输出“解析失败与补传建议”，不生成风险结论。

---

## 17. 开发启动清单（可直接执行）

1. 初始化仓库目录与 `docker-compose`（PostgreSQL/Redis/MinIO）。
2. 建立数据表与迁移脚本。
3. 完成 6 个核心 API 与任务队列。
4. 接入 `pdftotext` 并完成段落定位输出。
5. 完成抽取 JSON Schema 与校验器。
6. 落地评分规则配置文件与执行器。
7. 开发前端 3 个关键页面（上传/总览/详情）。
8. 增加日志、监控、错误码体系。
9. 完成测试样本准备与验收用例。

---

## 18. 附录

### 18.1 术语

1. 既往症：在投保前已存在或已被诊断的疾病/异常状态。
2. 告知义务：投保时对保险公司询问事项真实说明的义务。
3. 免责条款：保险公司在特定条件下不承担赔付责任的约定。

### 18.2 版本演进建议

1. v1.1：增加 OCR 支持与扫描件质量检测。
2. v1.2：增加险种模板（医疗险/重疾险/寿险）差异化规则。
3. v1.3：加入用户反馈闭环，驱动规则迭代。

---

## 19. 异常主题词典（10 个主题完整配置）

每个主题定义：识别关键词、指标阈值、同义词归一、与条款类型映射关系。

### 19.1 高血压 / 血压偏高

```yaml
topic: hypertension
keywords:
  - 高血压
  - 血压偏高
  - 收缩压升高
  - 舒张压升高
  - 血压异常
  - 建议降压
synonyms:
  - 血压偏高 → hypertension
  - 高血压病 → hypertension
thresholds:
  sbp_abnormal: 140      # 收缩压 >= 140 mmHg 视为异常
  dbp_abnormal: 90       # 舒张压 >= 90 mmHg 视为异常
  sbp_borderline: 130    # 130-139 为黄色边界值
red_signals:
  - 明确诊断高血压
  - 长期服用降压药
  - 合并心脑血管并发症
yellow_signals:
  - 单次血压偏高，建议复查
  - 无确诊记录，仅指标超标
policy_types:
  - preexisting_definition
  - exclusion
  - underwriting_disclosure
questions:
  - 是否已被医生明确诊断为高血压？
  - 是否正在长期服用降压药物？
  - 体检日期是否在本次投保前？
  - 投保告知问卷是否询问了血压情况？
  - 是否有心脑血管相关并发症检查记录？
```

### 19.2 血糖异常 / 糖尿病前期 / 糖尿病

```yaml
topic: diabetes
keywords:
  - 血糖
  - 空腹血糖
  - 餐后血糖
  - 糖化血红蛋白
  - HbA1c
  - 糖尿病
  - 糖耐量异常
  - 胰岛素抵抗
thresholds:
  fasting_abnormal: 7.0      # 空腹血糖 >= 7.0 mmol/L 确诊
  fasting_borderline: 6.1    # 6.1-6.9 为糖尿病前期
  hba1c_abnormal: 6.5        # HbA1c >= 6.5% 确诊
red_signals:
  - 明确诊断 2 型糖尿病
  - 长期服用降糖药或注射胰岛素
  - HbA1c >= 7.0%
yellow_signals:
  - 空腹血糖 6.1-6.9 mmol/L
  - 糖耐量异常（IGT）
  - 仅建议复查，无确诊
policy_types:
  - preexisting_definition
  - exclusion
  - specific_disease_definition
  - underwriting_disclosure
questions:
  - 是否已确诊糖尿病（1型/2型）？
  - 是否正在服用降糖药或注射胰岛素？
  - 最近一次 HbA1c 数值是多少？
  - 是否有糖尿病并发症（视网膜/肾脏/神经）相关检查？
```

### 19.3 血脂异常

```yaml
topic: dyslipidemia
keywords:
  - 血脂
  - 总胆固醇
  - 甘油三酯
  - 低密度脂蛋白
  - 高密度脂蛋白
  - LDL
  - HDL
  - TC
  - TG
  - 血脂偏高
thresholds:
  tc_abnormal: 6.2          # 总胆固醇 >= 6.2 mmol/L
  tg_abnormal: 2.3          # 甘油三酯 >= 2.3 mmol/L
  ldl_abnormal: 4.1         # LDL >= 4.1 mmol/L
red_signals:
  - 明确诊断高脂血症
  - 长期服用他汀类药物
  - 合并动脉粥样硬化或心血管疾病
yellow_signals:
  - 指标边界偏高，建议复查
  - 无确诊，仅单项异常
policy_types:
  - preexisting_definition
  - exclusion
  - underwriting_disclosure
questions:
  - 是否已确诊高脂血症？
  - 是否正在服用调脂药物？
  - 是否合并其他心血管疾病风险因素？
```

### 19.4 BMI 异常 / 肥胖

```yaml
topic: obesity
keywords:
  - BMI
  - 体重指数
  - 肥胖
  - 超重
  - 体重超标
thresholds:
  bmi_overweight: 24.0      # BMI >= 24 为超重（中国标准）
  bmi_obese: 28.0           # BMI >= 28 为肥胖
red_signals:
  - BMI >= 32，合并代谢综合征
yellow_signals:
  - BMI 24-27.9，仅体重超标
policy_types:
  - underwriting_disclosure
  - exclusion
questions:
  - 体重是否在近期发生明显变化？
  - 是否合并高血压、高血糖等代谢相关疾病？
```

### 19.5 脂肪肝 / 肝功能异常

```yaml
topic: fatty_liver
keywords:
  - 脂肪肝
  - ALT
  - AST
  - 谷丙转氨酶
  - 谷草转氨酶
  - 肝功能异常
  - 肝酶升高
  - 轻度脂肪肝
  - 中度脂肪肝
  - 重度脂肪肝
thresholds:
  alt_abnormal: 40          # ALT > 40 U/L
  ast_abnormal: 40          # AST > 40 U/L
red_signals:
  - 中/重度脂肪肝
  - ALT/AST 超过正常值 3 倍以上
  - 合并肝硬化或肝炎病史
yellow_signals:
  - 轻度脂肪肝
  - ALT/AST 轻度升高（1-3 倍）
policy_types:
  - preexisting_definition
  - exclusion
  - underwriting_disclosure
questions:
  - 是否有乙肝/丙肝病史或携带者状态？
  - 是否饮酒及饮酒频率/量？
  - 是否有肝脏相关用药史？
```

### 19.6 甲状腺结节

```yaml
topic: thyroid_nodule
keywords:
  - 甲状腺结节
  - TI-RADS
  - 甲状腺占位
  - 甲状腺低回声
  - 甲状腺钙化
  - 甲状腺肿物
thresholds:
  ti_rads_watch: 3           # TI-RADS 3 类建议随访
  ti_rads_biopsy: 4          # TI-RADS 4 类建议穿刺
red_signals:
  - TI-RADS 4 类及以上
  - 已行穿刺活检或手术
  - 确诊甲状腺癌或甲状腺功能异常
yellow_signals:
  - TI-RADS 3 类，建议随访
  - 单纯结节，无功能异常
policy_types:
  - preexisting_definition
  - exclusion
  - specific_disease_definition
questions:
  - 结节大小与 TI-RADS 分级是多少？
  - 是否已做穿刺活检？结果如何？
  - 是否合并甲亢或甲减，是否服药？
```

### 19.7 肺结节

```yaml
topic: pulmonary_nodule
keywords:
  - 肺结节
  - 肺部结节
  - 肺部阴影
  - 磨玻璃结节
  - GGO
  - 肺占位
  - 建议随访
thresholds:
  size_watch_mm: 6           # 直径 >= 6mm 建议随访
  size_biopsy_mm: 15         # 直径 >= 15mm 建议活检
red_signals:
  - 已行穿刺/手术，或确诊肺癌
  - 结节快速增长，高度可疑恶性
yellow_signals:
  - 磨玻璃结节 6-14mm，建议随访
  - 实性结节，良性特征明显
policy_types:
  - preexisting_definition
  - exclusion
  - specific_disease_definition
questions:
  - 结节大小、性质（磨玻璃/实性/混合）是什么？
  - 是否已随访，随访结果有无变化？
  - 是否有吸烟史？
  - 投保告知是否询问过肺部检查结果？
```

### 19.8 心电图异常 / 心律失常

```yaml
topic: ecg_abnormal
keywords:
  - 心电图异常
  - 心律不齐
  - 窦性心动过速
  - 窦性心动过缓
  - 房颤
  - 室性早搏
  - ST 段改变
  - T 波异常
  - 束支传导阻滞
  - 心肌缺血提示
red_signals:
  - 确诊心房颤动/室颤
  - 明确心肌缺血或冠心病
  - 已行心脏手术或植入起搏器
yellow_signals:
  - 偶发室早，无症状
  - ST/T 轻度改变，建议复查
  - 窦性心律不齐（青少年常见）
policy_types:
  - preexisting_definition
  - exclusion
  - specific_disease_definition
  - underwriting_disclosure
questions:
  - 是否有心悸、胸闷、晕厥等症状？
  - 是否做过 24 小时动态心电图或超声心动图？
  - 是否服用抗心律失常药物？
```

### 19.9 尿酸高 / 痛风

```yaml
topic: hyperuricemia
keywords:
  - 尿酸
  - 尿酸升高
  - 高尿酸血症
  - 痛风
  - 痛风性关节炎
thresholds:
  uric_acid_male: 420        # 男性 > 420 μmol/L 为高尿酸
  uric_acid_female: 360      # 女性 > 360 μmol/L 为高尿酸
red_signals:
  - 明确确诊痛风性关节炎
  - 长期服用降尿酸药物
  - 痛风石或肾结石病史
yellow_signals:
  - 尿酸偏高，无症状性高尿酸血症
policy_types:
  - preexisting_definition
  - exclusion
  - underwriting_disclosure
questions:
  - 是否有痛风发作史（关节红肿热痛）？
  - 是否长期服用非布司他/别嘌醇等药物？
  - 是否合并肾功能异常或肾结石？
```

### 19.10 肾功能异常 / 蛋白尿

```yaml
topic: renal_abnormal
keywords:
  - 肾功能异常
  - 肌酐升高
  - 尿素氮升高
  - 蛋白尿
  - 尿蛋白阳性
  - eGFR 降低
  - 肾小球滤过率
thresholds:
  creatinine_male: 115       # 男性肌酐 > 115 μmol/L 异常
  creatinine_female: 97      # 女性肌酐 > 97 μmol/L 异常
  egfr_watch: 90             # eGFR < 90 需关注
red_signals:
  - 明确确诊慢性肾病（CKD 分期）
  - 肌酐持续升高或 eGFR < 60
  - 蛋白尿 2+ 及以上
yellow_signals:
  - 单次肌酐轻度升高，建议复查
  - 尿蛋白微量阳性
policy_types:
  - preexisting_definition
  - exclusion
  - underwriting_disclosure
questions:
  - 是否确诊慢性肾病？分期是什么？
  - 是否有糖尿病肾病或高血压肾病病史？
  - 是否在服用保护肾功能的药物？
```

---

## 20. LLM 提示词模板

### 20.1 HealthFacts 抽取提示词

```text
你是一名医疗文档结构化分析助手。请从以下体检报告文本中抽取健康异常事实，输出严格的 JSON 格式，不得包含任何自然语言说明。

【抽取规则】
1. 只抽取与以下类别相关的异常项：hypertension, diabetes, dyslipidemia, obesity, fatty_liver, thyroid_nodule, pulmonary_nodule, ecg_abnormal, hyperuricemia, renal_abnormal。
2. 每条事实必须包含原文片段（evidence.text）、段落定位（evidence.loc）、检查日期（evidence.date，如无则填 "unknown"）。
3. 每条事实必须包含 confidence（0.0-1.0），低于 0.6 时必须填写 uncertain_reason。
4. 对无法确认的字段，填写 "unknown"，禁止推断补全。
5. 如报告中无任何相关异常，返回 {"facts": []}。

【输出格式】
{
  "facts": [
    {
      "category": "<类别英文标识>",
      "label": "<中文标签>",
      "evidence": {
        "text": "<原文片段>",
        "date": "<检查日期或 unknown>",
        "loc": "<段落索引，如 para_12>",
        "source": "report"
      },
      "values": {},
      "diagnosed": true | false | "unknown",
      "long_term_medication": true | false | "unknown",
      "confidence": 0.0-1.0,
      "uncertain_reason": "<置信度低于 0.6 时填写原因>"
    }
  ]
}

【体检报告文本】
{{REPORT_TEXT}}
```

### 20.2 PolicyFacts 抽取提示词

```text
你是一名保险条款结构化分析助手。请从以下保险合同文本中抽取关键条款内容，输出严格的 JSON 格式，不得包含任何自然语言说明。

【抽取规则】
1. 只抽取以下类型的条款：preexisting_definition（既往症定义）、exclusion（责任免除）、waiting_period（等待期）、underwriting_disclosure（投保告知）、specific_disease_definition（特定疾病定义）、renewal（续保条款）。
2. 每条条款必须包含原文内容（content）与段落定位（loc）。
3. 每条条款必须包含 confidence（0.0-1.0）。
4. 投保告知类型需额外提取问题列表（questions 字段）。
5. 如无法定位某类条款，不输出该类型，禁止补全。

【输出格式】
{
  "sections": [
    {
      "type": "<条款类型>",
      "title": "<条款标题>",
      "content": "<条款原文>",
      "loc": "<段落索引，如 para_120>",
      "confidence": 0.0-1.0,
      "questions": ["<投保告知问题1>", "<投保告知问题2>"]
    }
  ]
}

【保险合同文本】
{{POLICY_TEXT}}
```

### 20.3 风险解释生成提示词

```text
你是一名保险条款解读助手。根据以下体检异常证据与对应条款内容，用简洁的中文解释潜在拒赔/除外风险，并生成追问问题清单。

【约束】
1. 只能基于提供的证据作出解释，禁止引入未提供的信息。
2. 不得给出"一定赔/一定不赔"的确定性结论。
3. 每条解释不超过 80 字。
4. 追问问题 2-5 条，聚焦"用户需要补充确认的关键信息"。
5. 行动建议 1-3 条，聚焦"用户可以采取的下一步操作"。
6. 输出严格 JSON，禁止自然语言混杂。

【体检证据】
{{HEALTH_EVIDENCE}}

【条款证据】
{{POLICY_EVIDENCE}}

【输出格式】
{
  "explanation": "<风险解释>",
  "questions": ["<追问问题1>", "<追问问题2>"],
  "actions": ["<建议行动1>", "<建议行动2>"]
}
```

---

## 21. 错误码体系

### 21.1 错误码规范

格式：`PFIT-{模块码}-{错误序号}`

| 模块码 | 模块名称 |
|--------|----------|
| 1000   | 通用/系统 |
| 2000   | 文档上传与解析 |
| 3000   | LLM 抽取 |
| 4000   | 风险评分引擎 |
| 5000   | 任务管理 |

### 21.2 错误码清单

| 错误码 | HTTP 状态 | 说明 | 用户提示 |
|--------|-----------|------|----------|
| PFIT-1001 | 400 | 请求参数缺失或格式错误 | 请检查输入参数 |
| PFIT-1002 | 401 | 未授权/Token 失效 | 请重新登录 |
| PFIT-1003 | 403 | 无权限访问该资源 | 无权限访问 |
| PFIT-1004 | 429 | 请求频率超限 | 请求过于频繁，请稍后重试 |
| PFIT-1005 | 500 | 服务内部错误 | 系统繁忙，请稍后重试 |
| PFIT-2001 | 400 | 文件格式不支持（非 PDF） | 仅支持 PDF 格式文件 |
| PFIT-2002 | 400 | 文件大小超过限制（30MB） | 文件大小超过 30MB，请压缩后重传 |
| PFIT-2003 | 422 | PDF 无可提取文本（疑似扫描件） | 该文件可能为扫描件，暂不支持，请上传文字版 PDF |
| PFIT-2004 | 422 | PDF 解析失败（文件损坏或加密） | 文件解析失败，请检查文件是否加密或损坏 |
| PFIT-2005 | 422 | 提取文本内容过短（< 200 字符） | 文档内容过少，无法完成分析，请确认上传正确文件 |
| PFIT-3001 | 500 | LLM 调用超时 | 分析超时，系统将自动重试 |
| PFIT-3002 | 500 | LLM 返回非 JSON 格式 | 结构化分析失败，系统将自动重试 |
| PFIT-3003 | 500 | LLM 返回 JSON Schema 校验失败 | 分析结果异常，系统将自动重试 |
| PFIT-3004 | 422 | HealthFacts 抽取结果为空 | 未在体检报告中识别到支持的异常项，请确认报告内容 |
| PFIT-3005 | 422 | PolicyFacts 抽取结果为空 | 未在条款文档中识别到支持的条款类型，请确认文档内容 |
| PFIT-4001 | 500 | 规则引擎执行失败 | 风险评分异常，请重试 |
| PFIT-4002 | 422 | 无匹配风险主题 | 根据当前文档未发现明显风险冲突 |
| PFIT-5001 | 404 | 任务不存在 | 任务不存在或已被删除 |
| PFIT-5002 | 409 | 任务状态不允许该操作 | 当前任务状态不支持此操作 |
| PFIT-5003 | 400 | 任务缺少必要文档（体检报告/条款） | 请上传体检报告与保险条款后再开始分析 |

### 21.3 错误响应格式

```json
{
  "code": "PFIT-2003",
  "message": "该文件可能为扫描件，暂不支持，请上传文字版 PDF",
  "request_id": "req_abc123",
  "timestamp": "2026-02-28T10:00:00Z"
}
```

---

## 22. 前端交互设计细节

### 22.1 上传页交互流程

```
用户进入上传页
  │
  ├─ 拖拽/点击上传体检报告（必选）
  │    ├─ 校验：格式（PDF）、大小（≤30MB）
  │    ├─ 校验失败 → 内联错误提示（红色边框 + 提示文案）
  │    └─ 校验通过 → 显示文件名 + 文件大小 + 绿色勾
  │
  ├─ 拖拽/点击上传保险条款（必选）
  │    └─ 同上校验逻辑
  │
  ├─ 拖拽/点击上传投保告知（可选，有"跳过"链接）
  │
  ├─ 点击"开始分析"按钮
  │    ├─ 前端预校验（两个必选文件均已上传）
  │    └─ 提交 → 跳转任务处理页
  │
  └─ 文本直贴入口（折叠，默认隐藏，点击展开）
```

### 22.2 任务处理页状态流转

| 状态 | 展示内容 | 用户操作 |
|------|----------|----------|
| pending | 排队中，进度条（0%） | 可取消 |
| parsing | 正在读取文档，进度（30%） | 可取消 |
| extracting | 正在分析健康事实与条款，进度（60%） | 仅展示 |
| matching | 正在匹配风险规则，进度（90%） | 仅展示 |
| success | 分析完成，跳转报告页 | 查看报告 |
| failed | 显示失败原因（对应错误码文案）+ 重试按钮 | 重试 / 返回上传 |

**轮询策略**：前 30s 每 3s 轮询一次，30s 后每 10s 轮询一次，超过 5min 提示用户刷新。

### 22.3 报告总览页布局

```
┌──────────────────────────────────────────────────────┐
│  任务名称 / 创建时间 / 生成时间                        │
├──────────────────────────────────────────────────────┤
│  风险总览：🔴 2 高风险  🟡 3 待确认  🟢 5 暂无冲突   │
├──────────────────────────────────────────────────────┤
│  风险列表（按红→黄→绿排序）                            │
│  ┌──────────────────────────────────────────────────┐│
│  │ 🔴 高血压风险 · 可能触发既往症告知冲突            ││
│  │ 体检证据：血压 155/95 mmHg（2026-01-10）         ││
│  │ 条款证据：既往症定义第 3 条                       ││
│  │ [查看详情 →]                                      ││
│  └──────────────────────────────────────────────────┘│
│  ┌──────────────────────────────────────────────────┐│
│  │ 🟡 血脂异常 · 指标超标但诊断不明确               ││
│  │ ...                                               ││
│  └──────────────────────────────────────────────────┘│
├──────────────────────────────────────────────────────┤
│  [导出报告 PDF]  [删除任务]                            │
├──────────────────────────────────────────────────────┤
│  ⚠️ 免责声明：本报告仅为辅助解读，不构成理赔结论...   │
└──────────────────────────────────────────────────────┘
```

### 22.4 风险详情页布局

```
┌──────────────────────────────────────────────────────┐
│  [← 返回报告]   🔴 高血压风险                         │
├──────────────────────────────────────────────────────┤
│  风险说明                                              │
│  "体检报告显示血压偏高且有复查建议，与保单既往症定义   │
│   及投保告知要求可能存在冲突，建议补充核实。"          │
├──────────────────────────────────────────────────────┤
│  体检证据                                              │
│  ┌────────────────────────────────────────────────┐  │
│  │ [原文] 血压 155/95 mmHg，建议复查/就诊          │  │
│  │ 检查日期：2026-01-10  段落位置：第 12 段        │  │
│  └────────────────────────────────────────────────┘  │
├──────────────────────────────────────────────────────┤
│  条款证据                                              │
│  ┌────────────────────────────────────────────────┐  │
│  │ [原文] 既往症是指被保险人在本合同生效前已...    │  │
│  │ 所在位置：第 120 段（责任免除章节）             │  │
│  └────────────────────────────────────────────────┘  │
├──────────────────────────────────────────────────────┤
│  需要补充确认的问题                                     │
│  □ 是否已被医生明确诊断为高血压？                      │
│  □ 是否正在长期服用降压药物？                          │
│  □ 体检日期是否在本次投保前？                          │
├──────────────────────────────────────────────────────┤
│  建议行动                                              │
│  1. 补充投保告知截图，核实是否如实告知                  │
│  2. 联系保险公司核保确认除外或加费条件                  │
└──────────────────────────────────────────────────────┘
```

### 22.5 加载态与空状态设计

| 场景 | 设计 |
|------|------|
| 文件上传中 | 进度条 + 文件名 + 百分比 |
| 等待分析结果 | 骨架屏（Skeleton）+ 阶段文案 |
| 无风险发现（全绿） | 绿色图标 + "暂未发现高风险冲突，请注意以下待确认项" |
| 历史记录为空 | 引导插画 + "还没有分析记录，立即创建第一个分析" |
| 网络错误 | 错误图标 + 错误文案 + 重试按钮 |

---

## 23. 成本预估

### 23.1 LLM 调用成本预估

基于每次分析任务的典型调用量估算（以 GPT-4o / Claude 3.5 Sonnet 量级为参考）：

| 调用阶段 | 预估 Token 数（输入+输出） | 调用次数/任务 | 备注 |
|----------|---------------------------|---------------|------|
| HealthFacts 抽取 | ~3,000 tokens | 1-2 次 | 体检报告分段处理 |
| PolicyFacts 抽取 | ~8,000 tokens | 2-4 次 | 条款文档较长，需分批 |
| 风险解释生成 | ~1,000 tokens/条 | 3-8 次 | 按风险条数 |
| **单任务合计** | **~20,000-35,000 tokens** | | |

**费用估算（参考价格，以实际模型为准）**：

| 用量场景 | 月任务量 | 预估月费用（USD） |
|----------|----------|------------------|
| 个人测试 | 50 次 | < $5 |
| 小规模开源试用 | 500 次 | $25-50 |
| 中等规模 | 5,000 次 | $250-500 |

**成本控制策略**：
1. 分层模型：HealthFacts 抽取用小模型（如 GPT-4o-mini），风险解释用大模型。
2. 缓存中间结果：相同文档哈希跳过重复抽取。
3. 长文档分批处理，避免单次超长上下文。

### 23.2 基础设施成本预估（云部署）

| 组件 | 规格 | 月费用估算（USD） |
|------|------|------------------|
| API 服务（2C4G） | 单节点 | $20-40 |
| PostgreSQL（基础型） | 20GB | $15-30 |
| Redis | 1GB | $10-15 |
| MinIO / 对象存储 | 50GB | $5-10 |
| **合计** | | **$50-95/月** |

> 自部署（VPS）可控制在 $20-40/月；个人开发阶段可用 Docker Compose 本地运行，成本为零。

---

## 24. 开源协议与贡献指南

### 24.1 开源协议

推荐使用 **Apache License 2.0**，理由：
1. 允许商业使用，方便用户在企业内部部署。
2. 要求保留版权声明，保护项目品牌。
3. 与 Go 生态主流项目（如 Kubernetes、Gin）保持一致。

> 备选方案：MIT（更宽松）；若希望防止闭源分发，可选 AGPL-3.0。

### 24.2 仓库基础结构

```text
policy-fit/
  .github/
    ISSUE_TEMPLATE/
      bug_report.md
      feature_request.md
    PULL_REQUEST_TEMPLATE.md
    workflows/
      ci.yml          # Go test + lint
  docs/
    PRD-保单避坑雷达-v1.md
  LICENSE
  README.md
  CONTRIBUTING.md
  CHANGELOG.md
```

### 24.3 README 必要内容

1. 项目一句话介绍与截图/演示 GIF。
2. 功能特性列表（红黄绿风险报告、证据对照、追问清单）。
3. 快速启动（Docker Compose 一键运行）。
4. 环境变量配置说明（LLM API Key、数据库连接）。
5. 贡献指南链接。
6. 免责声明（中英文）。

### 24.4 贡献规范（CONTRIBUTING.md）

1. **Issue 先行**：新功能/大改动先开 Issue 讨论，小 bug 直接 PR。
2. **分支规范**：`feat/xxx`、`fix/xxx`、`docs/xxx`、`refactor/xxx`。
3. **Commit 规范**：遵循 Conventional Commits（`feat:`、`fix:`、`docs:` 等）。
4. **测试要求**：新功能需附带单元测试，覆盖率不低于当前基线。
5. **代码风格**：Go 代码通过 `gofmt` + `golangci-lint`；前端通过 ESLint + Prettier。
6. **PR 模板**：描述变更内容、关联 Issue、测试方式、截图（如有 UI 变更）。

### 24.5 版本发布规范

遵循语义化版本（SemVer）：`MAJOR.MINOR.PATCH`

| 版本类型 | 触发条件 | 示例 |
|----------|----------|------|
| PATCH | Bug 修复、文档更新 | 1.0.1 |
| MINOR | 新功能（向后兼容） | 1.1.0 |
| MAJOR | 破坏性变更、架构调整 | 2.0.0 |

每个 Release 需包含：CHANGELOG 条目、Docker 镜像 Tag、迁移脚本（如有 DB 变更）。

