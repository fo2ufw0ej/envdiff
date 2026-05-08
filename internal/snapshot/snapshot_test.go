package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/snapshot"
)

func results() []comparator.Result {
	return []comparator.Result{
		{Key: "DB_HOST", Env: "production", Status: comparator.StatusIdentical, Value: "db.prod"},
		{Key: "API_KEY", Env: "staging", Status: comparator.StatusMissing, Value: ""},
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	if err := snapshot.Save(path, "test-label", results()); err != nil {
		t.Fatalf("Save: %v", err)
	}

	s, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if s.Label != "test-label" {
		t.Errorf("label: got %q, want %q", s.Label, "test-label")
	}
	if len(s.Results) != 2 {
		t.Errorf("results len: got %d, want 2", len(s.Results))
	}
	if s.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestCompare_NoChanges(t *testing.T) {
	old := &snapshot.Snapshot{CreatedAt: time.Now(), Label: "a", Results: results()}
	new := &snapshot.Snapshot{CreatedAt: time.Now(), Label: "b", Results: results()}

	deltas := snapshot.Compare(old, new)
	if len(deltas) != 0 {
		t.Errorf("expected 0 deltas, got %d", len(deltas))
	}
}

func TestCompare_StatusChanged(t *testing.T) {
	old := &snapshot.Snapshot{Results: []comparator.Result{
		{Key: "DB_HOST", Env: "production", Status: comparator.StatusMissing},
	}}
	new := &snapshot.Snapshot{Results: []comparator.Result{
		{Key: "DB_HOST", Env: "production", Status: comparator.StatusIdentical, Value: "db.prod"},
	}}

	deltas := snapshot.Compare(old, new)
	if len(deltas) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(deltas))
	}
	if deltas[0].Added {
		t.Error("expected changed, not added")
	}
}

func TestCompare_AddedEntry(t *testing.T) {
	old := &snapshot.Snapshot{Results: []comparator.Result{}}
	new := &snapshot.Snapshot{Results: []comparator.Result{
		{Key: "NEW_KEY", Env: "staging", Status: comparator.StatusMissing},
	}}

	deltas := snapshot.Compare(old, new)
	if len(deltas) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(deltas))
	}
	if !deltas[0].Added {
		t.Error("expected Added=true")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	err := snapshot.Save(filepath.Join(os.TempDir(), "nonexistent", "dir", "snap.json"), "x", nil)
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}
