// Package trimmer removes keys from comparison results that match a
// provided allowlist, blocklist, or value-length constraint. It is useful
// for narrowing large result sets before reporting.
package trimmer

import (
	"strings"

	"github.com/your-org/envdiff/internal/comparator"
)

// Options controls which results are kept after trimming.
type Options struct {
	// AllowKeys, if non-empty, retains only results whose key is in the set.
	AllowKeys []string

	// BlockKeys removes results whose key is in the set.
	BlockKeys []string

	// MaxValueLen drops results where every known value exceeds this length.
	// Zero means no limit.
	MaxValueLen int
}

// Trim applies the given Options to results and returns the filtered slice.
// The original slice is never mutated.
func Trim(results []comparator.Result, opts Options) []comparator.Result {
	allowSet := toSet(opts.AllowKeys)
	blockSet := toSet(opts.BlockKeys)

	out := make([]comparator.Result, 0, len(results))
	for _, r := range results {
		key := r.Key

		if len(allowSet) > 0 && !allowSet[strings.ToUpper(key)] {
			continue
		}
		if blockSet[strings.ToUpper(key)] {
			continue
		}
		if opts.MaxValueLen > 0 && allValuesExceedLen(r.Values, opts.MaxValueLen) {
			continue
		}

		out = append(out, r)
	}
	return out
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[strings.ToUpper(k)] = true
	}
	return s
}

func allValuesExceedLen(values map[string]string, max int) bool {
	if len(values) == 0 {
		return false
	}
	for _, v := range values {
		if len(v) <= max {
			return false
		}
	}
	return true
}
