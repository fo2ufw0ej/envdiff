package reporter

import (
	"fmt"
	"io"
	"time"

	"github.com/yourorg/envdiff/internal/comparator"
	"github.com/yourorg/envdiff/internal/watcher"
)

// WatchReporter wraps an existing Reporter and re-runs a comparison
// each time a ChangeEvent arrives, writing results to the provided writer.
type WatchReporter struct {
	inner   Reporter
	compare func() ([]comparator.Result, error)
	out     io.Writer
}

// Reporter is the shared interface implemented by all envdiff reporters.
type Reporter interface {
	Write(results []comparator.Result) error
}

// NewWatchReporter creates a WatchReporter that delegates formatting to inner.
func NewWatchReporter(inner Reporter, compare func() ([]comparator.Result, error), out io.Writer) *WatchReporter {
	return &WatchReporter{inner: inner, compare: compare, out: out}
}

// Run listens for ChangeEvents and re-renders the diff report on each event.
// It returns when the events channel is closed.
func (wr *WatchReporter) Run(events <-chan watcher.ChangeEvent) error {
	for ev := range events {
		fmt.Fprintf(wr.out, "\n--- envdiff reloading [%s] ---\n", ev.At.Format(time.RFC3339))
		for _, p := range ev.ChangedPaths {
			fmt.Fprintf(wr.out, "  changed: %s\n", p)
		}
		results, err := wr.compare()
		if err != nil {
			fmt.Fprintf(wr.out, "error: %v\n", err)
			continue
		}
		if err := wr.inner.Write(results); err != nil {
			fmt.Fprintf(wr.out, "reporter error: %v\n", err)
		}
	}
	return nil
}
