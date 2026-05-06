// Package summary provides aggregation of comparison results into
// high-level statistics useful for reporting and exit-code decisions.
package summary

import "github.com/user/envdiff/internal/comparator"

// Stats holds aggregated counts derived from a comparison result set.
type Stats struct {
	TotalKeys    int
	MissingKeys  int
	Mismatched   int
	Identical    int
	EnvCount     int
	AffectedEnvs []string
}

// HasDifferences reports whether any missing or mismatched keys were found.
func (s Stats) HasDifferences() bool {
	return s.MissingKeys > 0 || s.Mismatched > 0
}

// Compute derives Stats from a slice of comparator.Result entries.
func Compute(results []comparator.Result) Stats {
	envSet := make(map[string]struct{})
	affectedSet := make(map[string]struct{})

	var missing, mismatched, identical int

	for _, r := range results {
		for env := range r.Values {
			envSet[env] = struct{}{}
		}

		switch r.Status {
		case comparator.StatusMissing:
			missing++
			for env := range r.Values {
				affectedSet[env] = struct{}{}
			}
			for _, env := range r.MissingIn {
				affectedSet[env] = struct{}{}
			}
		case comparator.StatusMismatch:
			mismatched++
			for env := range r.Values {
				affectedSet[env] = struct{}{}
			}
		case comparator.StatusOK:
			identical++
		}
	}

	affected := make([]string, 0, len(affectedSet))
	for env := range affectedSet {
		affected = append(affected, env)
	}

	return Stats{
		TotalKeys:    len(results),
		MissingKeys:  missing,
		Mismatched:   mismatched,
		Identical:    identical,
		EnvCount:     len(envSet),
		AffectedEnvs: affected,
	}
}
