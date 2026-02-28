package domain

import "time"

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusParsing    TaskStatus = "parsing"
	TaskStatusExtracting TaskStatus = "extracting"
	TaskStatusMatching   TaskStatus = "matching"
	TaskStatusSuccess    TaskStatus = "success"
	TaskStatusFailed     TaskStatus = "failed"
)

// AnalysisTask 分析任务
type AnalysisTask struct {
	ID          int64          `json:"id"`
	UserID      int64          `json:"user_id"`
	RequestID   string         `json:"request_id,omitempty"`
	Status      TaskStatus     `json:"status"`
	RiskSummary map[string]int `json:"risk_summary,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// DocumentType 文档类型
type DocumentType string

const (
	DocTypeReport     DocumentType = "report"
	DocTypePolicy     DocumentType = "policy"
	DocTypeDisclosure DocumentType = "disclosure"
)

// ParseStatus 解析状态
type ParseStatus string

const (
	ParseStatusPending ParseStatus = "pending"
	ParseStatusSuccess ParseStatus = "success"
	ParseStatusFailed  ParseStatus = "failed"
)

// Document 文档
type Document struct {
	ID          int64        `json:"id"`
	TaskID      int64        `json:"task_id"`
	DocType     DocumentType `json:"doc_type"`
	FileName    string       `json:"file_name"`
	StorageKey  string       `json:"storage_key"`
	ParseStatus ParseStatus  `json:"parse_status"`
	ParsedText  string       `json:"parsed_text,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
}

// RiskLevel 风险等级
type RiskLevel string

const (
	RiskLevelRed    RiskLevel = "red"
	RiskLevelYellow RiskLevel = "yellow"
	RiskLevelGreen  RiskLevel = "green"
)

// RiskFinding 风险发现
type RiskFinding struct {
	ID             int64      `json:"id"`
	TaskID         int64      `json:"task_id"`
	Level          RiskLevel  `json:"level"`
	Topic          string     `json:"topic"`
	Summary        string     `json:"summary"`
	HealthEvidence []Evidence `json:"health_evidence"`
	PolicyEvidence []Evidence `json:"policy_evidence"`
	Questions      []string   `json:"questions"`
	Actions        []string   `json:"actions,omitempty"`
	Confidence     float64    `json:"confidence"`
	CreatedAt      time.Time  `json:"created_at"`
}

// Evidence 证据
type Evidence struct {
	Loc  string `json:"loc"`
	Text string `json:"text"`
}

// HealthFact 健康事实
type HealthFact struct {
	Category           string                 `json:"category"`
	Label              string                 `json:"label"`
	Evidence           EvidenceDetail         `json:"evidence"`
	Values             map[string]interface{} `json:"values,omitempty"`
	Diagnosed          *bool                  `json:"diagnosed,omitempty"`
	LongTermMedication *bool                  `json:"long_term_medication,omitempty"`
	Confidence         float64                `json:"confidence"`
	UncertainReason    string                 `json:"uncertain_reason,omitempty"`
}

// EvidenceDetail 证据详情
type EvidenceDetail struct {
	Text   string `json:"text"`
	Date   string `json:"date"`
	Loc    string `json:"loc"`
	Source string `json:"source"`
}

// PolicyFact 条款事实
type PolicyFact struct {
	Type       string   `json:"type"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Loc        string   `json:"loc"`
	Confidence float64  `json:"confidence"`
	Questions  []string `json:"questions,omitempty"`
}

// AuditLog 审计日志
type AuditLog struct {
	ID         int64                  `json:"id"`
	TaskID     *int64                 `json:"task_id,omitempty"`
	ActorID    *int64                 `json:"actor_id,omitempty"`
	Action     string                 `json:"action"`
	TargetType string                 `json:"target_type"`
	TargetID   string                 `json:"target_id,omitempty"`
	Detail     map[string]interface{} `json:"detail"`
	CreatedAt  time.Time              `json:"created_at"`
}

// RuleConfigVersion 规则配置版本
type RuleConfigVersion struct {
	ID        int64                  `json:"id"`
	Version   string                 `json:"version"`
	Changelog string                 `json:"changelog"`
	Content   map[string]interface{} `json:"content"`
	IsActive  bool                   `json:"is_active"`
	IsGray    bool                   `json:"is_gray"`
	CreatedBy *int64                 `json:"created_by,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

// AnalyticsEvent 运营埋点事件
type AnalyticsEvent struct {
	ID         int64                  `json:"id"`
	UserID     *int64                 `json:"user_id,omitempty"`
	TaskID     *int64                 `json:"task_id,omitempty"`
	EventName  string                 `json:"event_name"`
	Properties map[string]interface{} `json:"properties"`
	CreatedAt  time.Time              `json:"created_at"`
}
