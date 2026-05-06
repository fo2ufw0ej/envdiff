package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/comparator"
)

func TestMarkdownReporter_NoDifferences(t *testing.T) {
	var buf bytes.Buffer
	r := NewMarkdownReporter(&buf)

	if err := r.Report(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "No differences") {
		t.Errorf("expected no-differences message, got: %q", got)
	}
}

func TestMarkdownReporter_MissingKey(t *testing.T) {
	var buf bytes.Buffer
	r := NewMarkdownReporter(&buf)

	results := []comparator.Result{
		{
			Key:      "DB_HOST",
			Missing:  []string{"staging"},
			Values:   map[string]string{"production": "db.prod.internal"},
		},
	}

	if err := r.Report(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "DB_HOST") {
		t.Errorf("expected key DB_HOST in output, got:\n%s", got)
	}
	if !strings.Contains(got, "_(missing)_") {
		t.Errorf("expected missing marker in output, got:\n%s", got)
	}
	if !strings.Contains(got, "db.prod.internal") {
		t.Errorf("expected value in output, got:\n%s", got)
	}
}

func TestMarkdownReporter_MismatchedKey(t *testing.T) {
	var buf bytes.Buffer
	r := NewMarkdownReporter(&buf)

	results := []comparator.Result{
		{
			Key:      "LOG_LEVEL",
			Missing:  nil,
			Values:   map[string]string{"dev": "debug", "prod": "error"},
		},
	}

	if err := r.Report(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "LOG_LEVEL") {
		t.Errorf("expected key LOG_LEVEL in output, got:\n%s", got)
	}
	if !strings.Contains(got, "debug") || !strings.Contains(got, "error") {
		t.Errorf("expected both values in output, got:\n%s", got)
	}
	// Verify header contains both env names
	if !strings.Contains(got, "dev") || !strings.Contains(got, "prod") {
		t.Errorf("expected env names in header, got:\n%s", got)
	}
}

func TestMarkdownReporter_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	r := NewMarkdownReporter(&buf)

	results := []comparator.Result{
		{Key: "Z_KEY", Values: map[string]string{"env1": "z"}},
		{Key: "A_KEY", Values: map[string]string{"env1": "a"}},
	}

	if err := r.Report(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	aIdx := strings.Index(got, "A_KEY")
	zIdx := strings.Index(got, "Z_KEY")
	if aIdx == -1 || zIdx == -1 {
		t.Fatalf("expected both keys in output, got:\n%s", got)
	}
	if aIdx > zIdx {
		t.Errorf("expected A_KEY before Z_KEY in output, got:\n%s", got)
	}
}
