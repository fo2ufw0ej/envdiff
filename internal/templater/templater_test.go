package templater_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/templater"
)

func TestRender_AllPresent(t *testing.T) {
	tmpl := templater.New()
	env := map[string]string{"DB_HOST": "localhost", "PORT": "8080"}
	result := tmpl.Render("dev", []string{"DB_HOST", "PORT"}, env)

	if result.EnvName != "dev" {
		t.Fatalf("expected env name dev, got %s", result.EnvName)
	}
	if len(result.Missing) != 0 {
		t.Fatalf("expected no missing keys, got %v", result.Missing)
	}
	if len(result.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(result.Lines))
	}
	for _, line := range result.Lines {
		if !strings.Contains(line, "=<") {
			t.Errorf("expected placeholder in line: %s", line)
		}
	}
}

func TestRender_MissingKey(t *testing.T) {
	tmpl := templater.New()
	env := map[string]string{"DB_HOST": "localhost"}
	result := tmpl.Render("prod", []string{"DB_HOST", "SECRET_KEY"}, env)

	if len(result.Missing) != 1 || result.Missing[0] != "SECRET_KEY" {
		t.Fatalf("expected SECRET_KEY in missing, got %v", result.Missing)
	}
	if len(result.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(result.Lines))
	}
}

func TestRender_CustomPlaceholder(t *testing.T) {
	tmpl := templater.NewWithPlaceholder(func(key string) string {
		return "CHANGEME"
	})
	env := map[string]string{"API_KEY": "abc"}
	result := tmpl.Render("staging", []string{"API_KEY"}, env)

	if result.Lines[0] != "API_KEY=CHANGEME" {
		t.Errorf("unexpected line: %s", result.Lines[0])
	}
}

func TestRenderAll_SortsEnvsAndKeys(t *testing.T) {
	tmpl := templater.New()
	envs := map[string]map[string]string{
		"prod": {"Z_KEY": "1", "A_KEY": "2"},
		"dev":  {"Z_KEY": "3"},
	}
	results := tmpl.RenderAll(envs)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].EnvName != "dev" {
		t.Errorf("expected dev first, got %s", results[0].EnvName)
	}
	// dev is missing A_KEY
	if len(results[0].Missing) != 1 || results[0].Missing[0] != "A_KEY" {
		t.Errorf("expected A_KEY missing in dev, got %v", results[0].Missing)
	}
	// first line should be A_KEY (sorted)
	if !strings.HasPrefix(results[0].Lines[0], "A_KEY=") {
		t.Errorf("expected A_KEY first, got %s", results[0].Lines[0])
	}
}

func TestNewWithPlaceholder_NilFallsBack(t *testing.T) {
	tmpl := templater.NewWithPlaceholder(nil)
	env := map[string]string{"FOO": "bar"}
	result := tmpl.Render("test", []string{"FOO"}, env)
	if !strings.Contains(result.Lines[0], "<FOO>") {
		t.Errorf("expected default placeholder, got %s", result.Lines[0])
	}
}
