package loader

import (
	"os"
	"path/filepath"
	"testing"
)

func writeEnvFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeEnvFile: %v", err)
	}
	return p
}

func TestLoadFiles_Basic(t *testing.T) {
	dir := t.TempDir()
	p1 := writeEnvFile(t, dir, ".env", "KEY1=val1\nKEY2=val2\n")
	p2 := writeEnvFile(t, dir, ".env.staging", "KEY1=stg1\nKEY3=stg3\n")

	em, err := LoadFiles([]string{p1, p2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(em) != 2 {
		t.Fatalf("expected 2 envs, got %d", len(em))
	}
	if em["default"]["KEY1"] != "val1" {
		t.Errorf("default KEY1: got %q, want %q", em["default"]["KEY1"], "val1")
	}
	if em["staging"]["KEY3"] != "stg3" {
		t.Errorf("staging KEY3: got %q, want %q", em["staging"]["KEY3"], "stg3")
	}
}

func TestLoadFiles_MissingFile(t *testing.T) {
	_, err := LoadFiles([]string{"/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadDir_Basic(t *testing.T) {
	dir := t.TempDir()
	writeEnvFile(t, dir, ".env", "A=1\n")
	writeEnvFile(t, dir, ".env.production", "A=prod\n")
	writeEnvFile(t, dir, "README.md", "not an env file\n")

	em, err := LoadDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := em["default"]; !ok {
		t.Error("expected 'default' env")
	}
	if _, ok := em["production"]; !ok {
		t.Error("expected 'production' env")
	}
	if _, ok := em["README.md"]; ok {
		t.Error("README.md should not be loaded")
	}
}

func TestLoadDir_Empty(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadDir(dir)
	if err == nil {
		t.Fatal("expected error for empty dir, got nil")
	}
}

func TestEnvName(t *testing.T) {
	cases := []struct{ path, want string }{
		{".env", "default"},
		{".env.dev", "dev"},
		{"some/path/.env.production", "production"},
		{"custom.env", "custom.env"},
	}
	for _, c := range cases {
		got := envName(c.path)
		if got != c.want {
			t.Errorf("envName(%q) = %q, want %q", c.path, got, c.want)
		}
	}
}
