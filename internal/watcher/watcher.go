package watcher

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

// FileState holds the last-known checksum and modification time of a file.
type FileState struct {
	Path    string
	ModTime time.Time
	Checksum string
}

// ChangeEvent is emitted when one or more watched files change.
type ChangeEvent struct {
	ChangedPaths []string
	At           time.Time
}

// Watcher polls a set of files and notifies via a channel when any file changes.
type Watcher struct {
	paths    []string
	interval time.Duration
	states   map[string]FileState
}

// New creates a Watcher for the given file paths and poll interval.
func New(paths []string, interval time.Duration) *Watcher {
	return &Watcher{
		paths:    paths,
		interval: interval,
		states:   make(map[string]FileState),
	}
}

// Watch starts polling and sends ChangeEvents to the returned channel.
// It stops when the done channel is closed.
func (w *Watcher) Watch(done <-chan struct{}) <-chan ChangeEvent {
	ch := make(chan ChangeEvent, 1)
	// Seed initial state.
	for _, p := range w.paths {
		if s, err := stat(p); err == nil {
			w.states[p] = s
		}
	}
	go func() {
		defer close(ch)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				changed := w.poll()
				if len(changed) > 0 {
					ch <- ChangeEvent{ChangedPaths: changed, At: t}
				}
			}
		}
	}()
	return ch
}

// Add appends a new file path to the set of watched paths. If the file is
// accessible, its current state is recorded so the next poll can detect changes.
func (w *Watcher) Add(path string) {
	for _, p := range w.paths {
		if p == path {
			return // already watching
		}
	}
	w.paths = append(w.paths, path)
	if s, err := stat(path); err == nil {
		w.states[path] = s
	}
}

func (w *Watcher) poll() []string {
	var changed []string
	for _, p := range w.paths {
		s, err := stat(p)
		if err != nil {
			continue
		}
		prev, known := w.states[p]
		if !known || s.Checksum != prev.Checksum {
			w.states[p] = s
			changed = append(changed, p)
		}
	}
	return changed
}

func stat(path string) (FileState, error) {
	info, err := os.Stat(path)
	if err != nil {
		return FileState{}, err
	}
	checksum, err := fileChecksum(path)
	if err != nil {
		return FileState{}, err
	}
	return FileState{Path: path, ModTime: info.ModTime(), Checksum: checksum}, nil
}

func fileChecksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
