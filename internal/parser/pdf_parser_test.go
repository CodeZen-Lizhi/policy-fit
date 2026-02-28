package parser

import (
	"context"
	"errors"
	"testing"
)

func TestPDFParserParseSuccess(t *testing.T) {
	p := NewPDFParser("pdftotext")
	p.runner = func(context.Context, string, string) (string, error) {
		return "第一段内容\n\n第二段内容\f第三页段落", nil
	}

	result, err := p.Parse(context.Background(), "dummy.pdf")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(result.Paragraphs) != 3 {
		t.Fatalf("unexpected paragraph length: %d", len(result.Paragraphs))
	}
	if result.Paragraphs[0].Loc != "para_1" || result.Paragraphs[2].Loc != "para_3" {
		t.Fatalf("unexpected paragraph location: %#v", result.Paragraphs)
	}
	if result.Paragraphs[2].Page != 2 {
		t.Fatalf("expected third paragraph on page 2, got %d", result.Paragraphs[2].Page)
	}
}

func TestPDFParserParseEmptyText(t *testing.T) {
	p := NewPDFParser("pdftotext")
	p.runner = func(context.Context, string, string) (string, error) {
		return "   \n\n", nil
	}

	_, err := p.Parse(context.Background(), "dummy.pdf")
	if !errors.Is(err, ErrEmptyText) {
		t.Fatalf("expected ErrEmptyText, got %v", err)
	}
}

func TestPDFParserParseCommandError(t *testing.T) {
	p := NewPDFParser("pdftotext")
	p.runner = func(context.Context, string, string) (string, error) {
		return "", errors.New("command failed")
	}

	_, err := p.Parse(context.Background(), "dummy.pdf")
	if !errors.Is(err, ErrUnreadablePDF) {
		t.Fatalf("expected ErrUnreadablePDF, got %v", err)
	}
}
