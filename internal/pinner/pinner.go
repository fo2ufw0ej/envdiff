// Package pinner allows users to pin specific keys so that their values
// are treated as authoritative and never flagged as mismatched during
// comparison. Pinned keys are excluded from mismatch reporting while
// still being included in missing-key checks.
package pinner

import (
	"strings"

	"github.com/user/envdiff/internal/comparator"
)

// Pinner holds the set of pinned keys.
type Pinner struct {
	keys map[string]struct{}
}

// New creates a Pinner from the supplied list of key names.
// Key matching is case-insensitive.
func New(keys []string) *Pinner {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[strings.ToUpper(k)] = struct{}{}
	}
	return &Pinner{keys: m}
}

// IsPinned reports whether key is pinned.
func (p *Pinner) IsPinned(key string) bool {
	_, ok := p.keys[strings.ToUpper(key)]
	return ok
}

// Apply returns a new slice of results with mismatch entries removed for
// pinned keys. Missing entries are preserved regardless of pin status.
func (p *Pinner) Apply(results []comparator.Result) []comparator.Result {
	out := make([]comparator.Result, 0, len(results))
	for _, r := range results {
		if r.Status == comparator.StatusMismatch && p.IsPinned(r.Key) {
			continue
		}
		out = append(out, r)
	}
	return out
}
