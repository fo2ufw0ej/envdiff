// Package grouper groups comparison results by a chosen dimension
// (key prefix, status, or environment name) to aid bulk analysis.
package grouper

import (
	"sort"
	"strings"

	"github.com/yourorg/envdiff/internal/comparator"
)

// By controls the grouping dimension.
type By string

const (
	// ByPrefix groups results by the first segment of the key (before "_").
	ByPrefix By = "prefix"
	// ByStatus groups results by their comparison status.
	ByStatus By = "status"
	// ByEnv groups results by the first environment name in which a difference appears.
	ByEnv By = "env"
)

// Group is a named collection of comparison results.
type Group struct {
	Name    string
	Results []comparator.Result
}

// GroupBy partitions results according to the chosen dimension and returns
// groups sorted by name for deterministic output.
func GroupBy(results []comparator.Result, by By) []Group {
	buckets := make(map[string][]comparator.Result)

	for _, r := range results {
		key := bucketKey(r, by)
		buckets[key] = append(buckets[key], r)
	}

	names := make([]string, 0, len(buckets))
	for name := range buckets {
		names = append(names, name)
	}
	sort.Strings(names)

	groups := make([]Group, 0, len(names))
	for _, name := range names {
		groups = append(groups, Group{Name: name, Results: buckets[name]})
	}
	return groups
}

func bucketKey(r comparator.Result, by By) string {
	switch by {
	case ByStatus:
		return string(r.Status)
	case ByEnv:
		for env := range r.Values {
			return env
		}
		return "(none)"
	default: // ByPrefix
		parts := strings.SplitN(r.Key, "_", 2)
		if len(parts) == 0 || parts[0] == "" {
			return "(none)"
		}
		return parts[0]
	}
}
