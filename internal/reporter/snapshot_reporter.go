package reporter

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/user/envdiff/internal/snapshot"
)

// NewSnapshotReporter returns a Reporter that writes snapshot deltas in a
// human-readable tabular format.
func NewSnapshotReporter(w io.Writer) *SnapshotReporter {
	return &SnapshotReporter{w: w}
}

// SnapshotReporter writes the differences between two snapshots.
type SnapshotReporter struct {
	w io.Writer
}

// Report writes a formatted diff of the provided deltas.
func (r *SnapshotReporter) Report(oldLabel, newLabel string, deltas []snapshot.Delta) error {
	if len(deltas) == 0 {
		_, err := fmt.Fprintf(r.w, "No changes between snapshots %q and %q.\n", oldLabel, newLabel)
		return err
	}

	tw := tabwriter.NewWriter(r.w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Snapshot diff: %q → %q\n", oldLabel, newLabel)
	fmt.Fprintln(tw, "KEY\tENV\tOLD STATUS\tNEW STATUS")
	fmt.Fprintln(tw, "---\t---\t----------\t----------")

	for _, d := range deltas {
		oldStatus := string(d.Old.Status)
		if d.Added {
			oldStatus = "(new)"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			d.New.Key,
			d.New.Env,
			oldStatus,
			string(d.New.Status),
		)
	}
	return tw.Flush()
}
