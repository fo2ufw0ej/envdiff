package reporter_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/patcher"
	"github.com/yourorg/envdiff/internal/reporter"
)

func TestPatchReporter_NoResults(t *testing.T) {
	var sb strings.Builder
	r := reporter.NewPatchReporter(&sb)
	if err := r.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "No patch actions") {
		t.Errorf("expected empty message, got: %q", sb.String())
	}
}

func TestPatchReporter_Added(t *testing.T) {
	results := []patcher.Result{
		{Key: "NEW_KEY", OldVal: "", NewVal: "hello", Action: "added"},
	}
	var sb strings.Builder
	r := reporter.NewPatchReporter(&sb)
	if err := r.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "NEW_KEY") {
		t.Errorf("expected NEW_KEY in output, got: %q", out)
	}
	if !strings.Contains(out, "added") {
		t.Errorf("expected action 'added' in output, got: %q", out)
	}
	if !strings.Contains(out, "<absent>") {
		t.Errorf("expected <absent> placeholder for old value, got: %q", out)
	}
}

func TestPatchReporter_Overwritten(t *testing.T) {
	results := []patcher.Result{
		{Key: "DB_HOST", OldVal: "localhost", NewVal: "prod.db", Action: "overwritten"},
		{Key: "PORT", OldVal: "5432", NewVal: "5432", Action: "skipped"},
	}
	var sb strings.Builder
	r := reporter.NewPatchReporter(&sb)
	if err := r.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "overwritten") {
		t.Errorf("expected 'overwritten' in output, got: %q", out)
	}
	if !strings.Contains(out, "skipped") {
		t.Errorf("expected 'skipped' in output, got: %q", out)
	}
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected old value 'localhost' in output, got: %q", out)
	}
}
