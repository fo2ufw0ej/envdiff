package reporter_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/reporter"
	"github.com/user/envdiff/internal/snapshot"
)

func TestSnapshotReporter_NoDeltas(t *testing.T) {
	var buf strings.Builder
	r := reporter.NewSnapshotReporter(&buf)

	if err := r.Report("v1", "v2", nil); err != nil {
		t.Fatalf("Report: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "No changes") {
		t.Errorf("expected 'No changes', got: %q", out)
	}
}

func TestSnapshotReporter_StatusChanged(t *testing.T) {
	var buf strings.Builder
	r := reporter.NewSnapshotReporter(&buf)

	deltas := []snapshot.Delta{
		{
			Old: comparator.Result{Key: "DB_HOST", Env: "prod", Status: comparator.StatusMissing},
			New: comparator.Result{Key: "DB_HOST", Env: "prod", Status: comparator.StatusIdentical, Value: "db"},
		},
	}

	if err := r.Report("v1", "v2", deltas); err != nil {
		t.Fatalf("Report: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %q", out)
	}
	if !strings.Contains(out, "missing") {
		t.Errorf("expected old status 'missing' in output, got: %q", out)
	}
	if !strings.Contains(out, "identical") {
		t.Errorf("expected new status 'identical' in output, got: %q", out)
	}
}

func TestSnapshotReporter_AddedEntry(t *testing.T) {
	var buf strings.Builder
	r := reporter.NewSnapshotReporter(&buf)

	deltas := []snapshot.Delta{
		{
			New:   comparator.Result{Key: "NEW_KEY", Env: "staging", Status: comparator.StatusMissing},
			Added: true,
		},
	}

	if err := r.Report("v1", "v2", deltas); err != nil {
		t.Fatalf("Report: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "NEW_KEY") {
		t.Errorf("expected NEW_KEY in output, got: %q", out)
	}
	if !strings.Contains(out, "(new)") {
		t.Errorf("expected '(new)' marker in output, got: %q", out)
	}
}
