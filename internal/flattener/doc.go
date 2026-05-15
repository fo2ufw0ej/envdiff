// Package flattener converts a slice of comparator.Result values into a flat
// list of Entry structs, each representing a single (key, environment) cell.
//
// This is useful when downstream consumers need a row-oriented view of the
// comparison matrix rather than the key-centric Result slice returned by the
// comparator.
//
// Usage:
//
//	entries := flattener.Flatten(results)
//	idx     := flattener.Index(entries)   // optional O(1) lookup by "KEY@ENV"
package flattener
