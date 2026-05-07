package reporter

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/linter"
)

// LintReporter writes lint findings in a human-readable text format.
type LintReporter struct {
	w io.Writer
}

// NewLintReporter creates a LintReporter that writes to w.
func NewLintReporter(w io.Writer) *LintReporter {
	return &LintReporter{w: w}
}

// Write outputs all findings grouped by severity, sorted by key.
func (r *LintReporter) Write(env string, findings []linter.Finding) error {
	if len(findings) == 0 {
		_, err := fmt.Fprintf(r.w, "[%s] no lint findings\n", env)
		return err
	}

	sorted := make([]linter.Finding, len(findings))
	copy(sorted, findings)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Severity != sorted[j].Severity {
			// errors before warnings
			return sorted[i].Severity == linter.SeverityError
		}
		return sorted[i].Key < sorted[j].Key
	})

	_, err := fmt.Fprintf(r.w, "[%s] %d lint finding(s):\n", env, len(sorted))
	if err != nil {
		return err
	}
	for _, f := range sorted {
		_, err = fmt.Fprintf(r.w, "  %-7s %s: %s\n", f.Severity, f.Key, f.Message)
		if err != nil {
			return err
		}
	}
	return nil
}
