package archiver_test

import (
	"testing"
	"time"

	"github.com/user/envdiff/internal/archiver"
	"github.com/user/envdiff/internal/comparator"
)

func sampleResults() []comparator.Result {
	return []comparator.Result{
		{
			Key:    "DB_HOST",
			Status: comparator.StatusMissing,
			Values: map[string]string{"production": "db.prod"},
		},
		{
			Key:    "APP_ENV",
			Status: comparator.StatusIdentical,
			Values: map[string]string{"production": "prod", "staging": "prod"},
		},
	}
}

func TestSaveAndList(t *testing.T) {
	dir := t.TempDir()
	a, err := archiver.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = a.Save(sampleResults(), "first run")
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	entries, err := a.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Label != "first run" {
		t.Errorf("label mismatch: got %q", entries[0].Label)
	}
	if len(entries[0].Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(entries[0].Results))
	}
}

func TestList_SortedOldestFirst(t *testing.T) {
	dir := t.TempDir()
	a, err := archiver.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	for _, label := range []string{"alpha", "beta", "gamma"} {
		time.Sleep(2 * time.Millisecond)
		if _, err := a.Save(sampleResults(), label); err != nil {
			t.Fatalf("Save %s: %v", label, err)
		}
	}

	entries, err := a.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Label != "alpha" || entries[2].Label != "gamma" {
		t.Errorf("wrong order: %v", []string{entries[0].Label, entries[1].Label, entries[2].Label})
	}
}

func TestLatest_Empty(t *testing.T) {
	dir := t.TempDir()
	a, err := archiver.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_, err = a.Latest()
	if err == nil {
		t.Error("expected error for empty archive")
	}
}

func TestLatest_ReturnsNewest(t *testing.T) {
	dir := t.TempDir()
	a, err := archiver.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for _, label := range []string{"old", "new"} {
		time.Sleep(2 * time.Millisecond)
		if _, err := a.Save(sampleResults(), label); err != nil {
			t.Fatalf("Save: %v", err)
		}
	}
	e, err := a.Latest()
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if e.Label != "new" {
		t.Errorf("expected latest label 'new', got %q", e.Label)
	}
}
