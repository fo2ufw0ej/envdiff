package redactor_test

import (
	"testing"

	"github.com/user/envdiff/internal/redactor"
)

func TestIsSensitive_DefaultPatterns(t *testing.T) {
	r := redactor.New(nil, "")

	sensitive := []string{"DB_PASSWORD", "API_KEY", "AUTH_TOKEN", "SECRET_KEY", "PRIVATE_KEY"}
	for _, k := range sensitive {
		if !r.IsSensitive(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}

	plain := []string{"APP_NAME", "PORT", "HOST", "LOG_LEVEL"}
	for _, k := range plain {
		if r.IsSensitive(k) {
			t.Errorf("expected %q to NOT be sensitive", k)
		}
	}
}

func TestIsSensitive_CustomPatterns(t *testing.T) {
	r := redactor.New([]string{"INTERNAL"}, "")
	if !r.IsSensitive("INTERNAL_URL") {
		t.Error("expected INTERNAL_URL to be sensitive")
	}
	if r.IsSensitive("API_KEY") {
		t.Error("expected API_KEY NOT to be sensitive with custom patterns")
	}
}

func TestRedactValue_MasksSensitive(t *testing.T) {
	r := redactor.New(nil, "REDACTED")
	got := r.RedactValue("DB_PASSWORD", "supersecret")
	if got != "REDACTED" {
		t.Errorf("got %q, want REDACTED", got)
	}
}

func TestRedactValue_LeavesPlainIntact(t *testing.T) {
	r := redactor.New(nil, "***")
	got := r.RedactValue("APP_NAME", "myapp")
	if got != "myapp" {
		t.Errorf("got %q, want myapp", got)
	}
}

func TestRedactMap(t *testing.T) {
	r := redactor.New(nil, "***")
	input := map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "secret123",
		"PORT":        "8080",
		"API_KEY":     "abc-def",
	}
	out := r.RedactMap(input)

	if out["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should be unchanged, got %q", out["APP_NAME"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("PORT should be unchanged, got %q", out["PORT"])
	}
	if out["DB_PASSWORD"] != "***" {
		t.Errorf("DB_PASSWORD should be redacted, got %q", out["DB_PASSWORD"])
	}
	if out["API_KEY"] != "***" {
		t.Errorf("API_KEY should be redacted, got %q", out["API_KEY"])
	}
	// original must not be mutated
	if input["DB_PASSWORD"] != "secret123" {
		t.Error("original map was mutated")
	}
}
