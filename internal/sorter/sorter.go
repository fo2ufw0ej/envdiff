// Package sorter provides utilities for sorting and ordering
// comparison results from the envdiff comparator.
package sorter

import (
	"sort"

	"github.com/user/envdiff/internal/comparator"
)

// SortBy defines the field by which results are sorted.
type SortBy int

const (
	// SortByKey sorts entries alphabetically by key name.
	SortByKey SortBy = iota
	// SortByStatus sorts entries by their difference status (missing before mismatched).
	SortByStatus
	// SortByEnv sorts entries by the first environment name that references them.
	SortByEnv
)

// Sort returns a new slice of Result entries ordered by the given SortBy strategy.
// The original slice is not modified.
func Sort(results []comparator.Result, by SortBy) []comparator.Result {
	copied := make([]comparator.Result, len(results))
	copy(copied, results)

	switch by {
	case SortByStatus:
		sort.SliceStable(copied, func(i, j int) bool {
			return statusRank(copied[i]) < statusRank(copied[j])
		})
	case SortByEnv:
		sort.SliceStable(copied, func(i, j int) bool {
			ei := firstEnv(copied[i])
			ej := firstEnv(copied[j])
			if ei != ej {
				return ei < ej
			}
			return copied[i].Key < copied[j].Key
		})
	default: // SortByKey
		sort.SliceStable(copied, func(i, j int) bool {
			return copied[i].Key < copied[j].Key
		})
	}

	return copied
}

// statusRank returns a numeric rank for a result's status to enable ordering.
func statusRank(r comparator.Result) int {
	for _, v := range r.Values {
		if v == "" {
			return 0 // missing
		}
	}
	return 1 // mismatched
}

// firstEnv returns the lexicographically smallest environment name in a result.
func firstEnv(r comparator.Result) string {
	keys := make([]string, 0, len(r.Values))
	for k := range r.Values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if len(keys) == 0 {
		return ""
	}
	return keys[0]
}
