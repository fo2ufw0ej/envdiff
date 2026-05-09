package reporter

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/aliaser"
	"github.com/user/envdiff/internal/comparator"
)

// AliasReporter wraps another Reporter and applies alias rules before
// delegating to the inner reporter. It also prints a preamble listing
// active rules so the output is self-documenting.
type AliasReporter struct {
	inner   Reporter
	aliaser *aliaser.Aliaser
	rules   []aliaser.Rule
}

// NewAliasReporter returns an AliasReporter that rewrites keys according to
// rules before delegating to inner.
func NewAliasReporter(inner Reporter, rules []aliaser.Rule) (*AliasReporter, error) {
	a, err := aliaser.New(rules)
	if err != nil {
		return nil, err
	}
	return &AliasReporter{inner: inner, aliaser: a, rules: rules}, nil
}

// Report rewrites result keys via the alias rules, prints a summary of
// active aliases, then delegates to the inner reporter.
func (r *AliasReporter) Report(w io.Writer, results []comparator.Result) error {
	if len(r.rules) > 0 {
		active := make([]aliaser.Rule, len(r.rules))
		copy(active, r.rules)
		sort.Slice(active, func(i, j int) bool {
			return active[i].From < active[j].From
		})
		fmt.Fprintln(w, "# Active key aliases:")
		for _, rule := range active {
			fmt.Fprintf(w, "#   %s -> %s\n", rule.From, rule.To)
		}
		fmt.Fprintln(w)
	}
	return r.inner.Report(w, r.aliaser.Apply(results))
}
