package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/comparator"
)

func TestCSVReporter_NoDifferences(t *testing.T) {
	var buf bytes.Buffer
	r := NewCSVReporter(&buf)

	if err := r.Report([]comparator.Result{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "key,status,env,value") {
		t.Errorf("expected CSV header, got: %q", output)
	}
	// Only the header line should be present
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 line (header only), got %d", len(lines))
	}
}

func TestCSVReporter_MissingKey(t *testing.T) {
	var buf bytes.Buffer
	r := NewCSVReporter(&buf)

	results := []comparator.Result{
		{
			Key:    "DB_HOST",
			Status: comparator.StatusMissing,
			Values: map[string]string{
				"production": "db.prod.example.com",
			},
		},
	}

	if err := r.Report(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %q", output)
	}
	if !strings.Contains(output, "production") {
		t.Errorf("expected env name 'production' in output, got: %q", output)
	}
}

func TestCSVReporter_MismatchedKey(t *testing.T) {
	var buf bytes.Buffer
	r := NewCSVReporter(&buf)

	results := []comparator.Result{
		{
			Key:    "LOG_LEVEL",
			Status: comparator.StatusMismatch,
			Values: map[string]string{
				"staging":    "debug",
				"production": "warn",
			},
		},
	}

	if err := r.Report(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "mismatch") {
		t.Errorf("expected 'mismatch' status in output, got: %q", output)
	}
	if !strings.Contains(output, "LOG_LEVEL") {
		t.Errorf("expected key LOG_LEVEL in output, got: %q", output)
	}
	// Two data rows: one per env
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 3 { // header + 2 envs
		t.Errorf("expected 3 lines, got %d:\n%s", len(lines), output)
	}
}

func TestCSVReporter_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	r := NewCSVReporter(&buf)

	results := []comparator.Result{
		{
			Key:    "ZEBRA",
			Status: comparator.StatusMismatch,
			Values: map[string]string{"z-env": "z", "a-env": "a"},
		},
	}

	if err := r.Report(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	// header + a-env + z-env (sorted)
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[1], "ZEBRA,mismatch,a-env") {
		t.Errorf("expected a-env first, got: %s", lines[1])
	}
}
