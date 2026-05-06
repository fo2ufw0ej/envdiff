package filter_test

import (
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/filter"
)

func entries() []comparator.DiffEntry {
	return []comparator.DiffEntry{
		{Key: "DB_HOST", Values: map[string]string{"dev": "localhost", "prod": "db.prod"}},
		{Key: "DB_PORT", Values: map[string]string{"dev": "5432", "prod": "5432"}},
		{Key: "API_KEY", Values: map[string]string{"dev": "abc", "prod": ""}},
		{Key: "DEBUG", Values: map[string]string{"dev": "true", "prod": ""}},
	}
}

func TestApply_NoFilter(t *testing.T) {
	result := filter.Apply(entries(), filter.Options{})
	if len(result) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(result))
	}
}

func TestApply_OnlyMissing(t *testing.T) {
	result := filter.Apply(entries(), filter.Options{OnlyMissing: true})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Key != "API_KEY" && e.Key != "DEBUG" {
			t.Errorf("unexpected key %q in missing-only results", e.Key)
		}
	}
}

func TestApply_OnlyMismatched(t *testing.T) {
	result := filter.Apply(entries(), filter.Options{OnlyMismatched: true})
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
}

func TestApply_KeyPrefix(t *testing.T) {
	result := filter.Apply(entries(), filter.Options{KeyPrefix: "DB_"})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestApply_KeyPrefix_CaseInsensitive(t *testing.T) {
	result := filter.Apply(entries(), filter.Options{KeyPrefix: "db_"})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries with lowercase prefix, got %d", len(result))
	}
}

func TestApply_CombinedFilters(t *testing.T) {
	result := filter.Apply(entries(), filter.Options{OnlyMissing: true, KeyPrefix: "API"})
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Key != "API_KEY" {
		t.Errorf("expected API_KEY, got %q", result[0].Key)
	}
}

func TestApply_EmptyEntries(t *testing.T) {
	result := filter.Apply(nil, filter.Options{OnlyMissing: true})
	if len(result) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result))
	}
}
