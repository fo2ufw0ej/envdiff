package reporter

import (
	"fmt"
	"html"
	"io"
	"sort"

	"github.com/user/envdiff/internal/comparator"
)

// NewHTMLReporter returns a Reporter that writes an HTML table of differences.
func NewHTMLReporter(w io.Writer) func([]comparator.Diff) error {
	return func(diffs []comparator.Diff) error {
		fmt.Fprintln(w, "<!DOCTYPE html>")
		fmt.Fprintln(w, "<html><head><meta charset=\"UTF-8\">")
		fmt.Fprintln(w, "<title>envdiff report</title>")
		fmt.Fprintln(w, "<style>")
		fmt.Fprintln(w, "  body { font-family: sans-serif; padding: 1rem; }")
		fmt.Fprintln(w, "  table { border-collapse: collapse; width: 100%; }")
		fmt.Fprintln(w, "  th, td { border: 1px solid #ccc; padding: 0.4rem 0.8rem; text-align: left; }")
		fmt.Fprintln(w, "  th { background: #f0f0f0; }")
		fmt.Fprintln(w, "  .missing { background: #fff3cd; }")
		fmt.Fprintln(w, "  .mismatch { background: #f8d7da; }")
		fmt.Fprintln(w, "</style></head><body>")
		fmt.Fprintln(w, "<h1>envdiff Report</h1>")

		if len(diffs) == 0 {
			fmt.Fprintln(w, "<p>No differences found.</p>")
			fmt.Fprintln(w, "</body></html>")
			return nil
		}

		// Collect all env names for column headers.
		envSet := map[string]struct{}{}
		for _, d := range diffs {
			for env := range d.Values {
				envSet[env] = struct{}{}
			}
		}
		envNames := make([]string, 0, len(envSet))
		for e := range envSet {
			envNames = append(envNames, e)
		}
		sort.Strings(envNames)

		fmt.Fprintln(w, "<table>")
		fmt.Fprint(w, "  <tr><th>Key</th><th>Type</th>")
		for _, e := range envNames {
			fmt.Fprintf(w, "<th>%s</th>", html.EscapeString(e))
		}
		fmt.Fprintln(w, "</tr>")

		sort.Slice(diffs, func(i, j int) bool { return diffs[i].Key < diffs[j].Key })

		for _, d := range diffs {
			rowClass := "mismatch"
			if d.Type == comparator.Missing {
				rowClass = "missing"
			}
			fmt.Fprintf(w, "  <tr class=\"%s\"><td>%s</td><td>%s</td>",
				rowClass,
				html.EscapeString(d.Key),
				html.EscapeString(string(d.Type)),
			)
			for _, e := range envNames {
				val, ok := d.Values[e]
				if !ok {
					fmt.Fprint(w, "<td><em>absent</em></td>")
				} else {
					fmt.Fprintf(w, "<td>%s</td>", html.EscapeString(val))
				}
			}
			fmt.Fprintln(w, "</tr>")
		}

		fmt.Fprintln(w, "</table>")
		fmt.Fprintln(w, "</body></html>")
		return nil
	}
}
