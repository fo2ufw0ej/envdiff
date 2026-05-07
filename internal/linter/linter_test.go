package linter_test

import (
	"testing"

	"github.com/user/envdiff/internal/linter"
)

func findingMessages(findings []linter.Finding) []string {
	out := make([]string, len(findings))
	for i, f := range findings {
		out[i] = f.Message
	}
	return out
}

func TestLint_CleanEnv(t *testing.T) {
	l := linter.New()
	env := map[string]string{"DB_HOST": "localhost", "PORT": "8080"}
	findings := l.Lint(env)
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %v", findingMessages(findings))
	}
}

func TestLint_EmptyValue(t *testing.T) {
	l := linter.New()
	env := map[string]string{"DB_PASS": ""}
	findings := l.Lint(env)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != linter.SeverityWarn {
		t.Errorf("expected warn severity, got %s", findings[0].Severity)
	}
}

func TestLint_WhitespaceInKey(t *testing.T) {
	l := linter.New()
	env := map[string]string{"bad key": "value"}
	findings := l.Lint(env)
	// expect at least the whitespace error (lower-case rule may also fire)
	var found bool
	for _, f := range findings {
		if f.Severity == linter.SeverityError {
			found = true
		}
	}
	if !found {
		t.Error("expected at least one error-severity finding for whitespace in key")
	}
}

func TestLint_LowerCaseKey(t *testing.T) {
	l := linter.New()
	env := map[string]string{"db_host": "localhost"}
	findings := l.Lint(env)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != linter.SeverityWarn {
		t.Errorf("expected warn, got %s", findings[0].Severity)
	}
}

func TestLint_CustomRule(t *testing.T) {
	noUnderscoreRule := func(key, _ string) []linter.Finding {
		if len(key) > 0 && key[0] == '_' {
			return []linter.Finding{{Key: key, Message: "key starts with underscore", Severity: linter.SeverityWarn}}
		}
		return nil
	}
	l := linter.NewWithRules(noUnderscoreRule)
	env := map[string]string{"_HIDDEN": "yes", "VISIBLE": "no"}
	findings := l.Lint(env)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "_HIDDEN" {
		t.Errorf("unexpected key %s", findings[0].Key)
	}
}
