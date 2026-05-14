// Package deduper removes duplicate comparison results, keeping the first
// occurrence of each (key, env) pair. Duplicates can arise when multiple
// loaders contribute overlapping entries for the same environment file.
package deduper

import "github.com/yourorg/envdiff/internal/comparator"

// seen key for deduplication.
type entryKey struct {
	Key string
	Env string
}

// Dedupe returns a new slice with duplicate (Key, Env) pairs removed.
// The first occurrence of each pair is retained; subsequent ones are dropped.
// The relative order of retained entries is preserved.
func Dedupe(results []comparator.Result) []comparator.Result {
	seen := make(map[entryKey]struct{}, len(results))
	out := make([]comparator.Result, 0, len(results))

	for _, r := range results {
		// A Result may reference multiple envs via Values map.
		// We treat the whole Result as keyed by its Key plus the sorted
		// canonical set of envs it covers.
		for env := range r.Values {
			ek := entryKey{Key: r.Key, Env: env}
			if _, exists := seen[ek]; !exists {
				seen[ek] = struct{}{}
			}
		}

		// Deduplicate at the Result level: same Key + same env set.
		primaryKey := entryKey{Key: r.Key, Env: canonicalEnvs(r.Values)}
		if _, exists := seen[primaryKey]; exists {
			continue
		}
		seen[primaryKey] = struct{}{}
		out = append(out, r)
	}

	return out
}

// DedupeStrict removes any Result whose exact (Key, Values) pair has already
// appeared. Two results are considered identical when their Key matches and
// every env value is equal.
func DedupeStrict(results []comparator.Result) []comparator.Result {
	type strictKey struct {
		Key    string
		Digest string
	}

	seen := make(map[strictKey]struct{}, len(results))
	out := make([]comparator.Result, 0, len(results))

	for _, r := range results {
		sk := strictKey{Key: r.Key, Digest: valuesDigest(r.Values)}
		if _, exists := seen[sk]; exists {
			continue
		}
		seen[sk] = struct{}{}
		out = append(out, r)
	}

	return out
}

// canonicalEnvs returns a single string that encodes all env names present in
// the values map, sorted, so it can be used as a map key.
func canonicalEnvs(values map[string]string) string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sortStrings(keys)
	var b []byte
	for i, k := range keys {
		if i > 0 {
			b = append(b, ',')
		}
		b = append(b, k...)
	}
	return string(b)
}

// valuesDigest returns a deterministic string encoding of the env→value map.
func valuesDigest(values map[string]string) string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sortStrings(keys)
	var b []byte
	for _, k := range keys {
		b = append(b, k...)
		b = append(b, '=')
		b = append(b, values[k]...)
		b = append(b, ';')
	}
	return string(b)
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
