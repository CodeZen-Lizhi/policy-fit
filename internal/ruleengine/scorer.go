package ruleengine

import (
	"fmt"

	"github.com/zhenglizhi/policy-fit/internal/domain"
)

// Scorer 评分器
type Scorer struct {
	lowConfidenceThreshold float64
}

// NewScorer 创建评分器
func NewScorer() *Scorer {
	return &Scorer{
		lowConfidenceThreshold: 0.7,
	}
}

// Score 红黄绿评分
func (s *Scorer) Score(matches []MatchResult) []domain.RiskFinding {
	findings := make([]domain.RiskFinding, 0, len(matches))
	for _, match := range matches {
		level := s.resolveLevel(match)
		healthEvidence := []domain.Evidence{
			{
				Loc:  match.HealthFact.Evidence.Loc,
				Text: match.HealthFact.Evidence.Text,
			},
		}

		policyEvidence := make([]domain.Evidence, 0, len(match.PolicyFacts))
		for _, pf := range match.PolicyFacts {
			policyEvidence = append(policyEvidence, domain.Evidence{
				Loc:  pf.Loc,
				Text: pf.Content,
			})
		}
		// 证据缺失时提供显式占位，保证结构完整
		if len(policyEvidence) == 0 {
			policyEvidence = append(policyEvidence, domain.Evidence{
				Loc:  "unknown",
				Text: "未命中条款证据",
			})
		}

		questions := match.Rule.Questions
		if len(questions) == 0 {
			questions = []string{
				"该异常是否有明确诊断记录？",
				"体检发生时间是否早于投保时间？",
			}
		}

		actions := match.Rule.Actions
		if len(actions) == 0 {
			actions = []string{
				"补充投保告知材料",
				"联系保险公司核保确认",
			}
		}

		findings = append(findings, domain.RiskFinding{
			TaskID:         0,
			Level:          level,
			Topic:          match.Topic,
			Summary:        buildSummary(match.Topic, level),
			HealthEvidence: healthEvidence,
			PolicyEvidence: policyEvidence,
			Questions:      trimToMax(questions, 5),
			Actions:        trimToMax(actions, 3),
			Confidence:     match.HealthFact.Confidence,
		})
	}
	return findings
}

func (s *Scorer) resolveLevel(match MatchResult) domain.RiskLevel {
	hasPolicyEvidence := len(match.PolicyFacts) > 0
	diagnosed := match.HealthFact.Diagnosed != nil && *match.HealthFact.Diagnosed
	longTermMedication := match.HealthFact.LongTermMedication != nil && *match.HealthFact.LongTermMedication
	dateMissing := match.HealthFact.Evidence.Date == "" || match.HealthFact.Evidence.Date == "unknown"
	lowConfidence := match.HealthFact.Confidence < s.lowConfidenceThreshold

	if (diagnosed || longTermMedication) && hasPolicyEvidence {
		if lowConfidence || dateMissing {
			return domain.RiskLevelYellow
		}
		return domain.RiskLevelRed
	}

	if dateMissing || !hasPolicyEvidence || lowConfidence {
		return domain.RiskLevelYellow
	}

	return domain.RiskLevelGreen
}

func buildSummary(topic string, level domain.RiskLevel) string {
	switch level {
	case domain.RiskLevelRed:
		return fmt.Sprintf("%s 可能触发既往症/免责冲突", topic)
	case domain.RiskLevelYellow:
		return fmt.Sprintf("%s 存在待确认风险信息", topic)
	default:
		return fmt.Sprintf("%s 暂未发现高风险冲突", topic)
	}
}

func trimToMax(values []string, max int) []string {
	if len(values) <= max {
		return values
	}
	return values[:max]
}
