// Package counter tallies comparison results by status across environments.
package counter

import (
	"sort"

	"github.com/yourorg/envdiff/internal/comparator"
)

// Counts holds the tally of result statuses for a single environment.
type Counts struct {
	Identical  int
	Missing    int
	Mismatched int
	Total      int
}

// Report maps each environment name to its Counts.
type Report map[string]Counts

// EnvNames returns the environment names in the report, sorted alphabetically.
func (r Report) EnvNames() []string {
	names := make([]string, 0, len(r))
	for k := range r {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// Count tallies comparator results per environment.
// Each result may reference multiple environments; every env seen in Values
// is counted independently.
func Count(results []comparator.Result) Report {
	report := Report{}

	for _, res := range results {
		for env := range res.Values {
			c := report[env]
			c.Total++
			switch res.Status {
			case comparator.StatusIdentical:
				c.Identical++
			case comparator.StatusMissing:
				c.Missing++
			case comparator.StatusMismatched:
				c.Mismatched++
			}
			report[env] = c
		}
	}

	return report
}
