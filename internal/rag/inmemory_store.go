package rag

import (
	"context"
	"errors"
	"sort"
)

type vectorItem struct {
	chunk  Chunk
	vector []float64
}

// InMemoryStore 内存向量存储（测试/开发）
type InMemoryStore struct {
	items []vectorItem
}

// NewInMemoryStore 创建内存向量存储
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{items: make([]vectorItem, 0)}
}

// Upsert 写入向量
func (s *InMemoryStore) Upsert(_ context.Context, chunks []Chunk, vectors [][]float64) error {
	if len(chunks) != len(vectors) {
		return errors.New("chunks and vectors length mismatch")
	}
	for i := range chunks {
		s.items = append(s.items, vectorItem{chunk: chunks[i], vector: vectors[i]})
	}
	return nil
}

// Search 向量检索
func (s *InMemoryStore) Search(_ context.Context, vector []float64, topK int) ([]SearchResult, error) {
	if topK <= 0 {
		topK = 5
	}
	results := make([]SearchResult, 0, len(s.items))
	for _, item := range s.items {
		score := cosineSimilarity(vector, item.vector)
		results = append(results, SearchResult{
			Chunk:  item.chunk,
			Score:  score,
			Source: "inmemory",
		})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	if len(results) > topK {
		return results[:topK], nil
	}
	return results, nil
}
