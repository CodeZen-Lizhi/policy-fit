package i18n

import (
	"embed"
	"encoding/json"
	"strings"
)

//go:embed messages/*.json
var messageFS embed.FS

var catalogs = map[string]map[string]string{}

func init() {
	catalogs["zh-CN"] = mustLoadCatalog("messages/zh-CN.json")
	catalogs["en-US"] = mustLoadCatalog("messages/en-US.json")
}

func mustLoadCatalog(path string) map[string]string {
	raw, err := messageFS.ReadFile(path)
	if err != nil {
		panic(err)
	}
	out := map[string]string{}
	if err := json.Unmarshal(raw, &out); err != nil {
		panic(err)
	}
	return out
}

// ResolveLanguage 解析请求语言
func ResolveLanguage(acceptLanguage string) string {
	lang := strings.ToLower(strings.TrimSpace(acceptLanguage))
	if strings.HasPrefix(lang, "en") {
		return "en-US"
	}
	return "zh-CN"
}

// TranslateErrorMessage 根据错误码翻译消息
func TranslateErrorMessage(code, acceptLanguage, fallback string) string {
	lang := ResolveLanguage(acceptLanguage)
	if msg, ok := catalogs[lang][code]; ok {
		return msg
	}
	if fallback != "" {
		return fallback
	}
	if msg, ok := catalogs["zh-CN"][code]; ok {
		return msg
	}
	return "Unknown error"
}
