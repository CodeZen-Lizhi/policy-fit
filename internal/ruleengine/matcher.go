package ruleengine

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/zhenglizhi/policy-fit/internal/domain"
)

// TopicRule 主题规则
type TopicRule struct {
	HitKeywords []string `yaml:"hit_keywords"`
	PolicyTypes []string `yaml:"policy_types"`
	Questions   []string `yaml:"questions"`
	Actions     []string `yaml:"actions"`
}

// TopicsConfig 主题配置
type TopicsConfig struct {
	Topics map[string]TopicRule `yaml:"topics"`
}

// MatchResult 匹配结果
type MatchResult struct {
	Topic       string
	Rule        TopicRule
	HealthFact  domain.HealthFact
	PolicyFacts []domain.PolicyFact
}

// Matcher 规则匹配器
type Matcher struct {
	config *TopicsConfig
}

// LoadTopicsConfig 加载主题配置
func LoadTopicsConfig(filePath string) (*TopicsConfig, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var cfg TopicsConfig
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}
	if cfg.Topics == nil {
		cfg.Topics = map[string]TopicRule{}
	}
	return &cfg, nil
}

// NewMatcher 创建匹配器
func NewMatcher(cfg *TopicsConfig) *Matcher {
	if cfg == nil {
		cfg = &TopicsConfig{Topics: map[string]TopicRule{}}
	}
	return &Matcher{config: cfg}
}

// Match 执行 topic 与 policy 类型匹配
func (m *Matcher) Match(healthFacts []domain.HealthFact, policyFacts []domain.PolicyFact) []MatchResult {
	results := make([]MatchResult, 0)
	for _, hf := range healthFacts {
		topic, rule := m.matchTopicRule(hf)
		if topic == "" {
			continue
		}
		matchedPolicy := make([]domain.PolicyFact, 0)
		for _, pf := range policyFacts {
			if contains(rule.PolicyTypes, pf.Type) {
				matchedPolicy = append(matchedPolicy, pf)
			}
		}

		results = append(results, MatchResult{
			Topic:       topic,
			Rule:        rule,
			HealthFact:  hf,
			PolicyFacts: matchedPolicy,
		})
	}
	return results
}

func (m *Matcher) matchTopicRule(healthFact domain.HealthFact) (string, TopicRule) {
	if rule, ok := m.config.Topics[healthFact.Category]; ok {
		return healthFact.Category, rule
	}

	content := strings.ToLower(healthFact.Label + " " + healthFact.Evidence.Text)
	for topic, rule := range m.config.Topics {
		for _, keyword := range rule.HitKeywords {
			if strings.Contains(content, strings.ToLower(keyword)) {
				return topic, rule
			}
		}
	}
	return "", TopicRule{}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
