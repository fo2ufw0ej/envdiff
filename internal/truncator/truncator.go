// Package truncator shortens long env values to a maximum display length,
// optionally appending a suffix such as "..." to indicate truncation.
package truncator

import "github.com/your-org/envdiff/internal/comparator"

// Options configures truncation behaviour.
type Options struct {
	// MaxLen is the maximum number of runes to keep. Defaults to 64.
	MaxLen int
	// Suffix is appended when a value is truncated. Defaults to "...".
	Suffix string
}

func defaults(o Options) Options {
	if o.MaxLen <= 0 {
		o.MaxLen = 64
	}
	if o.Suffix == "" {
		o.Suffix = "..."
	}
	return o
}

// Truncate returns a copy of results where every value longer than
// Options.MaxLen runes is cut and suffixed with Options.Suffix.
// The original slice is never mutated.
func Truncate(results []comparator.Result, o Options) []comparator.Result {
	o = defaults(o)
	out := make([]comparator.Result, len(results))
	for i, r := range results {
		out[i] = comparator.Result{
			Key:    r.Key,
			Status: r.Status,
			Values: truncateValues(r.Values, o),
		}
	}
	return out
}

func truncateValues(vals map[string]string, o Options) map[string]string {
	if vals == nil {
		return nil
	}
	out := make(map[string]string, len(vals))
	for env, v := range vals {
		out[env] = truncateString(v, o)
	}
	return out
}

func truncateString(s string, o Options) string {
	runes := []rune(s)
	if len(runes) <= o.MaxLen {
		return s
	}
	return string(runes[:o.MaxLen]) + o.Suffix
}
