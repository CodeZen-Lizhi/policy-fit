package rag

import "testing"

func TestChunkByClause(t *testing.T) {
	text := "第1条 既往症定义。\n内容A。\n\n第2条 等待期。\n内容B。"
	chunks := ChunkByClause("doc-1", text, 20)
	if len(chunks) < 2 {
		t.Fatalf("expected at least 2 chunks, got %d", len(chunks))
	}
	if chunks[0].DocumentID != "doc-1" {
		t.Fatalf("unexpected document id: %s", chunks[0].DocumentID)
	}
}
