package reporter

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/your-org/envdiff/internal/aggregator"
	"github.com/your-org/envdiff/internal/comparator"
)

// AggregatorReporter writes a tabular summary produced by the aggregator package.
type AggregatorReporter struct {
	w io.Writer
}

// NewAggregatorReporter returns a reporter that writes an aggregated summary
// table to w.
func NewAggregatorReporter(w io.Writer) *AggregatorReporter {
	return &AggregatorReporter{w: w}
}

// Report aggregates results and writes a summary table.
func (r *AggregatorReporter) Report(results []comparator.Result) error {
	rep := aggregator.Aggregate(results)

	tw := tabwriter.NewWriter(r.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ENV\tTOTAL\tIDENTICAL\tMISMATCHED\tMISSING")
	fmt.Fprintln(tw, "---\t-----\t---------\t----------\t-------")
	for _, e := range rep.Envs {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%d\n",
			e.Env, e.Total, e.Identical, e.Mismatched, e.Missing)
	}
	if err := tw.Flush(); err != nil {
		return err
	}
	fmt.Fprintf(r.w, "\nTotal unique keys: %d  |  Total issues: %d\n",
		rep.TotalKeys, rep.TotalIssues)
	return nil
}
