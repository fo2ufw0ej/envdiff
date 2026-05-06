package exporter_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/exporter"
)

func makeResult(key string, values map[string]string, status comparator.Status) comparator.Result {
	return comparator.Result{Key: key, Values: values, Status: status}
}

func TestExporter_EnvFormat_WritesReferenceValues(t *testing.T) {
	results := []comparator.Result{
		makeResult("DB_HOST", map[string]string{"prod": "db.prod.example.com", "staging": "db.staging.example.com"}, comparator.StatusMismatch),
		makeResult("API_KEY", map[string]string{"prod": "secret123"}, comparator.StatusMissing),
	}

	ex := exporter.New(exporter.FormatEnv, "prod")
	var buf strings.Builder
	if err := ex.Export(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "API_KEY=secret123") {
		t.Errorf("expected API_KEY in output, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST=db.prod.example.com") {
		t.Errorf("expected DB_HOST in output, got:\n%s", out)
	}
}

func TestExporter_EnvFormat_QuotesSpaces(t *testing.T) {
	results := []comparator.Result{
		makeResult("APP_NAME", map[string]string{"prod": "My App"}, comparator.StatusMissing),
	}

	ex := exporter.New(exporter.FormatEnv, "prod")
	var buf strings.Builder
	if err := ex.Export(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `APP_NAME="My App"`) {
		t.Errorf("expected quoted value, got:\n%s", out)
	}
}

func TestExporter_PatchFormat(t *testing.T) {
	results := []comparator.Result{
		makeResult("TIMEOUT", map[string]string{"prod": "30"}, comparator.StatusMissing),
	}

	ex := exporter.New(exporter.FormatPatch, "prod")
	var buf strings.Builder
	if err := ex.Export(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "+ TIMEOUT=30") {
		t.Errorf("expected patch line, got:\n%s", out)
	}
}

func TestExporter_UnsupportedFormat(t *testing.T) {
	ex := exporter.New("xml", "prod")
	var buf strings.Builder
	err := ex.Export(&buf, nil)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestExporter_SkipsMissingReference(t *testing.T) {
	results := []comparator.Result{
		makeResult("ONLY_STAGING", map[string]string{"staging": "value"}, comparator.StatusMissing),
	}

	ex := exporter.New(exporter.FormatEnv, "prod")
	var buf strings.Builder
	if err := ex.Export(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf("expected empty output when reference key absent, got: %s", buf.String())
	}
}
