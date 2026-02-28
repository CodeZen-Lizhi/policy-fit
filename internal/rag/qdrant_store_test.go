package rag

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestQdrantStoreSearch(t *testing.T) {
	store := NewQdrantStore("http://qdrant", "policy_chunks", 0)
	store.client = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path == "/collections/policy_chunks/points/search" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"result":[{"id":"chunk-1","score":0.91,"payload":{"document_id":"p1","section":"第1条","text":"既往症定义"}}]}`)),
				Header:     make(http.Header),
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"result":[]}`)),
			Header:     make(http.Header),
		}, nil
	})}

	results, err := store.Search(context.Background(), []float64{0.1, 0.2}, 3)
	if err != nil {
		t.Fatalf("Search error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Source != "qdrant" {
		t.Fatalf("unexpected source: %s", results[0].Source)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
