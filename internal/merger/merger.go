// Package merger provides functionality for merging multiple env maps
// into a single reference map, resolving conflicts by priority or union.
package merger

import "sort"

// Strategy controls how conflicting values are resolved during a merge.
type Strategy int

const (
	// StrategyFirst keeps the value from the first env that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the value from the last env that defines the key.
	StrategyLast
)

// Result holds the merged key-value map and metadata about conflicts.
type Result struct {
	// Values is the merged environment map.
	Values map[string]string
	// Conflicts maps each conflicting key to the list of envs that defined it
	// with differing values.
	Conflicts map[string][]string
}

// Merge combines multiple named env maps into a single Result.
// The envs parameter is an ordered slice of (name, map) pairs represented
// as a slice of Env structs so callers control priority order.
func Merge(envs []Env, strategy Strategy) Result {
	result := Result{
		Values:    make(map[string]string),
		Conflicts: make(map[string][]string),
	}

	// Track which envs set each key and with what value.
	type entry struct {
		envName string
		value   string
	}
	seen := make(map[string][]entry)

	for _, env := range envs {
		for k, v := range env.Values {
			seen[k] = append(seen[k], entry{envName: env.Name, value: v})
		}
	}

	keys := sortedKeys(seen)
	for _, k := range keys {
		entries := seen[k]
		switch strategy {
		case StrategyLast:
			result.Values[k] = entries[len(entries)-1].value
		default: // StrategyFirst
			result.Values[k] = entries[0].value
		}

		// Detect conflicts: multiple envs with differing values.
		if hasConflict(entries) {
			names := make([]string, len(entries))
			for i, e := range entries {
				names[i] = e.envName
			}
			result.Conflicts[k] = names
		}
	}

	return result
}

// Env pairs a name with its key-value map.
type Env struct {
	Name   string
	Values map[string]string
}

func hasConflict(entries []struct {
	envName string
	value   string
}) bool {
	if len(entries) < 2 {
		return false
	}
	first := entries[0].value
	for _, e := range entries[1:] {
		if e.value != first {
			return true
		}
	}
	return false
}

func sortedKeys(m map[string][]struct {
	envName string
	value   string
}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
