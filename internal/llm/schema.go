package llm

import (
	"fmt"
	"strings"

	"github.com/zhenglizhi/policy-fit/internal/domain"
)

// NormalizeAndValidateHealthFacts 规范化并校验 HealthFacts
func NormalizeAndValidateHealthFacts(facts []domain.HealthFact) error {
	for i := range facts {
		f := &facts[i]
		if strings.TrimSpace(f.Category) == "" {
			return fmt.Errorf("%w: facts[%d].category is required", ErrSchemaValidation, i)
		}
		if strings.TrimSpace(f.Label) == "" {
			f.Label = f.Category
		}
		if strings.TrimSpace(f.Evidence.Text) == "" {
			return fmt.Errorf("%w: facts[%d].evidence.text is required", ErrSchemaValidation, i)
		}
		if strings.TrimSpace(f.Evidence.Loc) == "" {
			return fmt.Errorf("%w: facts[%d].evidence.loc is required", ErrSchemaValidation, i)
		}
		if strings.TrimSpace(f.Evidence.Date) == "" {
			f.Evidence.Date = "unknown"
		}
		if strings.TrimSpace(f.Evidence.Source) == "" {
			f.Evidence.Source = "report"
		}
		if f.Confidence <= 0 {
			f.Confidence = 0.5
		}
		if f.Confidence > 1 {
			return fmt.Errorf("%w: facts[%d].confidence > 1", ErrSchemaValidation, i)
		}
	}
	return nil
}

// NormalizeAndValidatePolicyFacts 规范化并校验 PolicyFacts
func NormalizeAndValidatePolicyFacts(sections []domain.PolicyFact) error {
	for i := range sections {
		s := &sections[i]
		if strings.TrimSpace(s.Type) == "" {
			return fmt.Errorf("%w: sections[%d].type is required", ErrSchemaValidation, i)
		}
		if strings.TrimSpace(s.Title) == "" {
			s.Title = s.Type
		}
		if strings.TrimSpace(s.Content) == "" {
			return fmt.Errorf("%w: sections[%d].content is required", ErrSchemaValidation, i)
		}
		if strings.TrimSpace(s.Loc) == "" {
			return fmt.Errorf("%w: sections[%d].loc is required", ErrSchemaValidation, i)
		}
		if s.Confidence <= 0 {
			s.Confidence = 0.5
		}
		if s.Confidence > 1 {
			return fmt.Errorf("%w: sections[%d].confidence > 1", ErrSchemaValidation, i)
		}
	}
	return nil
}
