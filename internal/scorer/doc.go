// Package scorer computes a numeric health score for a set of environment
// comparison results produced by the comparator package.
//
// The score is expressed as a percentage (0–100) where 100 indicates that
// every key is present and identical across all compared environments.
// Alongside the percentage, a Score value exposes raw counts for identical,
// missing, and mismatched keys so callers can render detailed summaries.
package scorer
