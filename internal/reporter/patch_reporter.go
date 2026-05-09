package reporter

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/yourorg/envdiff/internal/patcher"
)

// PatchReporter writes a human-readable summary of patch results.
type PatchReporter struct {
	w io.Writer
}

// NewPatchReporter returns a PatchReporter that writes to w.
func NewPatchReporter(w io.Writer) *PatchReporter {
	return &PatchReporter{w: w}
}

// Write formats results as an aligned table and writes it to the underlying
// writer. It returns the first error encountered, if any.
func (r *PatchReporter) Write(results []patcher.Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(r.w, "No patch actions recorded.")
		return err
	}

	tw := tabwriter.NewWriter(r.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "KEY\tACTION\tOLD\tNEW")
	fmt.Fprintln(tw, "---\t------\t---\t---")

	for _, res := range results {
		old := res.OldVal
		if old == "" && res.Action == "added" {
			old = "<absent>"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", res.Key, res.Action, old, res.NewVal)
	}

	return tw.Flush()
}
