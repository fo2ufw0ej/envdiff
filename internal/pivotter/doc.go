// Package pivotter transforms a flat slice of comparison results into a
// key-major pivot table.
//
// Given results that look like:
//
//	{Key: "PORT", Values: {"dev": "3000", "prod": "8080"}}
//	{Key: "HOST", Values: {"dev": "localhost"}}
//
// Pivot returns a Table whose Rows are sorted by key and whose Values map
// contains an entry for every environment seen across all results.  If an
// environment did not define a particular key, its cell is set to the
// Absent sentinel string ("<absent>").
//
// This is useful for rendering tabular reports (e.g. HTML tables or CSV
// output) where each column represents one environment.
package pivotter
