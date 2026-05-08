package annotator_test

import (
	"testing"

	"github.com/user/envdiff/internal/annotator"
	"github.com/user/envdiff/internal/comparator"
)

func makeResult(key string, statuses map[string]comparator.Status, values map[string]string) comparator.Result {
	return comparator.Result{Key: key, Statuses: statuses, Values: values}
}

func TestAnnotate_NoIssues(t *testing.T) {
	a := annotator.New()
	results := []comparator.Result{
		makeResult("APP_ENV", map[string]comparator.Status{"prod": comparator.StatusMatch}, map[string]string{"prod": "production"}),
	}
	anns := a.Annotate(results)
	if len(anns) != 0 {
		t.Fatalf("expected 0 annotations, got %d", len(anns))
	}
}

func TestAnnotate_MissingKey(t *testing.T) {
	a := annotator.New()
	results := []comparator.Result{
		makeResult("DB_HOST", map[string]comparator.Status{"staging": comparator.StatusMissing}, map[string]string{}),
	}
	anns := a.Annotate(results)
	if len(anns) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(anns))
	}
	if anns[0].Severity != "error" {
		t.Errorf("expected severity 'error', got %q", anns[0].Severity)
	}
	if anns[0].Key != "DB_HOST" {
		t.Errorf("expected key 'DB_HOST', got %q", anns[0].Key)
	}
}

func TestAnnotate_MismatchedPlainKey(t *testing.T) {
	a := annotator.New()
	results := []comparator.Result{
		makeResult("LOG_LEVEL", map[string]comparator.Status{"prod": comparator.StatusMismatch}, map[string]string{"prod": "debug"}),
	}
	anns := a.Annotate(results)
	if len(anns) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(anns))
	}
	if anns[0].Severity != "warn" {
		t.Errorf("expected severity 'warn', got %q", anns[0].Severity)
	}
}

func TestAnnotate_MismatchedSensitiveKey(t *testing.T) {
	a := annotator.New()
	results := []comparator.Result{
		makeResult("API_TOKEN", map[string]comparator.Status{"prod": comparator.StatusMismatch}, map[string]string{"prod": "abc"}),
	}
	anns := a.Annotate(results)
	if len(anns) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(anns))
	}
	if anns[0].Severity != "error" {
		t.Errorf("expected severity 'error' for sensitive key, got %q", anns[0].Severity)
	}
}

func TestAnnotate_CustomPatterns(t *testing.T) {
	a := annotator.NewWithPatterns([]string{"internal"})
	results := []comparator.Result{
		makeResult("INTERNAL_URL", map[string]comparator.Status{"dev": comparator.StatusMismatch}, map[string]string{"dev": "http://x"}),
	}
	anns := a.Annotate(results)
	if len(anns) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(anns))
	}
	if anns[0].Severity != "error" {
		t.Errorf("expected severity 'error' for custom sensitive key, got %q", anns[0].Severity)
	}
}
