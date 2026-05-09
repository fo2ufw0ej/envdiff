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
// If there are no deltas, a message indicating no changes is written instead.
func (r *SnapshotReporter) Report(oldLabel, newLabel string, deltas []snapshot.Delta) error {
	if len(deltas) == 0 {
		_, err := fmt.Fprintf(r.w, "No changes between snapshots %q and %q.\n", oldLabel, newLabel)
		return err
	}

	tw := tabwriter.NewWriter(r.w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Snapshot diff: %q \u2192 %q\n", oldLabel, newLabel)
	fmt.Fprintln(tw, "KEY\tENV\tOLD STATUS\tNEW STATUS")
	fmt.Fprintln(tw, "---\t---\t----------\t----------")

	for _, d := range deltas {
		if err := r.writeDelta(tw, d); err != nil {
			return err
		}
	}
	return tw.Flush()
}

// writeDelta formats and writes a single Delta row to the given writer.
func (r *SnapshotReporter) writeDelta(w io.Writer, d snapshot.Delta) error {
	oldStatus := string(d.Old.Status)
	if d.Added {
		oldStatus = "(new)"
	}
	_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
		d.New.Key,
		d.New.Env,
		oldStatus,
		string(d.New.Status),
	)
	return err
}
