// Package scoper restricts comparison results to a defined set of
// environment names, discarding any envs not in the allowed scope.
package scoper

import "github.com/yourorg/envdiff/internal/comparator"

// Scoper filters comparator results to only include data from
// environments whose names appear in the allowed set.
type Scoper struct {
	allowed map[string]struct{}
}

// New creates a Scoper that keeps only the envs listed in envNames.
// An empty envNames list produces a no-op scoper (all envs kept).
func New(envNames []string) *Scoper {
	allowed := make(map[string]struct{}, len(envNames))
	for _, n := range envNames {
		allowed[n] = struct{}{}
	}
	return &Scoper{allowed: allowed}
}

// Apply returns a new slice of results where each result's Values map
// contains only the scoped environments. Results that become empty
// (no matching envs remain) are dropped entirely.
func (s *Scoper) Apply(results []comparator.Result) []comparator.Result {
	if len(s.allowed) == 0 {
		return results
	}

	out := make([]comparator.Result, 0, len(results))
	for _, r := range results {
		scoped := make(map[string]string)
		for env, val := range r.Values {
			if _, ok := s.allowed[env]; ok {
				scoped[env] = val
			}
		}
		if len(scoped) == 0 {
			continue
		}
		out = append(out, comparator.Result{
			Key:    r.Key,
			Status: r.Status,
			Values: scoped,
		})
	}
	return out
}

// Envs returns the sorted list of allowed environment names.
func (s *Scoper) Envs() []string {
	if len(s.allowed) == 0 {
		return nil
	}
	names := make([]string, 0, len(s.allowed))
	for n := range s.allowed {
		names = append(names, n)
	}
	sortStrings(names)
	return names
}

func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
