package ruleengine

import (
	"testing"

	"github.com/zhenglizhi/policy-fit/internal/domain"
)

func TestMatcherMatchByCategoryAndPolicyType(t *testing.T) {
	cfg := &TopicsConfig{
		Topics: map[string]TopicRule{
			"hypertension": {
				HitKeywords: []string{"高血压"},
				PolicyTypes: []string{"exclusion"},
			},
		},
	}
	matcher := NewMatcher(cfg)

	results := matcher.Match(
		[]domain.HealthFact{
			{
				Category: "hypertension",
				Label:    "高血压",
				Evidence: domain.EvidenceDetail{Text: "血压偏高", Loc: "para_1", Date: "2026-01-01"},
			},
		},
		[]domain.PolicyFact{
			{Type: "exclusion", Content: "免责条款", Loc: "para_100"},
			{Type: "renewal", Content: "续保条款", Loc: "para_120"},
		},
	)

	if len(results) != 1 {
		t.Fatalf("expected 1 match, got %d", len(results))
	}
	if len(results[0].PolicyFacts) != 1 {
		t.Fatalf("expected 1 matched policy fact, got %d", len(results[0].PolicyFacts))
	}
	if results[0].PolicyFacts[0].Type != "exclusion" {
		t.Fatalf("unexpected policy type: %s", results[0].PolicyFacts[0].Type)
	}
}
