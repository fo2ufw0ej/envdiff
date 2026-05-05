// Package comparator provides functionality for comparing parsed .env files
// across multiple environments and identifying discrepancies.
package comparator

// Result holds the outcome of comparing two or more .env files.
type Result struct {
	// MissingIn maps an environment name to keys missing from that environment.
	MissingIn map[string][]string
	// Mismatched contains keys whose values differ across environments.
	Mismatched []MismatchedKey
}

// MismatchedKey describes a key that exists in all environments but has
// differing values.
type MismatchedKey struct {
	Key    string
	Values map[string]string // env name -> value
}

// Compare takes a map of environment name -> parsed key/value pairs and
// returns a Result describing missing and mismatched keys.
func Compare(envs map[string]map[string]string) Result {
	result := Result{
		MissingIn: make(map[string][]string),
	}

	// Collect the union of all keys.
	allKeys := unionKeys(envs)

	for key := range allKeys {
		values := make(map[string]string)
		presentIn := []string{}

		for envName, pairs := range envs {
			if val, ok := pairs[key]; ok {
				values[envName] = val
				presentIn = append(presentIn, envName)
			} else {
				result.MissingIn[envName] = append(result.MissingIn[envName], key)
			}
		}

		if len(presentIn) < 2 {
			continue
		}

		if hasMismatch(values) {
			result.Mismatched = append(result.Mismatched, MismatchedKey{
				Key:    key,
				Values: values,
			})
		}
	}

	return result
}

func unionKeys(envs map[string]map[string]string) map[string]struct{} {
	keys := make(map[string]struct{})
	for _, pairs := range envs {
		for k := range pairs {
			keys[k] = struct{}{}
		}
	}
	return keys
}

func hasMismatch(values map[string]string) bool {
	var first string
	set := false
	for _, v := range values {
		if !set {
			first = v
			set = true
			continue
		}
		if v != first {
			return true
		}
	}
	return false
}
