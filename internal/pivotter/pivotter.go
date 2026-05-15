// Package pivotter re-shapes a flat []comparator.Result slice into a
// key-major table: for every unique key, one row that lists every
// environment's value (or a sentinel when the key is absent).
package pivotter

import (
	"sort"

	"github.com/user/envdiff/internal/comparator"
)

// Absent is the sentinel placed in a cell when an environment does not
// define the key at all.
const Absent = "<absent>"

// Row holds one key and its per-environment values.
type Row struct {
	Key    string
	Values map[string]string // env name → value (or Absent)
}

// Table is the result of pivoting: an ordered slice of Rows plus the
// full sorted list of environment names seen across all results.
type Table struct {
	Envs []string
	Rows []Row
}

// Pivot converts a flat result slice into a key-major Table.
func Pivot(results []comparator.Result) Table {
	envSet := map[string]struct{}{}
	for _, r := range results {
		for env := range r.Values {
			envSet[env] = struct{}{}
		}
	}

	envs := make([]string, 0, len(envSet))
	for e := range envSet {
		envs = append(envs, e)
	}
	sort.Strings(envs)

	// Collect unique keys preserving insertion order, then sort.
	keySet := map[string]struct{}{}
	for _, r := range results {
		keySet[r.Key] = struct{}{}
	}
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build a lookup: key → (env → value).
	lookup := make(map[string]map[string]string, len(keys))
	for _, r := range results {
		if lookup[r.Key] == nil {
			lookup[r.Key] = make(map[string]string)
		}
		for env, val := range r.Values {
			lookup[r.Key][env] = val
		}
	}

	rows := make([]Row, 0, len(keys))
	for _, k := range keys {
		cells := make(map[string]string, len(envs))
		for _, e := range envs {
			v, ok := lookup[k][e]
			if !ok {
				v = Absent
			}
			cells[e] = v
		}
		rows = append(rows, Row{Key: k, Values: cells})
	}

	return Table{Envs: envs, Rows: rows}
}
