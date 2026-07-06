package main

import (
	"testing"
)

func TestDecode_PooledDecoderStateLeakage(t *testing.T) {
	type Target struct {
		ID   string `json:"id"`
		Name string `json:"name,omitempty"`
	}

	// Payload 1 has both fields
	payload1 := []byte(`{"id": "123", "name": "Alice"}`)
	// Payload 2 only has ID
	payload2 := []byte(`{"id": "456"}`)

	// Run decodes sequentially, simulating decoder reuse from the pool
	var doc1 Target
	if err := Unmarshal(payload1, &doc1); err != nil {
		t.Fatalf("failed to unmarshal payload 1: %v", err)
	}
	if doc1.Name != "Alice" {
		t.Errorf("expected doc1.Name to be 'Alice', got '%s'", doc1.Name)
	}

	var doc2 Target
	if err := Unmarshal(payload2, &doc2); err != nil {
		t.Fatalf("failed to unmarshal payload 2: %v", err)
	}

	// If state leaks, doc2.Name might incorrectly retain "Alice"
	if doc2.Name != "" {
		t.Errorf("expected doc2.Name to be empty, but got stale value: '%s'", doc2.Name)
	}
}
