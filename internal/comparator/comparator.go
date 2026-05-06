package comparator

import "sort"

// DiffEntry represents a single key and its values across all compared environments.
type DiffEntry struct {
	// Key is the environment variable name.
	Key string
	// Values maps each environment name to the value found for Key.
	// An empty string indicates the key is absent in that environment.
	Values map[string]string
}

// Result holds the full output of a comparison across multiple environments.
type Result struct {
	// Envs is the ordered list of environment names that were compared.
	Envs []string
	// Entries contains one DiffEntry per unique key found across all environments.
	Entries []DiffEntry
}

// Compare accepts a map of environment name → key/value pairs and returns a
// Result describing every key found across all environments.
func Compare(envs map[string]map[string]string) Result {
	if len(envs) == 0 {
		return Result{}
	}

	keys := unionKeys(envs)
	envNames := make([]string, 0, len(envs))
	for name := range envs {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)

	entries := make([]DiffEntry, 0, len(keys))
	for _, key := range keys {
		values := make(map[string]string, len(envNames))
		for _, name := range envNames {
			values[name] = envs[name][key]
		}
		entries = append(entries, DiffEntry{Key: key, Values: values})
	}

	return Result{Envs: envNames, Entries: entries}
}

// unionKeys returns a sorted slice of all unique keys across every environment.
func unionKeys(envs map[string]map[string]string) []string {
	seen := make(map[string]struct{})
	for _, kv := range envs {
		for k := range kv {
			seen[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// hasMismatch reports whether a DiffEntry has differing values across environments.
func hasMismatch(e DiffEntry) bool {
	var first *string
	for _, v := range e.Values {
		if first == nil {
			copy := v
			first = &copy
			continue
		}
		if v != *first {
			return true
		}
	}
	return false
}
