// Package censor replaces specific key values in comparison results
// with a fixed placeholder string, allowing sensitive data to be hidden
// from output without removing the key entirely.
package censor

import (
	"strings"

	"github.com/user/envdiff/internal/comparator"
)

const defaultPlaceholder = "[CENSORED]"

// Censor holds the set of keys to censor and the placeholder string.
type Censor struct {
	keys        map[string]struct{}
	placeholder string
}

// New returns a Censor that replaces values for the given keys with
// defaultPlaceholder. Key matching is case-insensitive.
func New(keys []string) *Censor {
	return NewWithPlaceholder(keys, defaultPlaceholder)
}

// NewWithPlaceholder returns a Censor with a custom placeholder string.
// If placeholder is empty, defaultPlaceholder is used.
func NewWithPlaceholder(keys []string, placeholder string) *Censor {
	if placeholder == "" {
		placeholder = defaultPlaceholder
	}
	km := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		km[strings.ToLower(k)] = struct{}{}
	}
	return &Censor{keys: km, placeholder: placeholder}
}

// Apply returns a new slice of results where any value belonging to a
// censored key is replaced with the placeholder. The original slice is
// not mutated.
func (c *Censor) Apply(results []comparator.Result) []comparator.Result {
	out := make([]comparator.Result, len(results))
	for i, r := range results {
		if _, ok := c.keys[strings.ToLower(r.Key)]; !ok {
			out[i] = r
			continue
		}
		// Deep-copy the Values map with censored values.
		censored := make(map[string]string, len(r.Values))
		for env, val := range r.Values {
			if val != "" {
				censored[env] = c.placeholder
			} else {
				censored[env] = val
			}
		}
		out[i] = comparator.Result{
			Key:    r.Key,
			Status: r.Status,
			Values: censored,
		}
	}
	return out
}
