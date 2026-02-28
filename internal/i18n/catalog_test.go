package i18n

import "testing"

func TestTranslateErrorMessage(t *testing.T) {
	msg := TranslateErrorMessage("PFIT-1001", "en-US", "")
	if msg != "Invalid request parameters" {
		t.Fatalf("unexpected en message: %s", msg)
	}

	msg = TranslateErrorMessage("PFIT-1001", "zh-CN", "")
	if msg == "" {
		t.Fatalf("expected zh message")
	}

	msg = TranslateErrorMessage("UNKNOWN", "en-US", "fallback")
	if msg != "fallback" {
		t.Fatalf("expected fallback, got %s", msg)
	}
}
