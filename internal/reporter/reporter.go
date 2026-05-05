// Package reporter formats and writes comparison results to an io.Writer.
package reporter

import (
	"fmt"
	"io"
	"sort"

	"github.com/yourorg/envdiff/internal/comparator"
)

// TextReporter writes a plain-text diff report.
type TextReporter struct {
	w io.Writer
}

// NewTextReporter creates a TextReporter that writes to w.
func NewTextReporter(w io.Writer) *TextReporter {
	return &TextReporter{w: w}
}

// Write renders the comparison result to the underlying writer.
func (r *TextReporter) Write(result comparator.Result) error {
	if len(result.MissingIn) == 0 && len(result.Mismatched) == 0 {
		_, err := fmt.Fprintln(r.w, "✓ No differences found.")
		return err
	}

	// Sort environment names for deterministic output.
	envNames := make([]string, 0, len(result.MissingIn))
	for env := range result.MissingIn {
		if len(result.MissingIn[env]) > 0 {
			envNames = append(envNames, env)
		}
	}
	sort.Strings(envNames)

	for _, env := range envNames {
		keys := result.MissingIn[env]
		sort.Strings(keys)
		for _, key := range keys {
			if _, err := fmt.Fprintf(r.w, "MISSING  [%s] %s\n", env, key); err != nil {
				return err
			}
		}
	}

	sort.Slice(result.Mismatched, func(i, j int) bool {
		return result.Mismatched[i].Key < result.Mismatched[j].Key
	})

	for _, mm := range result.Mismatched {
		if _, err := fmt.Fprintf(r.w, "MISMATCH %s\n", mm.Key); err != nil {
			return err
		}
		envs := make([]string, 0, len(mm.Values))
		for e := range mm.Values {
			envs = append(envs, e)
		}
		sort.Strings(envs)
		for _, e := range envs {
			if _, err := fmt.Fprintf(r.w, "         [%s] = %q\n", e, mm.Values[e]); err != nil {
				return err
			}
		}
	}

	return nil
}
