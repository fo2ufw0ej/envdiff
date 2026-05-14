package aggregator

import (
	"sort"

	"github.com/your-org/envdiff/internal/comparator"
)

// EnvSummary holds aggregated statistics for a single environment.
type EnvSummary struct {
	Env        string
	Total      int
	Missing    int
	Mismatched int
	Identical  int
}

// Report is the result of aggregating comparator results across all environments.
type Report struct {
	Envs       []EnvSummary
	TotalKeys  int
	TotalIssues int
}

// Aggregate summarises comparator.Result entries per environment.
func Aggregate(results []comparator.Result) Report {
	counts := map[string]*EnvSummary{}
	keySet := map[string]struct{}{}

	for _, r := range results {
		keySet[r.Key] = struct{}{}
		for env := range r.Values {
			if _, ok := counts[env]; !ok {
				counts[env] = &EnvSummary{Env: env}
			}
			s := counts[env]
			s.Total++
			switch r.Status {
			case comparator.StatusMissing:
				if _, present := r.Values[env]; !present {
					s.Missing++
				} else {
					s.Identical++
				}
			case comparator.StatusMismatch:
				s.Mismatched++
			default:
				s.Identical++
			}
		}
	}

	envs := make([]EnvSummary, 0, len(counts))
	totalIssues := 0
	for _, s := range counts {
		totalIssues += s.Missing + s.Mismatched
		envs = append(envs, *s)
	}
	sort.Slice(envs, func(i, j int) bool { return envs[i].Env < envs[j].Env })

	return Report{
		Envs:        envs,
		TotalKeys:   len(keySet),
		TotalIssues: totalIssues,
	}
}
