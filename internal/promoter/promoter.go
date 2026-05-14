// Package promoter copies key-value pairs from one environment map into
// another, optionally overwriting existing values.
package promoter

import "fmt"

// Strategy controls how conflicts are resolved when a key already exists in
// the destination environment.
type Strategy int

const (
	// SkipExisting leaves destination values untouched when the key is already
	// present.
	SkipExisting Strategy = iota
	// OverwriteExisting replaces destination values with the source value.
	OverwriteExisting
)

// Result describes the outcome of promoting a single key.
type Result struct {
	Key      string
	OldValue string // empty when the key was absent in destination
	NewValue string
	Skipped  bool // true when strategy == SkipExisting and key existed
}

// Promote copies keys from src into dst according to the given strategy.
// It returns a Result for every key present in src.
func Promote(src, dst map[string]string, strategy Strategy) (map[string]string, []Result, error) {
	if src == nil {
		return nil, nil, fmt.Errorf("promoter: src must not be nil")
	}
	if dst == nil {
		dst = make(map[string]string)
	}

	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	results := make([]Result, 0, len(src))
	for k, v := range src {
		old, exists := out[k]
		if exists && strategy == SkipExisting {
			results = append(results, Result{Key: k, OldValue: old, NewValue: old, Skipped: true})
			continue
		}
		out[k] = v
		results = append(results, Result{Key: k, OldValue: old, NewValue: v, Skipped: false})
	}
	return out, results, nil
}
