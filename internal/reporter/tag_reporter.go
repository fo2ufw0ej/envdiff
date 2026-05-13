package reporter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/tagger"
)

// tagReporter writes a human-readable summary of tagged results, grouping
// output by tag name.
type tagReporter struct {
	w io.Writer
}

// NewTagReporter returns a Reporter that renders tagged results grouped by
// tag. Results with no tags are printed under "(untagged)".
func NewTagReporter(w io.Writer) *tagReporter {
	return &tagReporter{w: w}
}

// Report writes tagged results to the underlying writer.
func (r *tagReporter) Report(results []tagger.TaggedResult) error {
	groups := map[string][]tagger.TaggedResult{}
	for _, res := range results {
		if len(res.Tags) == 0 {
			groups["(untagged)"] = append(groups["(untagged)"], res)
			continue
		}
		for _, tag := range res.Tags {
			groups[tag] = append(groups[tag], res)
		}
	}

	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, tag := range keys {
		fmt.Fprintf(r.w, "[%s]\n", tag)
		entries := groups[tag]
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Key < entries[j].Key
		})
		for _, e := range entries {
			fmt.Fprintf(r.w, "  %-30s %s\n", e.Key, strings.ToUpper(string(e.Status)))
		}
	}
	return nil
}
