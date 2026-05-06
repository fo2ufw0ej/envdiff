// Package redactor provides utilities for masking sensitive environment
// variable values before they are written to reports or displayed to users.
//
// A Redactor is configured with a list of key-name patterns (e.g. "PASSWORD",
// "SECRET") and a mask string (default "***"). Any key whose name contains one
// of the patterns (case-insensitive) is considered sensitive, and its value is
// replaced with the mask.
//
// Usage:
//
//	r := redactor.New(nil, "")          // use defaults
//	safe := r.RedactMap(envMap)         // returns a new map with masked values
package redactor
