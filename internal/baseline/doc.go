// Package baseline provides save, load, and diff operations for envdiff
// comparison snapshots.
//
// A baseline is a JSON file that records a previous run's comparison results.
// Subsequent runs can load the baseline and compute a delta — only the entries
// whose status has changed or that are entirely new are returned, making it
// easy to detect regressions in environment configuration over time.
//
// Usage:
//
//	// Save current results as the new baseline.
//	if err := baseline.Save("baseline.json", results); err != nil { ... }
//
//	// On the next run, load and diff.
//	b, err := baseline.Load("baseline.json")
//	delta := baseline.Diff(b, newResults)
package baseline
