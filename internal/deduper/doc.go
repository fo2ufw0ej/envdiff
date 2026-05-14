// Package deduper provides utilities for removing duplicate entries from a
// slice of comparator.Result values.
//
// Two deduplication strategies are available:
//
//   - Dedupe: removes results whose (Key, env-set) pair has already been seen.
//     This is useful when the same key appears for the same set of environments
//     more than once due to overlapping loader output.
//
//   - DedupeStrict: removes results whose Key and complete Values map are
//     identical to a previously seen result. Results with the same key but
//     differing values are both retained.
//
// In both cases the relative order of retained entries is preserved and the
// input slice is never mutated.
package deduper
