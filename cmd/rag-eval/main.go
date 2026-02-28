package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/zhenglizhi/policy-fit/internal/rag"
)

type evalCase struct {
	PolicyID        string `json:"policy_id"`
	PolicyText      string `json:"policy_text"`
	Query           string `json:"query"`
	ExpectedKeyword string `json:"expected_keyword"`
}

func main() {
	cases, err := loadCases("testdata/rag/eval_cases.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "load eval cases: %v\n", err)
		os.Exit(1)
	}

	store := rag.NewInMemoryStore()
	embedder := rag.NewHashEmbeddingProvider(32)
	retriever := rag.NewRetriever(store, embedder)

	ctx := context.Background()
	for _, c := range cases {
		if err := retriever.IndexPolicyText(ctx, c.PolicyID, c.PolicyText); err != nil {
			fmt.Fprintf(os.Stderr, "index policy %s: %v\n", c.PolicyID, err)
			os.Exit(1)
		}
	}

	hitCount := 0
	falsePositive := 0
	for _, c := range cases {
		results, err := retriever.Retrieve(ctx, c.Query, 3)
		if err != nil {
			fmt.Fprintf(os.Stderr, "retrieve %s: %v\n", c.Query, err)
			os.Exit(1)
		}
		if len(results) == 0 {
			continue
		}
		top := strings.ToLower(results[0].Chunk.Text)
		expected := strings.ToLower(c.ExpectedKeyword)
		if strings.Contains(top, expected) {
			hitCount++
		} else {
			falsePositive++
		}
	}

	recall := float64(hitCount)
	if len(cases) > 0 {
		recall = recall / float64(len(cases))
	}
	falsePositiveRate := 0.0
	if len(cases) > 0 {
		falsePositiveRate = float64(falsePositive) / float64(len(cases))
	}

	fmt.Printf("RAG offline eval\n")
	fmt.Printf("cases=%d hit=%d recall=%.3f false_positive_rate=%.3f\n", len(cases), hitCount, recall, falsePositiveRate)
}

func loadCases(path string) ([]evalCase, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cases []evalCase
	if err := json.Unmarshal(raw, &cases); err != nil {
		return nil, err
	}
	return cases, nil
}
