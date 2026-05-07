package watcher_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envdiff/internal/loader"
	"github.com/yourorg/envdiff/internal/watcher"
)

// TestWatcher_WithLoader verifies that a change event can be acted upon
// by re-loading the affected files through the loader package.
func TestWatcher_WithLoader(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env.staging")
	if err := os.WriteFile(p, []byte("DB_HOST=localhost\nDB_PORT=5432\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	w := watcher.New([]string{p}, 20*time.Millisecond)
	done := make(chan struct{})
	ch := w.Watch(done)
	defer close(done)

	time.Sleep(30 * time.Millisecond)

	// Simulate an environment update.
	if err := os.WriteFile(p, []byte("DB_HOST=db.prod\nDB_PORT=5432\nDB_NAME=mydb\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	select {
	case ev := <-ch:
		envs, err := loader.LoadFiles(ev.ChangedPaths)
		if err != nil {
			t.Fatalf("LoadFiles: %v", err)
		}
		env, ok := envs["staging"]
		if !ok {
			t.Fatalf("expected 'staging' env, got keys: %v", envs)
		}
		if env["DB_HOST"] != "db.prod" {
			t.Errorf("expected DB_HOST=db.prod, got %q", env["DB_HOST"])
		}
		if _, ok := env["DB_NAME"]; !ok {
			t.Error("expected DB_NAME key to be present after reload")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for change event")
	}
}
