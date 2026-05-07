// Package watcher provides file-system polling for .env files.
//
// It watches a set of files at a configurable interval and emits
// ChangeEvents whenever a file's content checksum differs from the
// previously observed state.  This allows tooling built on envdiff
// to react to live edits without relying on OS-level inotify APIs,
// keeping the implementation portable across platforms.
//
// Basic usage:
//
//	w := watcher.New([]string{".env", ".env.prod"}, 2*time.Second)
//	done := make(chan struct{})
//	for ev := range w.Watch(done) {
//		fmt.Println("changed:", ev.ChangedPaths)
//	}
package watcher
