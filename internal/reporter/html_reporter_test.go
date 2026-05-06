package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/reporter"
)

func TestHTMLReporter_NoDifferences(t *testing.T) {
	var buf bytes.Buffer
	report := reporter.NewHTMLReporter(&buf)
	if err := report(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "No differences found.") {
		t.Errorf("expected no-differences message, got:\n%s", out)
	}
	if !strings.Contains(out, "</html>") {
		t.Errorf("expected closing html tag, got:\n%s", out)
	}
}

func TestHTMLReporter_MissingKey(t *testing.T) {
	var buf bytes.Buffer
	report := reporter.NewHTMLReporter(&buf)
	diffs := []comparator.Diff{
		{
			Key:  "DB_HOST",
			Type: comparator.Missing,
			Values: map[string]string{
				"production": "db.prod.example.com",
			},
		},
	}
	if err := report(diffs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected key DB_HOST in output")
	}
	if !strings.Contains(out, "missing") {
		t.Errorf("expected row class 'missing' in output")
	}
	if !strings.Contains(out, "<em>absent</em>") {
		t.Errorf("expected absent cell for missing env")
	}
	if !strings.Contains(out, "db.prod.example.com") {
		t.Errorf("expected value in output")
	}
}

func TestHTMLReporter_MismatchedKey(t *testing.T) {
	var buf bytes.Buffer
	report := reporter.NewHTMLReporter(&buf)
	diffs := []comparator.Diff{
		{
			Key:  "LOG_LEVEL",
			Type: comparator.Mismatch,
			Values: map[string]string{
				"staging":    "debug",
				"production": "error",
			},
		},
	}
	if err := report(diffs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "LOG_LEVEL") {
		t.Errorf("expected key LOG_LEVEL in output")
	}
	if !strings.Contains(out, "mismatch") {
		t.Errorf("expected row class 'mismatch' in output")
	}
	if !strings.Contains(out, "debug") || !strings.Contains(out, "error") {
		t.Errorf("expected both values in output")
	}
	if !strings.Contains(out, "<table>") {
		t.Errorf("expected table element in output")
	}
}
