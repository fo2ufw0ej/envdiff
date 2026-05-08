package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/baseline"
	"github.com/user/envdiff/internal/comparator"
)

func results() []comparator.Result {
	return []comparator.Result{
		{Key: "APP_ENV", Env: "production", Status: "identical", Value: "prod"},
		{Key: "DB_URL", Env: "staging", Status: "missing", Value: ""},
		{Key: "SECRET", Env: "production", Status: "mismatch", Value: "abc"},
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	original := results()
	if err := baseline.Save(path, original); err != nil {
		t.Fatalf("Save: %v", err)
	}

	b, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(b.Entries) != len(original) {
		t.Fatalf("expected %d entries, got %d", len(original), len(b.Entries))
	}
	for i, e := range b.Entries {
		if e.Key != original[i].Key || e.Status != original[i].Status {
			t.Errorf("entry %d mismatch: got %+v", i, e)
		}
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := baseline.Load("/nonexistent/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestDiff_NewEntry(t *testing.T) {
	base := &baseline.Baseline{Entries: results()}
	current := append(results(), comparator.Result{Key: "NEW_KEY", Env: "production", Status: "missing"})

	delta := baseline.Diff(base, current)
	if len(delta) != 1 || delta[0].Key != "NEW_KEY" {
		t.Errorf("expected 1 new entry NEW_KEY, got %+v", delta)
	}
}

func TestDiff_ChangedStatus(t *testing.T) {
	base := &baseline.Baseline{Entries: results()}
	current := results()
	current[1].Status = "mismatch" // was "missing"

	delta := baseline.Diff(base, current)
	if len(delta) != 1 || delta[0].Key != "DB_URL" {
		t.Errorf("expected 1 changed entry DB_URL, got %+v", delta)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	base := &baseline.Baseline{Entries: results()}
	delta := baseline.Diff(base, results())
	if len(delta) != 0 {
		t.Errorf("expected no delta, got %+v", delta)
	}
}

func TestSave_InvalidPath(t *testing.T) {
	err := baseline.Save("/no/such/dir/baseline.json", results())
	if err == nil {
		t.Fatal("expected error writing to invalid path")
	}
	_ = os.Remove("/no/such/dir/baseline.json")
}
