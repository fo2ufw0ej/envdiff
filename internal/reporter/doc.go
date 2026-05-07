// Package reporter provides output formatters for envdiff comparison results.
//
// Three reporter implementations are available:
//
//   - TextReporter     – human-readable plain-text output (default)
//   - JSONReporter     – machine-readable JSON output
//   - MarkdownReporter – Markdown table output suitable for PR comments
//
// All reporters implement the Reporter interface:
//
//	type Reporter interface {
//	    Report(results []comparator.Result) error
//	}
//
// Choose a reporter based on the desired output format and pass it an
// io.Writer (typically os.Stdout) at construction time.
//
// Example usage:
//
//	w := os.Stdout
//	r := reporter.NewTextReporter(w)
//	if err := r.Report(results); err != nil {
//	    log.Fatal(err)
//	}
//
// For CI pipelines or tooling that consumes structured data, prefer
// JSONReporter. For pull-request comment automation, MarkdownReporter
// produces a table that renders natively on GitHub and GitLab.
package reporter
