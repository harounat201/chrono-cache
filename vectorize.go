// create a struct QVectorizer to turn queries into embeddings.
// preprocess inferences (strip JSON to surface query titles)

 package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// QVectorizer handles turning raw queries into vector embeddings
type QVectorizer struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewQVectorizer creates a new QVectorizer instance
func NewQVectorizer(apiKey string) *QVectorizer {
	return &QVectorizer{
		apiKey:     apiKey,
		model:      "text-embedding-3-small",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Vectorize takes a raw query, preprocesses it, and returns its embedding
func (v *QVectorizer) Vectorize(raw string) ([]float32, error) {
	query := v.preprocess(raw)
	if query == "" {
		return nil, fmt.Errorf("empty query after preprocessing")
	}
	return v.embed(query)
}

// preprocess strips JSON structure to surface the semantic query content
func (v *QVectorizer) preprocess(raw string) string {
	raw = strings.TrimSpace(raw)

	// attempt to parse as JSON and extract the query field
	var payload map[string]any
	if err := json.Unmarshal([]byte(raw), &payload); err == nil {
		// priority order: "query" → "content" → "prompt" → "text"
		for _, key := range []string{"query", "content", "prompt", "text"} {
			if val, ok := payload[key]; ok {
				if s, ok := val.(string); ok && strings.TrimSpace(s) != "" {
					return strings.TrimSpace(s)
				}
			}
		}

		// handle Anthropic messages format: { messages: [{ role, content }] }
		if msgs, ok := payload["messages"].([]any); ok && len(msgs) > 0 {
			// surface the last user message
			for i := len(msgs) - 1; i >= 0; i-- {
				if msg, ok := msgs[i].(map[string]any); ok {
					if msg["role"] == "user" {
						if content, ok := msg["content"].(string); ok {
							return strings.TrimSpace(content)
						}
					}
				}
			}
		}
	}

	// not JSON — treat as a raw string query
	return raw
}

// embed calls the OpenAI embeddings API and returns the vector
func (v *QVectorizer) embed(query string) ([]float32, error) {
	body, err := json.Marshal(map[string]any{
		"input": query,
		"model": v.model,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal embed request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create embed request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+v.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("embed request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("embed API error %d: %s", resp.StatusCode, string(b))
	}

	var result struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode embed response: %w", err)
	}
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return result.Data[0].Embedding, nil
}