package reporter

import (
	"fmt"
	"io"

	"github.com/user/envdiff/internal/renamer"
)

// NewRenameReporter returns a Reporter that summarises the outcome of rename
// rule application produced by renamer.Apply.
func NewRenameReporter(w io.Writer, ruleResults []renamer.Result) Reporter {
	return &renameReporter{w: w, ruleResults: ruleResults}
}

type renameReporter struct {
	w           io.Writer
	ruleResults []renamer.Result
}

func (r *renameReporter) Report(results interface{}) error {
	if len(r.ruleResults) == 0 {
		fmt.Fprintln(r.w, "No rename rules defined.")
		return nil
	}

	matched := 0
	for _, rr := range r.ruleResults {
		status := "unmatched"
		if rr.Matched {
			status = "matched"
			matched++
		}
		fmt.Fprintf(r.w, "  [%s] %s -> %s\n", status, rr.Rule.OldKey, rr.Rule.NewKey)
	}
	fmt.Fprintf(r.w, "\n%d/%d rules matched.\n", matched, len(r.ruleResults))
	return nil
}
