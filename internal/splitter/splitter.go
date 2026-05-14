// Package splitter splits a flat list of comparator results into
// per-environment maps, making it easy to produce per-env reports or
// drive further per-env processing pipelines.
package splitter

import "github.com/yourorg/envdiff/internal/comparator"

// EnvResults maps an environment name to the subset of results that
// involve that environment (either as a missing key or as one side of
// a mismatch).
type EnvResults map[string][]comparator.Result

// Split partitions results by environment name.  A result that covers
// multiple environments (e.g. a mismatch present in several envs) is
// placed under every relevant environment key.
//
// Results with no environment information are collected under the
// special key "_global".
func Split(results []comparator.Result) EnvResults {
	out := make(EnvResults)

	for _, r := range results {
		envs := envNames(r)
		if len(envs) == 0 {
			out["_global"] = append(out["_global"], r)
			continue
		}
		for _, env := range envs {
			out[env] = append(out[env], r)
		}
	}

	return out
}

// Envs returns a sorted slice of all environment names present in
// results, excluding the synthetic "_global" bucket.
func Envs(er EnvResults) []string {
	keys := make([]string, 0, len(er))
	for k := range er {
		if k == "_global" {
			continue
		}
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys
}

// envNames extracts all environment names referenced by a single result.
func envNames(r comparator.Result) []string {
	seen := make(map[string]struct{})
	var out []string
	for env := range r.Values {
		if _, ok := seen[env]; !ok {
			seen[env] = struct{}{}
			out = append(out, env)
		}
	}
	sortStrings(out)
	return out
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
