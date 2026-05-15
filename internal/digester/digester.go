// Package digester computes a deterministic fingerprint (digest) for a set
// of comparator results. The digest can be used to detect whether the diff
// output has changed between two runs without storing the full result set.
package digester

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/user/envdiff/internal/comparator"
)

// Digest returns a stable SHA-256 hex fingerprint of the provided results.
// Results are sorted by key then by environment name before hashing so that
// the digest is independent of the order in which entries are supplied.
func Digest(results []comparator.Result) string {
	type entry struct {
		key    string
		env    string
		status string
		value  string
	}

	var entries []entry
	for _, r := range results {
		envNames := make([]string, 0, len(r.Values))
		for env := range r.Values {
			envNames = append(envNames, env)
		}
		sort.Strings(envNames)

		for _, env := range envNames {
			entries = append(entries, entry{
				key:    r.Key,
				env:    env,
				status: string(r.Status),
				value:  r.Values[env],
			})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].key != entries[j].key {
			return entries[i].key < entries[j].key
		}
		return entries[i].env < entries[j].env
	})

	h := sha256.New()
	for _, e := range entries {
		fmt.Fprintf(h, "%s\x00%s\x00%s\x00%s\n", e.key, e.env, e.status, e.value)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Equal reports whether two slices of results produce the same digest.
func Equal(a, b []comparator.Result) bool {
	return Digest(a) == Digest(b)
}
