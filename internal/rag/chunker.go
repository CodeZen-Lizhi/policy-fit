package rag

import (
	"fmt"
	"regexp"
	"strings"
)

var clauseHeadingPattern = regexp.MustCompile(`(?m)^\s*(第[一二三四五六七八九十百\d]+条|[0-9]+\.|\([0-9]+\))`)

// ChunkByClause 按条款结构切分文本
func ChunkByClause(documentID string, text string, maxChars int) []Chunk {
	if strings.TrimSpace(text) == "" {
		return nil
	}
	if maxChars <= 0 {
		maxChars = 600
	}

	segments := clauseHeadingPattern.Split(text, -1)
	headings := clauseHeadingPattern.FindAllString(text, -1)

	chunks := make([]Chunk, 0, len(segments))
	idx := 0
	for i, seg := range segments {
		segment := strings.TrimSpace(seg)
		if segment == "" {
			continue
		}
		heading := "section"
		if i > 0 && i-1 < len(headings) {
			heading = strings.TrimSpace(headings[i-1])
		}

		parts := splitByLength(segment, maxChars)
		for _, part := range parts {
			idx++
			chunks = append(chunks, Chunk{
				ID:         fmt.Sprintf("%s-%d", documentID, idx),
				DocumentID: documentID,
				Section:    heading,
				Text:       part,
			})
		}
	}
	if len(chunks) == 0 {
		chunks = append(chunks, Chunk{
			ID:         fmt.Sprintf("%s-1", documentID),
			DocumentID: documentID,
			Section:    "full_text",
			Text:       strings.TrimSpace(text),
		})
	}
	return chunks
}

func splitByLength(text string, maxChars int) []string {
	if len(text) <= maxChars {
		return []string{text}
	}
	parts := make([]string, 0, (len(text)/maxChars)+1)
	for start := 0; start < len(text); start += maxChars {
		end := start + maxChars
		if end > len(text) {
			end = len(text)
		}
		parts = append(parts, strings.TrimSpace(text[start:end]))
	}
	return parts
}
