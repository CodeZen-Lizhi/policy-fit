package rag

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

// Retriever 召回+重排流程
type Retriever struct {
	store    VectorStore
	embedder EmbeddingProvider
}

// NewRetriever 创建召回器
func NewRetriever(store VectorStore, embedder EmbeddingProvider) *Retriever {
	return &Retriever{store: store, embedder: embedder}
}

// IndexPolicyText 切片并写入向量库
func (r *Retriever) IndexPolicyText(ctx context.Context, documentID string, text string) error {
	chunks := ChunkByClause(documentID, text, 600)
	if len(chunks) == 0 {
		return nil
	}
	texts := make([]string, 0, len(chunks))
	for _, c := range chunks {
		texts = append(texts, c.Text)
	}
	vectors, err := r.embedder.Embed(ctx, texts)
	if err != nil {
		return fmt.Errorf("embed chunks: %w", err)
	}
	if err := r.store.Upsert(ctx, chunks, vectors); err != nil {
		return fmt.Errorf("upsert chunks: %w", err)
	}
	return nil
}

// Retrieve 召回并重排
func (r *Retriever) Retrieve(ctx context.Context, query string, topK int) ([]SearchResult, error) {
	if topK <= 0 {
		topK = 5
	}
	vectors, err := r.embedder.Embed(ctx, []string{query})
	if err != nil {
		return nil, err
	}
	if len(vectors) == 0 {
		return nil, nil
	}
	results, err := r.store.Search(ctx, vectors[0], topK*2)
	if err != nil {
		return nil, err
	}
	results = rerankByTokenOverlap(query, results)
	if len(results) > topK {
		results = results[:topK]
	}
	return results, nil
}

func rerankByTokenOverlap(query string, results []SearchResult) []SearchResult {
	queryTokens := strings.Fields(strings.ToLower(query))
	if len(queryTokens) == 0 {
		return results
	}

	boosted := make([]SearchResult, 0, len(results))
	for _, item := range results {
		text := strings.ToLower(item.Chunk.Text)
		hits := 0
		for _, token := range queryTokens {
			if strings.Contains(text, token) {
				hits++
			}
		}
		if hits > 0 {
			item.Score += float64(hits) * 0.05
		}
		item.Source = "retriever"
		boosted = append(boosted, item)
	}
	sort.Slice(boosted, func(i, j int) bool {
		return boosted[i].Score > boosted[j].Score
	})
	return boosted
}
