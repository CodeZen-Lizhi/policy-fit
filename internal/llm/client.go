package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/zhenglizhi/policy-fit/internal/config"
	"github.com/zhenglizhi/policy-fit/internal/domain"
)

var (
	// ErrProviderTimeout 供应商调用超时
	ErrProviderTimeout = errors.New("llm provider timeout")
	// ErrProviderResponse 供应商响应异常
	ErrProviderResponse = errors.New("llm provider response error")
	// ErrInvalidJSON LLM 返回 JSON 非法
	ErrInvalidJSON = errors.New("llm invalid json")
	// ErrSchemaValidation JSON Schema 校验失败
	ErrSchemaValidation = errors.New("llm schema validation failed")
)

// Provider LLM Provider 抽象
type Provider interface {
	CompleteJSON(ctx context.Context, prompt string) (string, error)
}

// Client LLM 客户端
type Client struct {
	provider Provider
}

// NewClient 创建 LLM 客户端
func NewClient(cfg config.LLMConfig) (*Client, error) {
	provider, err := newProvider(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{provider: provider}, nil
}

// NewClientWithProvider 测试/扩展用构造函数
func NewClientWithProvider(provider Provider) *Client {
	return &Client{provider: provider}
}

// ExtractHealthFacts 提取健康事实
func (c *Client) ExtractHealthFacts(ctx context.Context, reportText string) ([]domain.HealthFact, error) {
	prompt, err := renderPrompt("health_facts.tmpl", map[string]string{
		"REPORT_TEXT": reportText,
	})
	if err != nil {
		return nil, err
	}

	raw, err := c.provider.CompleteJSON(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var payload struct {
		Facts []domain.HealthFact `json:"facts"`
	}
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}
	if err := NormalizeAndValidateHealthFacts(payload.Facts); err != nil {
		return nil, err
	}
	return payload.Facts, nil
}

// ExtractPolicyFacts 提取条款事实
func (c *Client) ExtractPolicyFacts(ctx context.Context, policyText string) ([]domain.PolicyFact, error) {
	prompt, err := renderPrompt("policy_facts.tmpl", map[string]string{
		"POLICY_TEXT": policyText,
	})
	if err != nil {
		return nil, err
	}

	raw, err := c.provider.CompleteJSON(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var payload struct {
		Sections []domain.PolicyFact `json:"sections"`
	}
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}
	if err := NormalizeAndValidatePolicyFacts(payload.Sections); err != nil {
		return nil, err
	}
	return payload.Sections, nil
}

func renderPrompt(name string, data map[string]string) (string, error) {
	raw, ok := PromptTemplates[name]
	if !ok {
		return "", fmt.Errorf("prompt template not found: %s", name)
	}
	tpl, err := template.New(name).Parse(raw)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func newProvider(cfg config.LLMConfig) (Provider, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case "openai":
		return newOpenAIProvider(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported llm provider: %s", cfg.Provider)
	}
}

type openAIProvider struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
	retries int
}

func newOpenAIProvider(cfg config.LLMConfig) Provider {
	timeout := time.Duration(cfg.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 120 * time.Second
	}
	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	return &openAIProvider{
		apiKey:  cfg.APIKey,
		baseURL: baseURL,
		model:   cfg.Model,
		client: &http.Client{
			Timeout: timeout,
		},
		retries: 2,
	}
}

func (p *openAIProvider) CompleteJSON(ctx context.Context, prompt string) (string, error) {
	reqBody := map[string]interface{}{
		"model": p.model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"response_format": map[string]string{
			"type": "json_object",
		},
	}

	rawReq, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	var lastErr error
	for i := 0; i <= p.retries; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/chat/completions", bytes.NewReader(rawReq))
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := p.client.Do(req)
		if err != nil {
			if isTimeoutError(err) {
				lastErr = fmt.Errorf("%w: %v", ErrProviderTimeout, err)
			} else {
				lastErr = fmt.Errorf("%w: %v", ErrProviderResponse, err)
			}
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode >= 500 && i < p.retries {
			lastErr = fmt.Errorf("%w: http=%d body=%s", ErrProviderResponse, resp.StatusCode, string(body))
			continue
		}
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("%w: http=%d body=%s", ErrProviderResponse, resp.StatusCode, string(body))
		}

		var parsed struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}
		if err := json.Unmarshal(body, &parsed); err != nil {
			return "", fmt.Errorf("%w: %v", ErrInvalidJSON, err)
		}
		if len(parsed.Choices) == 0 {
			return "", fmt.Errorf("%w: empty choices", ErrProviderResponse)
		}
		return parsed.Choices[0].Message.Content, nil
	}

	if lastErr != nil {
		return "", lastErr
	}
	return "", ErrProviderResponse
}

func isTimeoutError(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}
