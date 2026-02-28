package llm

import (
	"context"
	"errors"
	"testing"
)

func TestExtractHealthFacts(t *testing.T) {
	client := NewClientWithProvider(mockProvider(func(context.Context, string) (string, error) {
		return `{"facts":[{"category":"hypertension","label":"高血压风险","evidence":{"text":"血压 155/95","loc":"para_1"},"confidence":0.9}]}`, nil
	}))

	facts, err := client.ExtractHealthFacts(context.Background(), "dummy")
	if err != nil {
		t.Fatalf("ExtractHealthFacts error: %v", err)
	}
	if len(facts) != 1 {
		t.Fatalf("unexpected facts length: %d", len(facts))
	}
	if facts[0].Evidence.Date != "unknown" {
		t.Fatalf("expected unknown date, got %s", facts[0].Evidence.Date)
	}
}

func TestExtractPolicyFacts(t *testing.T) {
	client := NewClientWithProvider(mockProvider(func(context.Context, string) (string, error) {
		return `{"sections":[{"type":"exclusion","title":"责任免除","content":"xxx","loc":"para_2","confidence":0.88}]}`, nil
	}))

	sections, err := client.ExtractPolicyFacts(context.Background(), "dummy")
	if err != nil {
		t.Fatalf("ExtractPolicyFacts error: %v", err)
	}
	if len(sections) != 1 {
		t.Fatalf("unexpected sections length: %d", len(sections))
	}
}

func TestExtractHealthFactsProviderError(t *testing.T) {
	client := NewClientWithProvider(mockProvider(func(context.Context, string) (string, error) {
		return "", errors.New("provider failed")
	}))
	_, err := client.ExtractHealthFacts(context.Background(), "dummy")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestExtractPolicyFactsInvalidJSON(t *testing.T) {
	client := NewClientWithProvider(mockProvider(func(context.Context, string) (string, error) {
		return `{"sections":[`, nil
	}))
	_, err := client.ExtractPolicyFacts(context.Background(), "dummy")
	if !errors.Is(err, ErrInvalidJSON) {
		t.Fatalf("expected ErrInvalidJSON, got %v", err)
	}
}

type mockProvider func(ctx context.Context, prompt string) (string, error)

func (m mockProvider) CompleteJSON(ctx context.Context, prompt string) (string, error) {
	return m(ctx, prompt)
}
