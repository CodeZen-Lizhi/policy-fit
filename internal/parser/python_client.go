package parser

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	// ErrPythonService Python 解析服务调用失败
	ErrPythonService = errors.New("python parser service error")
)

// PythonServiceClient Python 解析服务客户端
type PythonServiceClient struct {
	baseURL string
	client  *http.Client
	retries int
}

// ParseDocumentRequest 解析文档请求
type ParseDocumentRequest struct {
	Filename      string `json:"filename"`
	MimeType      string `json:"mime_type"`
	ContentBase64 string `json:"content_base64,omitempty"`
	RawText       string `json:"raw_text,omitempty"`
	EnableOCR     bool   `json:"enable_ocr"`
}

// ParseDocumentResponse 解析文档响应
type ParseDocumentResponse struct {
	Text         string      `json:"text"`
	Paragraphs   []Paragraph `json:"paragraphs"`
	QualityScore float64     `json:"quality_score"`
	Hints        []string    `json:"hints"`
}

// ParseReportResponse 体检结构化响应
type ParseReportResponse struct {
	Facts        []map[string]interface{} `json:"facts"`
	QualityScore float64                  `json:"quality_score"`
	Hints        []string                 `json:"hints"`
}

// ParsePolicyResponse 条款结构化响应
type ParsePolicyResponse struct {
	Sections     []map[string]interface{} `json:"sections"`
	QualityScore float64                  `json:"quality_score"`
	Hints        []string                 `json:"hints"`
}

// NewPythonServiceClient 创建 Python 解析服务客户端
func NewPythonServiceClient(baseURL string, timeout time.Duration, retries int) *PythonServiceClient {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	if retries < 0 {
		retries = 0
	}
	return &PythonServiceClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: timeout,
		},
		retries: retries,
	}
}

// ParseFile 解析本地文件
func (c *PythonServiceClient) ParseFile(ctx context.Context, filePath string, enableOCR bool) (*ParseDocumentResponse, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: read file: %v", ErrPythonService, err)
	}
	req := ParseDocumentRequest{
		Filename:      filepath.Base(filePath),
		MimeType:      detectMimeType(filePath),
		ContentBase64: base64.StdEncoding.EncodeToString(content),
		EnableOCR:     enableOCR,
	}
	return c.ParseDocument(ctx, req)
}

// ParseDocument 调用 /parse/document
func (c *PythonServiceClient) ParseDocument(ctx context.Context, req ParseDocumentRequest) (*ParseDocumentResponse, error) {
	if strings.TrimSpace(c.baseURL) == "" {
		return nil, fmt.Errorf("%w: empty python parser base url", ErrPythonService)
	}
	if req.MimeType == "" {
		req.MimeType = "application/pdf"
	}
	if req.Filename == "" {
		req.Filename = "document"
	}
	if !req.EnableOCR {
		// keep false
	} else {
		req.EnableOCR = true
	}

	var resp ParseDocumentResponse
	if err := c.postJSON(ctx, "/parse/document", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ParseReport 调用 /parse/report
func (c *PythonServiceClient) ParseReport(ctx context.Context, text string) (*ParseReportResponse, error) {
	var resp ParseReportResponse
	if err := c.postJSON(ctx, "/parse/report", map[string]string{"text": text}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ParsePolicy 调用 /parse/policy
func (c *PythonServiceClient) ParsePolicy(ctx context.Context, text string) (*ParsePolicyResponse, error) {
	var resp ParsePolicyResponse
	if err := c.postJSON(ctx, "/parse/policy", map[string]string{"text": text}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *PythonServiceClient) postJSON(ctx context.Context, path string, requestBody interface{}, out interface{}) error {
	if strings.TrimSpace(c.baseURL) == "" {
		return fmt.Errorf("%w: empty python parser base url", ErrPythonService)
	}

	raw, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("%w: marshal request: %v", ErrPythonService, err)
	}

	var lastErr error
	for i := 0; i <= c.retries; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(raw))
		if err != nil {
			return fmt.Errorf("%w: create request: %v", ErrPythonService, err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.client.Do(req)
		if err != nil {
			if isRetryableNetworkErr(err) && i < c.retries {
				lastErr = err
				continue
			}
			return fmt.Errorf("%w: %v", ErrPythonService, err)
		}

		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if resp.StatusCode >= 500 && i < c.retries {
			lastErr = fmt.Errorf("status=%d body=%s", resp.StatusCode, string(body))
			continue
		}

		if resp.StatusCode >= 300 {
			return fmt.Errorf("%w: status=%d body=%s", ErrPythonService, resp.StatusCode, string(body))
		}

		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("%w: unmarshal response: %v", ErrPythonService, err)
		}
		return nil
	}

	if lastErr != nil {
		return fmt.Errorf("%w: %v", ErrPythonService, lastErr)
	}
	return fmt.Errorf("%w: unknown request failure", ErrPythonService)
}

func detectMimeType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".pdf" {
		return "application/pdf"
	}
	if byExt := mime.TypeByExtension(ext); byExt != "" {
		return byExt
	}
	return "application/octet-stream"
}

func isRetryableNetworkErr(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}
	return false
}
