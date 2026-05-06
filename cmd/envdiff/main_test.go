package main

import (
	"testing"
)

func TestSplitComma_Single(t *testing.T) {
	got := splitComma(".env.dev")
	if len(got) != 1 || got[0] != ".env.dev" {
		t.Fatalf("expected [".env.dev"], got %v", got)
	}
}

func TestSplitComma_Multiple(t *testing.T) {
	got := splitComma(".env.dev,.env.prod,.env.staging")
	want := []string{".env.dev", ".env.prod", ".env.staging"}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i, v := range want {
		if got[i] != v {
			t.Errorf("index %d: expected %q, got %q", i, v, got[i])
		}
	}
}

func TestSplitComma_Empty(t *testing.T) {
	got := splitComma("")
	if len(got) != 1 || got[0] != "" {
		t.Fatalf("expected single empty string, got %v", got)
	}
}

func TestSplitComma_TrailingComma(t *testing.T) {
	got := splitComma(".env.dev,")
	if len(got) != 2 {
		t.Fatalf("expected 2 elements, got %d: %v", len(got), got)
	}
	if got[0] != ".env.dev" {
		t.Errorf("expected .env.dev, got %q", got[0])
	}
	if got[1] != "" {
		t.Errorf("expected empty string, got %q", got[1])
	}
}

func TestLoadCommaSeparated_MissingFile(t *testing.T) {
	_, err := loadCommaSeparated("/nonexistent/path/.env.dev")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadCommaSeparated_EmptyString(t *testing.T) {
	envs, err := loadCommaSeparated("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(envs) != 0 {
		t.Fatalf("expected empty map, got %v", envs)
	}
}
