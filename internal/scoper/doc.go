// Package scoper provides environment-level scoping for comparison results.
//
// When working with many environments it is often useful to narrow the output
// to a specific subset — for example, comparing only "dev" and "staging"
// while ignoring "prod". Scoper accepts a list of environment names and
// strips all other environments from each comparator.Result, dropping any
// result that has no remaining data after filtering.
//
// Usage:
//
//	s := scoper.New([]string{"dev", "staging"})
//	narrowed := s.Apply(results)
package scoper
