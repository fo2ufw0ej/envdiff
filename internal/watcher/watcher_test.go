package watcher_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envdiff/internal/watcher"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestNew_NoChanges(t *testing.T) {
	dir := t.TempDir()
	p := writeTempFile(t, dir, ".env", "KEY=value\n")

	w := watcher.New([]string{p}, 20*time.Millisecond)
	done := make(chan struct{})
	ch := w.Watch(done)

	select {
	case ev := <-ch:
		t.Fatalf("unexpected change event: %v", ev)
	case <-time.After(80 * time.Millisecond):
		// expected: no changes
	}
	close(done)
}

func TestWatch_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	p := writeTempFile(t, dir, ".env", "KEY=original\n")

	w := watcher.New([]string{p}, 20*time.Millisecond)
	done := make(chan struct{})
	ch := w.Watch(done)
	defer close(done)

	// Give the watcher time to seed state.
	time.Sleep(30 * time.Millisecond)

	// Modify the file.
	if err := os.WriteFile(p, []byte("KEY=changed\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	select {
	case ev := <-ch:
		if len(ev.ChangedPaths) != 1 || ev.ChangedPaths[0] != p {
			t.Errorf("expected changed path %q, got %v", p, ev.ChangedPaths)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatch_MultipleFiles(t *testing.T) {
	dir := t.TempDir()
	p1 := writeTempFile(t, dir, ".env.dev", "A=1\n")
	p2 := writeTempFile(t, dir, ".env.prod", "A=2\n")

	w := watcher.New([]string{p1, p2}, 20*time.Millisecond)
	done := make(chan struct{})
	ch := w.Watch(done)
	defer close(done)

	time.Sleep(30 * time.Millisecond)

	// Modify only p2.
	if err := os.WriteFile(p2, []byte("A=99\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	select {
	case ev := <-ch:
		if len(ev.ChangedPaths) != 1 || ev.ChangedPaths[0] != p2 {
			t.Errorf("expected only p2 changed, got %v", ev.ChangedPaths)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for change event")
	}
}
