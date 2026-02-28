package parser

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var (
	// ErrUnreadablePDF PDF 不可读或命令失败
	ErrUnreadablePDF = errors.New("pdf unreadable")
	// ErrEmptyText 解析文本为空
	ErrEmptyText = errors.New("pdf extracted text is empty")
)

// Paragraph 段落结构
type Paragraph struct {
	Loc    string `json:"loc"`
	Page   int    `json:"page"`
	Index  int    `json:"index"`
	Text   string `json:"text"`
	Source string `json:"source"`
}

// ParseResult 解析结果
type ParseResult struct {
	Text         string      `json:"text"`
	Paragraphs   []Paragraph `json:"paragraphs"`
	QualityScore float64     `json:"quality_score,omitempty"`
	Hints        []string    `json:"hints,omitempty"`
}

// PDFParser PDF 解析器
type PDFParser struct {
	binaryPath string
	runner     func(ctx context.Context, binaryPath string, filePath string) (string, error)
}

// NewPDFParser 创建 PDF 解析器
func NewPDFParser(binaryPath string) *PDFParser {
	if strings.TrimSpace(binaryPath) == "" {
		binaryPath = "pdftotext"
	}
	return &PDFParser{
		binaryPath: binaryPath,
		runner:     runPDFToText,
	}
}

// Parse 解析 PDF 文件
func (p *PDFParser) Parse(ctx context.Context, filePath string) (*ParseResult, error) {
	if strings.TrimSpace(filePath) == "" {
		return nil, fmt.Errorf("%w: empty file path", ErrUnreadablePDF)
	}

	text, err := p.runner(ctx, p.binaryPath, filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnreadablePDF, err)
	}

	normalized := strings.TrimSpace(text)
	if normalized == "" {
		return nil, ErrEmptyText
	}

	paragraphs := splitParagraphs(text)
	if len(paragraphs) == 0 {
		return nil, ErrEmptyText
	}

	return &ParseResult{
		Text:         text,
		Paragraphs:   paragraphs,
		QualityScore: estimateQualityScore(paragraphs, normalized),
	}, nil
}

func runPDFToText(ctx context.Context, binaryPath string, filePath string) (string, error) {
	cmd := exec.CommandContext(ctx, binaryPath, "-layout", filePath, "-")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("exec %s failed: %w, output=%s", binaryPath, err, string(output))
	}
	return string(output), nil
}

func splitParagraphs(text string) []Paragraph {
	pages := strings.Split(text, "\f")
	paragraphs := make([]Paragraph, 0)
	paraIdx := 0

	for pageIdx, pageText := range pages {
		pageNumber := pageIdx + 1
		blocks := splitTextBlocks(pageText)
		for _, block := range blocks {
			trimmed := strings.TrimSpace(block)
			if trimmed == "" {
				continue
			}
			paraIdx++
			paragraphs = append(paragraphs, Paragraph{
				Loc:    fmt.Sprintf("para_%d", paraIdx),
				Page:   pageNumber,
				Index:  paraIdx,
				Text:   trimmed,
				Source: "report",
			})
		}
	}

	return paragraphs
}

func splitTextBlocks(text string) []string {
	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")

	var result []string
	var current strings.Builder
	for _, line := range strings.Split(normalized, "\n") {
		if strings.TrimSpace(line) == "" {
			if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
			continue
		}
		if current.Len() > 0 {
			current.WriteByte('\n')
		}
		current.WriteString(line)
	}
	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

func estimateQualityScore(paragraphs []Paragraph, normalizedText string) float64 {
	if len(paragraphs) == 0 || normalizedText == "" {
		return 0
	}
	textScore := float64(len(normalizedText))
	if textScore > 2000 {
		textScore = 2000
	}
	paraScore := float64(len(paragraphs))
	if paraScore > 20 {
		paraScore = 20
	}
	return ((textScore / 2000.0) * 0.7) + ((paraScore / 20.0) * 0.3)
}
