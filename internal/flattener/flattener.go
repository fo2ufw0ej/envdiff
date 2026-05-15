// Package flattener merges all comparison results into a single flat map
// keyed by "KEY@ENV", making it easy to feed results into tabular tools.
package flattener

import (
	"fmt"
	"sort"

	"github.com/user/envdiff/internal/comparator"
)

// Entry is a single flattened cell from the comparison matrix.
type Entry struct {
	Key    string
	Env    string
	Status comparator.Status
	Value  string // empty when status is Missing
}

// Flatten converts a slice of comparator.Result into a flat list of Entry
// values, one per (key, env) pair. Results are sorted by key then env for
// deterministic output.
func Flatten(results []comparator.Result) []Entry {
	var entries []Entry

	for _, r := range results {
		for env, val := range r.Values {
			entries = append(entries, Entry{
				Key:    r.Key,
				Env:    env,
				Status: r.Status,
				Value:  val,
			})
		}
		// Emit explicit Missing entries for envs that have no value recorded.
		for _, env := range r.MissingIn {
			if _, exists := r.Values[env]; !exists {
				entries = append(entries, Entry{
					Key:    r.Key,
					Env:    env,
					Status: comparator.StatusMissing,
					Value:  "",
				})
			}
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		ki := fmt.Sprintf("%s\x00%s", entries[i].Key, entries[i].Env)
		kj := fmt.Sprintf("%s\x00%s", entries[j].Key, entries[j].Env)
		return ki < kj
	})

	return entries
}

// Index returns a map from "KEY@ENV" to Entry for O(1) lookups.
func Index(entries []Entry) map[string]Entry {
	m := make(map[string]Entry, len(entries))
	for _, e := range entries {
		m[fmt.Sprintf("%s@%s", e.Key, e.Env)] = e
	}
	return m
}
