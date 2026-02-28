package rag

import (
	"context"
	"testing"
)

func TestRetrieverIndexAndRetrieve(t *testing.T) {
	store := NewInMemoryStore()
	embedder := NewHashEmbeddingProvider(32)
	r := NewRetriever(store, embedder)

	policyText := "第1条 既往症定义，投保前疾病不赔。\n第2条 等待期 90 天。"
	if err := r.IndexPolicyText(context.Background(), "policy-1", policyText); err != nil {
		t.Fatalf("IndexPolicyText error: %v", err)
	}

	results, err := r.Retrieve(context.Background(), "既往症 是否赔付", 3)
	if err != nil {
		t.Fatalf("Retrieve error: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("expected non-empty retrieval results")
	}
	if results[0].Chunk.Text == "" {
		t.Fatalf("retrieved chunk text should not be empty")
	}
}
