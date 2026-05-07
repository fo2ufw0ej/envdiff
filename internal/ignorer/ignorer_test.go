package ignorer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/ignorer"
)

func TestShouldIgnore_ExactKey(t *testing.T) {
	ig := ignorer.New([]string{"SECRET", "TOKEN"}, nil)
	if !ig.ShouldIgnore("SECRET") {
		t.Error("expected SECRET to be ignored")
	}
	if !ig.ShouldIgnore("token") { // case-insensitive
		t.Error("expected token (lowercase) to be ignored")
	}
	if ig.ShouldIgnore("HOST") {
		t.Error("expected HOST not to be ignored")
	}
}

func TestShouldIgnore_Prefix(t *testing.T) {
	ig := ignorer.New(nil, []string{"AWS_"})
	if !ig.ShouldIgnore("AWS_ACCESS_KEY_ID") {
		t.Error("expected AWS_ACCESS_KEY_ID to be ignored")
	}
	if !ig.ShouldIgnore("aws_secret") {
		t.Error("expected aws_secret (lowercase) to be ignored")
	}
	if ig.ShouldIgnore("DATABASE_URL") {
		t.Error("expected DATABASE_URL not to be ignored")
	}
}

func TestShouldIgnore_Empty(t *testing.T) {
	ig := ignorer.New(nil, nil)
	if ig.ShouldIgnore("ANYTHING") {
		t.Error("empty ignorer should not ignore any key")
	}
}

func writeIgnoreFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".envignore")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadFile_Basic(t *testing.T) {
	p := writeIgnoreFile(t, "# comment\nSECRET\nAWS_*\n\nTOKEN\n")
	ig, err := ignorer.LoadFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, key := range []string{"SECRET", "TOKEN", "AWS_REGION", "aws_key"} {
		if !ig.ShouldIgnore(key) {
			t.Errorf("expected %q to be ignored", key)
		}
	}
	if ig.ShouldIgnore("HOST") {
		t.Error("expected HOST not to be ignored")
	}
}

func TestLoadFile_Missing(t *testing.T) {
	_, err := ignorer.LoadFile("/nonexistent/.envignore")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
