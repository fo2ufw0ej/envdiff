package flattener_test

import (
	"fmt"
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/flattener"
)

func makeResult(key string, status comparator.Status, values map[string]string, missingIn []string) comparator.Result {
	return comparator.Result{
		Key:       key,
		Status:    status,
		Values:    values,
		MissingIn: missingIn,
	}
}

func TestFlatten_Empty(t *testing.T) {
	entries := flattener.Flatten(nil)
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestFlatten_IdenticalKey(t *testing.T) {
	results := []comparator.Result{
		makeResult("PORT", comparator.StatusIdentical, map[string]string{"dev": "8080", "prod": "8080"}, nil),
	}
	entries := flattener.Flatten(results)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Key != "PORT" {
			t.Errorf("unexpected key %q", e.Key)
		}
		if e.Value != "8080" {
			t.Errorf("unexpected value %q", e.Value)
		}
	}
}

func TestFlatten_MissingKey(t *testing.T) {
	results := []comparator.Result{
		makeResult("SECRET", comparator.StatusMissing, map[string]string{"dev": "abc"}, []string{"prod"}),
	}
	entries := flattener.Flatten(results)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	var missing *flattener.Entry
	for i := range entries {
		if entries[i].Env == "prod" {
			missing = &entries[i]
		}
	}
	if missing == nil {
		t.Fatal("expected a prod entry")
	}
	if missing.Status != comparator.StatusMissing {
		t.Errorf("expected Missing status, got %v", missing.Status)
	}
	if missing.Value != "" {
		t.Errorf("expected empty value, got %q", missing.Value)
	}
}

func TestFlatten_SortedOutput(t *testing.T) {
	results := []comparator.Result{
		makeResult("Z_KEY", comparator.StatusIdentical, map[string]string{"dev": "1"}, nil),
		makeResult("A_KEY", comparator.StatusIdentical, map[string]string{"dev": "2"}, nil),
	}
	entries := flattener.Flatten(results)
	if entries[0].Key != "A_KEY" {
		t.Errorf("expected A_KEY first, got %q", entries[0].Key)
	}
	if entries[1].Key != "Z_KEY" {
		t.Errorf("expected Z_KEY second, got %q", entries[1].Key)
	}
}

func TestIndex_Lookup(t *testing.T) {
	results := []comparator.Result{
		makeResult("DB_URL", comparator.StatusMismatch, map[string]string{"dev": "localhost", "prod": "db.example.com"}, nil),
	}
	entries := flattener.Flatten(results)
	idx := flattener.Index(entries)

	key := fmt.Sprintf("%s@%s", "DB_URL", "prod")
	e, ok := idx[key]
	if !ok {
		t.Fatalf("expected key %q in index", key)
	}
	if e.Value != "db.example.com" {
		t.Errorf("expected db.example.com, got %q", e.Value)
	}
}
