package reporter

import (
	"fmt"
	"io"
	"sort"

	"github.com/yourorg/envdiff/internal/classifier"
	"github.com/yourorg/envdiff/internal/comparator"
)

// classifierReporter groups results by their "category" tag and prints a
// per-category summary followed by individual key lines.
type classifierReporter struct {
	w io.Writer
	c *classifier.Classifier
}

// NewClassifierReporter returns a Reporter that groups output by semantic category.
func NewClassifierReporter(w io.Writer, c *classifier.Classifier) Reporter {
	if c == nil {
		c = classifier.New()
	}
	return &classifierReporter{w: w, c: c}
}

func (r *classifierReporter) Report(results []comparator.Result) error {
	annotated := r.c.ClassifyResults(results)

	groups := make(map[string][]comparator.Result)
	for _, res := range annotated {
		cat := res.Tags["category"]
		groups[cat] = append(groups[cat], res)
	}

	cats := make([]string, 0, len(groups))
	for cat := range groups {
		cats = append(cats, cat)
	}
	sort.Strings(cats)

	for _, cat := range cats {
		fmt.Fprintf(r.w, "[%s]\n", cat)
		entries := groups[cat]
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Key < entries[j].Key
		})
		for _, e := range entries {
			switch e.Status {
			case comparator.StatusIdentical:
				fmt.Fprintf(r.w, "  OK       %s\n", e.Key)
			case comparator.StatusMissing:
				fmt.Fprintf(r.w, "  MISSING  %s\n", e.Key)
			case comparator.StatusMismatch:
				fmt.Fprintf(r.w, "  MISMATCH %s\n", e.Key)
			default:
				fmt.Fprintf(r.w, "  %-8s %s\n", e.Status, e.Key)
			}
		}
	}
	return nil
}
