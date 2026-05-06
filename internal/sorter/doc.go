// Package sorter provides utilities for ordering comparator.Result slices
// produced by the envdiff comparison engine.
//
// Three sort strategies are available:
//
//   - SortByKey   – alphabetical order by environment variable name (default)
//   - SortByStatus – missing keys appear before mismatched keys
//   - SortByEnv   – grouped by the first environment that references each key
//
// Sort never modifies the input slice; it always returns a new slice.
package sorter
