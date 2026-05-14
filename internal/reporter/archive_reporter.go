package reporter

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/user/envdiff/internal/archiver"
)

// ArchiveReporter writes a human-readable summary of archived entries.
type ArchiveReporter struct {
	w io.Writer
}

// NewArchiveReporter returns a reporter that writes to w.
func NewArchiveReporter(w io.Writer) *ArchiveReporter {
	return &ArchiveReporter{w: w}
}

// Report prints a table of all entries in the archive.
func (r *ArchiveReporter) Report(entries []archiver.Entry) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(r.w, "archive: no entries found")
		return err
	}

	tw := tabwriter.NewWriter(r.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "#\tTIMESTAMP\tLABEL\tRESULTS")
	fmt.Fprintln(tw, "-\t---------\t-----\t-------")
	for i, e := range entries {
		label := e.Label
		if label == "" {
			label = "(unlabelled)"
		}
		fmt.Fprintf(tw, "%d\t%s\t%s\t%d\n",
			i+1,
			e.Timestamp.Format(time.RFC3339),
			label,
			len(e.Results),
		)
	}
	return tw.Flush()
}
