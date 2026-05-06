// Package validator checks parsed env maps for common issues such as
// empty values, duplicate keys (case-insensitive), and keys that contain
// invalid characters.
package validator

import (
	"fmt"
	"strings"
	"unicode"
)

// Issue represents a single validation problem found in an env map.
type Issue struct {
	Key     string
	Message string
}

func (i Issue) String() string {
	return fmt.Sprintf("key %q: %s", i.Key, i.Message)
}

// Validate inspects the provided env map and returns a (possibly empty)
// slice of Issues. It checks for:
//   - keys that are empty strings
//   - keys that contain whitespace or non-printable characters
//   - case-insensitive duplicate keys
//   - values that are entirely whitespace (suspicious but not fatal)
func Validate(env map[string]string) []Issue {
	var issues []Issue
	seen := make(map[string]string) // lower-case key -> original key

	for k, v := range env {
		if k == "" {
			issues = append(issues, Issue{Key: k, Message: "key must not be empty"})
			continue
		}

		if err := validateKeyChars(k); err != nil {
			issues = append(issues, Issue{Key: k, Message: err.Error()})
		}

		lower := strings.ToLower(k)
		if orig, dup := seen[lower]; dup {
			issues = append(issues, Issue{
				Key:     k,
				Message: fmt.Sprintf("case-insensitive duplicate of key %q", orig),
			})
		} else {
			seen[lower] = k
		}

		if v != strings.TrimSpace(v) && strings.TrimSpace(v) == "" {
			issues = append(issues, Issue{Key: k, Message: "value is blank (only whitespace)"})
		}
	}

	return issues
}

// validateKeyChars returns an error if the key contains characters that are
// not allowed in standard env variable names (letters, digits, underscore).
func validateKeyChars(key string) error {
	for _, r := range key {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return fmt.Errorf("key contains invalid character %q", r)
		}
	}
	return nil
}
