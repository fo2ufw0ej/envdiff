// Package profiler analyses a set of parsed env maps and produces
// statistics about key coverage, value lengths, and uniqueness.
package profiler

import "sort"

// KeyStat holds per-key statistics across all environments.
type KeyStat struct {
	Key          string
	EnvsPresent  int
	EnvsTotal    int
	UniqueValues int
	MaxValueLen  int
}

// Report is the output of Profile.
type Report struct {
	TotalKeys  int
	TotalEnvs  int
	KeyStats   []KeyStat
}

// Profile computes coverage and value statistics for the provided env maps.
// The map key is the environment name; the value is the parsed key→value map.
func Profile(envs map[string]map[string]string) Report {
	if len(envs) == 0 {
		return Report{}
	}

	// Collect union of all keys.
	keySet := make(map[string]struct{})
	for _, kv := range envs {
		for k := range kv {
			keySet[k] = struct{}{}
		}
	}

	totalEnvs := len(envs)
	stats := make([]KeyStat, 0, len(keySet))

	for key := range keySet {
		valueSet := make(map[string]struct{})
		present := 0
		maxLen := 0

		for _, kv := range envs {
			v, ok := kv[key]
			if !ok {
				continue
			}
			present++
			valueSet[v] = struct{}{}
			if len(v) > maxLen {
				maxLen = len(v)
			}
		}

		stats = append(stats, KeyStat{
			Key:          key,
			EnvsPresent:  present,
			EnvsTotal:    totalEnvs,
			UniqueValues: len(valueSet),
			MaxValueLen:  maxLen,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Key < stats[j].Key
	})

	return Report{
		TotalKeys: len(keySet),
		TotalEnvs: totalEnvs,
		KeyStats:  stats,
	}
}
