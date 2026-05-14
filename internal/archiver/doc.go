// Package archiver provides a simple file-based archive for persisting
// historical comparison results produced by envdiff.
//
// Each call to Save writes a timestamped JSON file to the configured
// directory. List returns all archived entries in chronological order,
// and Latest returns the most recent entry.
//
// Typical usage:
//
//	a, err := archiver.New(".envdiff/archive")
//	if err != nil { ... }
//
//	// save after every comparison run
//	path, err := a.Save(results, "ci-run-42")
//
//	// retrieve the most recent baseline for diffing
//	latest, err := a.Latest()
package archiver
