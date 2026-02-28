package rag

import "context"

// Chunk 文本切片
type Chunk struct {
	ID         string            `json:"id"`
	DocumentID string            `json:"document_id"`
	Section    string            `json:"section"`
	Text       string            `json:"text"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// SearchResult 召回结果
type SearchResult struct {
	Chunk  Chunk   `json:"chunk"`
	Score  float64 `json:"score"`
	Source string  `json:"source"`
}

// EmbeddingProvider 向量化抽象
type EmbeddingProvider interface {
	Embed(ctx context.Context, texts []string) ([][]float64, error)
}

// VectorStore 向量存储抽象
type VectorStore interface {
	Upsert(ctx context.Context, chunks []Chunk, vectors [][]float64) error
	Search(ctx context.Context, vector []float64, topK int) ([]SearchResult, error)
}
