package rag

import (
	"context"
	"hash/fnv"
	"math"
	"strings"
)

// HashEmbeddingProvider 开发环境用向量器（稳定可测）
type HashEmbeddingProvider struct {
	dim int
}

// NewHashEmbeddingProvider 创建 hash 向量器
func NewHashEmbeddingProvider(dim int) *HashEmbeddingProvider {
	if dim <= 0 {
		dim = 32
	}
	return &HashEmbeddingProvider{dim: dim}
}

// Embed 计算文本向量
func (p *HashEmbeddingProvider) Embed(_ context.Context, texts []string) ([][]float64, error) {
	vectors := make([][]float64, 0, len(texts))
	for _, text := range texts {
		vec := make([]float64, p.dim)
		tokens := tokenize(text)
		for _, token := range tokens {
			h := fnv.New64a()
			_, _ = h.Write([]byte(token))
			idx := int(h.Sum64() % uint64(p.dim))
			vec[idx] += 1
		}
		normalize(vec)
		vectors = append(vectors, vec)
	}
	return vectors, nil
}

func tokenize(text string) []string {
	cleaned := strings.ToLower(strings.TrimSpace(text))
	if cleaned == "" {
		return []string{"_empty_"}
	}
	tokens := strings.Fields(cleaned)
	if len(tokens) == 0 {
		return []string{cleaned}
	}
	return tokens
}

func normalize(vec []float64) {
	sum := 0.0
	for _, v := range vec {
		sum += v * v
	}
	if sum <= 0 {
		return
	}
	norm := math.Sqrt(sum)
	for i := range vec {
		vec[i] = vec[i] / norm
	}
}
