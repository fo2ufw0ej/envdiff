package reporter

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/comparator"
)

// markdownReporter formats comparison results as a Markdown table.
type markdownReporter struct {
	w io.Writer
}

// NewMarkdownReporter returns a Reporter that writes Markdown-formatted output
// to w. The output includes a summary table suitable for inclusion in
// documentation or pull-request comments.
func NewMarkdownReporter(w io.Writer) Reporter {
	return &markdownReporter{w: w}
}

func (r *markdownReporter) Report(results []comparator.Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(r.w, "_No differences found._")
		return err
	}

	// Collect all environment names for the header row.
	envSet := map[string]struct{}{}
	for _, res := range results {
		for env := range res.Values {
			envSet[env] = struct{}{}
		}
	}
	envs := sortedKeys(envSet)

	// Header
	header := "| Key |"
	separator := "| --- |"
	for _, env := range envs {
		header += fmt.Sprintf(" %s |", env)
		separator += " --- |"
	}
	fmt.Fprintln(r.w, header)
	fmt.Fprintln(r.w, separator)

	// Sort results by key for deterministic output.
	sorted := make([]comparator.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, res := range sorted {
		row := fmt.Sprintf("| `%s` |", res.Key)
		for _, env := range envs {
			val, ok := res.Values[env]
			if !ok {
				row += " _(missing)_ |"
			} else {
				row += fmt.Sprintf(" `%s` |", val)
			}
		}
		fmt.Fprintln(r.w, row)
	}

	return nil
}
