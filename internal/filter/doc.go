// Package filter provides utilities for narrowing down comparison results
// produced by the comparator package.
//
// Callers can restrict output to entries that are missing in one or more
// environments, entries whose values differ across environments, or entries
// whose key matches a given prefix. Filters can be combined freely.
package filter
