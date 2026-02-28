package parser

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zhenglizhi/policy-fit/internal/config"
)

// DocumentParser 通用文档解析抽象
type DocumentParser interface {
	Parse(ctx context.Context, filePath string) (*ParseResult, error)
}

type pythonServiceAdapter struct {
	client *PythonServiceClient
}

func (p *pythonServiceAdapter) Parse(ctx context.Context, filePath string) (*ParseResult, error) {
	resp, err := p.client.ParseFile(ctx, filePath, true)
	if err != nil {
		return nil, err
	}
	return &ParseResult{
		Text:         resp.Text,
		Paragraphs:   resp.Paragraphs,
		QualityScore: resp.QualityScore,
		Hints:        resp.Hints,
	}, nil
}

// NewConfiguredParser 按配置创建解析器（pdftotext/python-service）
func NewConfiguredParser(cfg config.ParserConfig) (DocumentParser, error) {
	switch strings.TrimSpace(cfg.PDFParser) {
	case "", "pdftotext":
		return NewPDFParser(cfg.PDFParser), nil
	case "python-service":
		if strings.TrimSpace(cfg.PythonServiceURL) == "" {
			return nil, fmt.Errorf("%w: PYTHON_SERVICE_URL is empty", ErrPythonService)
		}
		client := NewPythonServiceClient(cfg.PythonServiceURL, 30*time.Second, 2)
		return &pythonServiceAdapter{client: client}, nil
	default:
		return nil, fmt.Errorf("unsupported parser type: %s", cfg.PDFParser)
	}
}
