package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/comparator"
	"github.com/yourorg/envdiff/internal/reporter"
)

func TestTextReporter_NoDifferences(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.NewTextReporter(&buf)

	err := r.Write(comparator.Result{
		MissingIn:  map[string][]string{},
		Mismatched: nil,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestTextReporter_MissingKey(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.NewTextReporter(&buf)

	result := comparator.Result{
		MissingIn: map[string][]string{
			"prod": {"SECRET_KEY"},
		},
	}

	if err := r.Write(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "MISSING") || !strings.Contains(out, "prod") || !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestTextReporter_MismatchedKey(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.NewTextReporter(&buf)

	result := comparator.Result{
		MissingIn: map[string][]string{},
		Mismatched: []comparator.MismatchedKey{
			{
				Key:    "LOG_LEVEL",
				Values: map[string]string{"dev": "debug", "prod": "error"},
			},
		},
	}

	if err := r.Write(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "MISMATCH") || !strings.Contains(out, "LOG_LEVEL") {
		t.Errorf("unexpected output: %s", out)
	}
	if !strings.Contains(out, "debug") || !strings.Contains(out, "error") {
		t.Errorf("expected both values in output: %s", out)
	}
}
