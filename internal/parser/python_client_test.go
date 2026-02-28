package parser

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/zhenglizhi/policy-fit/internal/config"
)

func TestPythonServiceClientParseDocumentWithRetry(t *testing.T) {
	attempts := 0
	client := NewPythonServiceClient("http://python-parser", 2*time.Second, 2)
	client.client = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		attempts++
		if attempts == 1 {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader(`{"error":"temporary"}`)),
				Header:     make(http.Header),
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"text":"ok","paragraphs":[{"loc":"para_1","page":1,"index":1,"text":"ok","source":"report"}],"quality_score":0.88,"hints":[]}`)),
			Header:     make(http.Header),
		}, nil
	})}

	resp, err := client.ParseDocument(context.Background(), ParseDocumentRequest{
		Filename:  "sample.pdf",
		MimeType:  "application/pdf",
		RawText:   "mock text",
		EnableOCR: true,
	})
	if err != nil {
		t.Fatalf("ParseDocument error: %v", err)
	}
	if resp.Text != "ok" {
		t.Fatalf("unexpected text: %s", resp.Text)
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

func TestNewConfiguredParserPythonService(t *testing.T) {
	p, err := NewConfiguredParser(config.ParserConfig{
		PDFParser:        "python-service",
		PythonServiceURL: "http://localhost:8081",
	})
	if err != nil {
		t.Fatalf("NewConfiguredParser error: %v", err)
	}
	if p == nil {
		t.Fatalf("parser should not be nil")
	}
}

func TestNewConfiguredParserPythonServiceMissingURL(t *testing.T) {
	_, err := NewConfiguredParser(config.ParserConfig{
		PDFParser: "python-service",
	})
	if err == nil {
		t.Fatalf("expected error for missing PYTHON_SERVICE_URL")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
