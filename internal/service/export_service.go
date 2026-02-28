package service

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/zhenglizhi/policy-fit/internal/domain"
	"github.com/zhenglizhi/policy-fit/internal/i18n"
)

// ExportService 报告导出服务
type ExportService struct {
	taskSvc    exportTaskReader
	findingSvc exportFindingReader
}

type exportTaskReader interface {
	GetTask(ctx context.Context, taskID int64, actorID int64) (*domain.AnalysisTask, error)
}

type exportFindingReader interface {
	ListByTask(ctx context.Context, taskID int64) ([]domain.RiskFinding, error)
}

// NewExportService 创建导出服务
func NewExportService(taskSvc exportTaskReader, findingSvc exportFindingReader) *ExportService {
	return &ExportService{
		taskSvc:    taskSvc,
		findingSvc: findingSvc,
	}
}

// ExportMarkdown 导出 Markdown 报告
func (s *ExportService) ExportMarkdown(ctx context.Context, taskID int64, actorID int64, lang string) ([]byte, error) {
	task, err := s.taskSvc.GetTask(ctx, taskID, actorID)
	if err != nil {
		return nil, err
	}
	findings, err := s.findingSvc.ListByTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	var b strings.Builder
	generatedAt := time.Now().Format(time.RFC3339)
	labels := resolveExportLabels(lang)
	b.WriteString(fmt.Sprintf("# %s\n\n", labels.ReportTitle))
	b.WriteString(fmt.Sprintf("- %s: %d\n", labels.TaskID, task.ID))
	b.WriteString(fmt.Sprintf("- %s: v1.0\n", labels.ReportVersion))
	b.WriteString(fmt.Sprintf("- %s: %s\n\n", labels.GeneratedAt, generatedAt))
	b.WriteString(fmt.Sprintf("## %s\n\n", labels.RiskSummary))
	summary := map[string]int{"red": 0, "yellow": 0, "green": 0}
	for _, finding := range findings {
		switch finding.Level {
		case domain.RiskLevelRed:
			summary["red"]++
		case domain.RiskLevelYellow:
			summary["yellow"]++
		case domain.RiskLevelGreen:
			summary["green"]++
		}
	}
	b.WriteString(fmt.Sprintf("- 🔴 %s: %d\n- 🟡 %s: %d\n- 🟢 %s: %d\n\n",
		labels.RedRisk, summary["red"],
		labels.YellowRisk, summary["yellow"],
		labels.GreenRisk, summary["green"],
	))
	b.WriteString(fmt.Sprintf("## %s\n\n", labels.RiskDetail))
	for i, finding := range findings {
		b.WriteString(fmt.Sprintf("### %d. %s (%s)\n", i+1, finding.Topic, finding.Level))
		b.WriteString(fmt.Sprintf("- %s: %s\n", labels.Description, finding.Summary))
		if len(finding.HealthEvidence) > 0 {
			b.WriteString(fmt.Sprintf("- %s: %s (%s)\n", labels.HealthEvidence, finding.HealthEvidence[0].Text, finding.HealthEvidence[0].Loc))
		}
		if len(finding.PolicyEvidence) > 0 {
			b.WriteString(fmt.Sprintf("- %s: %s (%s)\n", labels.PolicyEvidence, finding.PolicyEvidence[0].Text, finding.PolicyEvidence[0].Loc))
		}
		if len(finding.Questions) > 0 {
			b.WriteString(fmt.Sprintf("- %s:\n", labels.Questions))
			for _, q := range finding.Questions {
				b.WriteString(fmt.Sprintf("  - %s\n", q))
			}
		}
		b.WriteString("\n")
	}
	b.WriteString(fmt.Sprintf("## %s\n\n", labels.DisclaimerTitle))
	b.WriteString(labels.DisclaimerBody + "\n")

	return []byte(b.String()), nil
}

// ExportPDF 导出 PDF 报告
func (s *ExportService) ExportPDF(ctx context.Context, taskID int64, actorID int64, lang string) ([]byte, error) {
	md, err := s.ExportMarkdown(ctx, taskID, actorID, lang)
	if err != nil {
		return nil, err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)
	pdf.AddPage()
	for _, line := range strings.Split(string(md), "\n") {
		pdf.MultiCell(0, 6, line, "", "L", false)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type exportLabels struct {
	ReportTitle     string
	TaskID          string
	ReportVersion   string
	GeneratedAt     string
	RiskSummary     string
	RedRisk         string
	YellowRisk      string
	GreenRisk       string
	RiskDetail      string
	Description     string
	HealthEvidence  string
	PolicyEvidence  string
	Questions       string
	DisclaimerTitle string
	DisclaimerBody  string
}

func resolveExportLabels(lang string) exportLabels {
	if i18n.ResolveLanguage(lang) == "en-US" {
		return exportLabels{
			ReportTitle:     "Policy Fit Report",
			TaskID:          "Task ID",
			ReportVersion:   "Report Version",
			GeneratedAt:     "Generated At",
			RiskSummary:     "Risk Summary",
			RedRisk:         "High Risk",
			YellowRisk:      "Need Confirmation",
			GreenRisk:       "No Significant Conflict",
			RiskDetail:      "Risk Findings",
			Description:     "Description",
			HealthEvidence:  "Health Evidence",
			PolicyEvidence:  "Policy Evidence",
			Questions:       "Follow-up Questions",
			DisclaimerTitle: "Disclaimer",
			DisclaimerBody:  "This report is for policy interpretation and risk hints only. It is not an underwriting or claim decision.",
		}
	}
	return exportLabels{
		ReportTitle:     "保单避坑雷达报告",
		TaskID:          "任务ID",
		ReportVersion:   "报告版本",
		GeneratedAt:     "生成时间",
		RiskSummary:     "风险摘要",
		RedRisk:         "高风险",
		YellowRisk:      "待确认",
		GreenRisk:       "暂无冲突",
		RiskDetail:      "风险明细",
		Description:     "说明",
		HealthEvidence:  "体检证据",
		PolicyEvidence:  "条款证据",
		Questions:       "追问",
		DisclaimerTitle: "免责声明",
		DisclaimerBody:  "本报告仅用于条款辅助解读与风险提示，不构成承保或理赔结论。",
	}
}
