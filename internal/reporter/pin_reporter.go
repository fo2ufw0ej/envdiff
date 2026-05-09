package reporter

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/pinner"
)

// NewPinReporter returns a Reporter that lists which keys are currently
// pinned before delegating to the wrapped reporter for full output.
func NewPinReporter(w io.Writer, p *pinner.Pinner, keys []string, inner Reporter) Reporter {
	return &pinReporter{w: w, pinner: p, keys: keys, inner: inner}
}

type pinReporter struct {
	w      io.Writer
	pinner *pinner.Pinner
	keys   []string
	inner  Reporter
}

func (r *pinReporter) Report(results []Result) error {
	pinned := make([]string, 0, len(r.keys))
	for _, k := range r.keys {
		if r.pinner.IsPinned(k) {
			pinned = append(pinned, k)
		}
	}
	sort.Strings(pinned)

	if len(pinned) == 0 {
		fmt.Fprintln(r.w, "# Pinned keys: none")
	} else {
		fmt.Fprintf(r.w, "# Pinned keys (%d): %v\n", len(pinned), pinned)
	}

	filtered := r.pinner.Apply(results)
	if r.inner != nil {
		return r.inner.Report(filtered)
	}
	return nil
}
