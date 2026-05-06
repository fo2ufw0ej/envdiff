// Package reporter provides output formatters for envdiff comparison results.
//
// Three reporter implementations are available:
//
//   - TextReporter  – human-readable plain-text output (default)
//   - JSONReporter  – machine-readable JSON output
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
package reporter
