package rag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// QdrantStore Qdrant 向量库存储实现
type QdrantStore struct {
	baseURL    string
	collection string
	client     *http.Client
}

// NewQdrantStore 创建 Qdrant 存储
func NewQdrantStore(baseURL, collection string, timeout time.Duration) *QdrantStore {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &QdrantStore{
		baseURL:    strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		collection: strings.TrimSpace(collection),
		client:     &http.Client{Timeout: timeout},
	}
}

// Upsert 写入 Qdrant
func (s *QdrantStore) Upsert(ctx context.Context, chunks []Chunk, vectors [][]float64) error {
	if len(chunks) != len(vectors) {
		return fmt.Errorf("chunks and vectors length mismatch")
	}
	points := make([]map[string]interface{}, 0, len(chunks))
	for i := range chunks {
		points = append(points, map[string]interface{}{
			"id":     chunks[i].ID,
			"vector": vectors[i],
			"payload": map[string]interface{}{
				"document_id": chunks[i].DocumentID,
				"section":     chunks[i].Section,
				"text":        chunks[i].Text,
			},
		})
	}
	body := map[string]interface{}{"points": points}
	return s.postJSON(ctx, fmt.Sprintf("/collections/%s/points", s.collection), body, nil)
}

// Search 检索 Qdrant
func (s *QdrantStore) Search(ctx context.Context, vector []float64, topK int) ([]SearchResult, error) {
	if topK <= 0 {
		topK = 5
	}
	var response struct {
		Result []struct {
			ID      interface{}            `json:"id"`
			Score   float64                `json:"score"`
			Payload map[string]interface{} `json:"payload"`
		} `json:"result"`
	}

	err := s.postJSON(ctx, fmt.Sprintf("/collections/%s/points/search", s.collection), map[string]interface{}{
		"vector": vector,
		"limit":  topK,
	}, &response)
	if err != nil {
		return nil, err
	}

	results := make([]SearchResult, 0, len(response.Result))
	for _, item := range response.Result {
		chunk := Chunk{
			ID:         fmt.Sprintf("%v", item.ID),
			DocumentID: asString(item.Payload["document_id"]),
			Section:    asString(item.Payload["section"]),
			Text:       asString(item.Payload["text"]),
		}
		results = append(results, SearchResult{
			Chunk:  chunk,
			Score:  item.Score,
			Source: "qdrant",
		})
	}
	return results, nil
}

func (s *QdrantStore) postJSON(ctx context.Context, path string, reqBody interface{}, out interface{}) error {
	if s.baseURL == "" || s.collection == "" {
		return fmt.Errorf("qdrant store not configured")
	}
	raw, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+path, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("qdrant status=%d body=%s", resp.StatusCode, string(body))
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return err
	}
	return nil
}

func asString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}
