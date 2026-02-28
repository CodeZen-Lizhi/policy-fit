package llm

import (
	"embed"
	"fmt"
)

//go:embed prompts/*.tmpl
var promptFS embed.FS

// PromptTemplates 提示词模板
var PromptTemplates = map[string]string{}

func init() {
	loadPrompt("health_facts.tmpl")
	loadPrompt("policy_facts.tmpl")
}

func loadPrompt(name string) {
	raw, err := promptFS.ReadFile("prompts/" + name)
	if err != nil {
		panic(fmt.Sprintf("load prompt %s failed: %v", name, err))
	}
	PromptTemplates[name] = string(raw)
}
