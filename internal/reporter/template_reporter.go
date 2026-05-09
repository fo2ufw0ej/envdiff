package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/templater"
)

// TemplateReporter writes rendered template results to an io.Writer.
type TemplateReporter struct {
	w io.Writer
}

// NewTemplateReporter returns a TemplateReporter that writes to w.
func NewTemplateReporter(w io.Writer) *TemplateReporter {
	return &TemplateReporter{w: w}
}

// Write outputs the template results for all environments.
// Each environment section is separated by a blank line.
func (r *TemplateReporter) Write(results []templater.Result) error {
	for i, res := range results {
		if i > 0 {
			if _, err := fmt.Fprintln(r.w); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(r.w, "# %s\n", res.EnvName); err != nil {
			return err
		}
		for _, line := range res.Lines {
			if _, err := fmt.Fprintln(r.w, line); err != nil {
				return err
			}
		}
		if len(res.Missing) > 0 {
			if _, err := fmt.Fprintf(r.w, "# missing: %s\n", strings.Join(res.Missing, ", ")); err != nil {
				return err
			}
		}
	}
	return nil
}
