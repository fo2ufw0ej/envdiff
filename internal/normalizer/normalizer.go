package normalizer

import (
	"strings"

	"github.com/user/envdiff/internal/comparator"
)

// Option controls how normalization is applied.
type Option func(*options)

type options struct {
	upperKeys   bool
	trimValues  bool
	stripPrefix string
}

// WithUpperKeys converts all keys to UPPER_CASE before comparison.
func WithUpperKeys() Option {
	return func(o *options) { o.upperKeys = true }
}

// WithTrimValues trims leading/trailing whitespace from all values.
func WithTrimValues() Option {
	return func(o *options) { o.trimValues = true }
}

// WithStripPrefix removes a common prefix from all keys.
func WithStripPrefix(prefix string) Option {
	return func(o *options) { o.stripPrefix = prefix }
}

// Normalize applies the given options to a slice of comparator.Result,
// returning a new slice with transformed keys and/or values.
func Normalize(results []comparator.Result, opts ...Option) []comparator.Result {
	cfg := &options{}
	for _, o := range opts {
		o(cfg)
	}

	out := make([]comparator.Result, 0, len(results))
	for _, r := range results {
		nr := comparator.Result{
			Key:    normalizeKey(r.Key, cfg),
			Values: make(map[string]string, len(r.Values)),
		}
		for env, val := range r.Values {
			if cfg.trimValues {
				val = strings.TrimSpace(val)
			}
			nr.Values[env] = val
		}
		nr.Status = r.Status
		out = append(out, nr)
	}
	return out
}

func normalizeKey(key string, cfg *options) string {
	if cfg.stripPrefix != "" {
		key = strings.TrimPrefix(key, cfg.stripPrefix)
	}
	if cfg.upperKeys {
		key = strings.ToUpper(key)
	}
	return key
}
