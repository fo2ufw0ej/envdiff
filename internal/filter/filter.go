package filter

import (
	"strings"

	"github.com/user/envdiff/internal/comparator"
)

// Options holds filtering criteria for comparison results.
type Options struct {
	// OnlyMissing restricts output to keys that are missing in at least one environment.
	OnlyMissing bool
	// OnlyMismatched restricts output to keys whose values differ across environments.
	OnlyMismatched bool
	// KeyPrefix filters results to keys that start with the given prefix (case-insensitive).
	KeyPrefix string
}

// Apply returns a new slice of DiffEntry values that match the given Options.
func Apply(entries []comparator.DiffEntry, opts Options) []comparator.DiffEntry {
	var out []comparator.DiffEntry
	for _, e := range entries {
		if opts.KeyPrefix != "" {
			if !strings.HasPrefix(strings.ToUpper(e.Key), strings.ToUpper(opts.KeyPrefix)) {
				continue
			}
		}
		if opts.OnlyMissing && !isMissing(e) {
			continue
		}
		if opts.OnlyMismatched && !isMismatched(e) {
			continue
		}
		out = append(out, e)
	}
	return out
}

// isMissing returns true when at least one environment has no value for the key.
func isMissing(e comparator.DiffEntry) bool {
	for _, v := range e.Values {
		if v == "" {
			return true
		}
	}
	return false
}

// isMismatched returns true when the values across environments are not all identical.
func isMismatched(e comparator.DiffEntry) bool {
	var first *string
	for _, v := range e.Values {
		if first == nil {
			copy := v
			first = &copy
			continue
		}
		if v != *first {
			return true
		}
	}
	return false
}
