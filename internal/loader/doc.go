// Package loader handles discovery and loading of one or more .env files,
// returning a unified EnvMap that downstream packages (comparator, reporter)
// can operate on.
//
// # Usage
//
// Load a specific list of files:
//
//	em, err := loader.LoadFiles([]string{".env", ".env.production"})
//
// Scan a directory for all .env* files:
//
//	em, err := loader.LoadDir("/path/to/project")
//
// The resulting EnvMap keys are short environment names derived from the
// file names (e.g. ".env.staging" → "staging", ".env" → "default").
package loader
