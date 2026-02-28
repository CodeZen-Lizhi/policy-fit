package ruleengine

import (
	"testing"

	"github.com/zhenglizhi/policy-fit/internal/domain"
)

func TestScorerCasesForAllTopics(t *testing.T) {
	cfg, err := LoadTopicsConfig("../../configs/topics.yaml")
	if err != nil {
		t.Fatalf("load topics config: %v", err)
	}

	scorer := NewScorer()
	for topic := range cfg.Topics {
		rule := cfg.Topics[topic]

		diagnosed := true
		notDiagnosed := false
		longMed := true

		redMatch := MatchResult{
			Topic: topic,
			Rule:  rule,
			HealthFact: domain.HealthFact{
				Category:   topic,
				Label:      topic,
				Evidence:   domain.EvidenceDetail{Loc: "para_1", Text: "abnormal", Date: "2026-01-01"},
				Diagnosed:  &diagnosed,
				Confidence: 0.95,
			},
			PolicyFacts: []domain.PolicyFact{
				{Type: firstPolicyType(rule), Loc: "para_100", Content: "条款"},
			},
		}
		redFinding := scorer.Score([]MatchResult{redMatch})
		if len(redFinding) != 1 || redFinding[0].Level != domain.RiskLevelRed {
			t.Fatalf("topic=%s expected red, got %#v", topic, redFinding)
		}

		yellowMatch := MatchResult{
			Topic: topic,
			Rule:  rule,
			HealthFact: domain.HealthFact{
				Category:   topic,
				Label:      topic,
				Evidence:   domain.EvidenceDetail{Loc: "para_2", Text: "abnormal", Date: "unknown"},
				Diagnosed:  &longMed,
				Confidence: 0.95,
			},
			PolicyFacts: []domain.PolicyFact{
				{Type: firstPolicyType(rule), Loc: "para_110", Content: "条款"},
			},
		}
		yellowFinding := scorer.Score([]MatchResult{yellowMatch})
		if len(yellowFinding) != 1 || yellowFinding[0].Level != domain.RiskLevelYellow {
			t.Fatalf("topic=%s expected yellow(date missing), got %#v", topic, yellowFinding)
		}

		greenMatch := MatchResult{
			Topic: topic,
			Rule:  rule,
			HealthFact: domain.HealthFact{
				Category:   topic,
				Label:      topic,
				Evidence:   domain.EvidenceDetail{Loc: "para_3", Text: "mild abnormal", Date: "2026-01-02"},
				Diagnosed:  &notDiagnosed,
				Confidence: 0.95,
			},
			PolicyFacts: []domain.PolicyFact{
				{Type: firstPolicyType(rule), Loc: "para_120", Content: "条款"},
			},
		}
		greenFinding := scorer.Score([]MatchResult{greenMatch})
		if len(greenFinding) != 1 || greenFinding[0].Level != domain.RiskLevelGreen {
			t.Fatalf("topic=%s expected green, got %#v", topic, greenFinding)
		}
	}
}

func firstPolicyType(rule TopicRule) string {
	if len(rule.PolicyTypes) == 0 {
		return "exclusion"
	}
	return rule.PolicyTypes[0]
}
