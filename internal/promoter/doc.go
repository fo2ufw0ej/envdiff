// Package promoter copies key-value pairs from one environment map into
// another. It supports two conflict-resolution strategies:
//
//   - SkipExisting: keys already present in the destination are left unchanged.
//   - OverwriteExisting: destination values are replaced with the source value.
//
// Promote returns a new map (the destination is never mutated) together with
// a slice of Result values that describe what happened to each key.
package promoter
