package redactor

import "strings"

// DefaultSensitivePatterns contains common substrings that indicate a key
// holds a sensitive value and should be redacted in output.
var DefaultSensitivePatterns = []string{
	"PASSWORD",
	"SECRET",
	"TOKEN",
	"API_KEY",
	"PRIVATE",
	"CREDENTIAL",
	"AUTH",
}

// Redactor masks sensitive values in comparison results.
type Redactor struct {
	patterns []string
	mask     string
}

// New returns a Redactor that replaces sensitive values with mask.
// If patterns is empty, DefaultSensitivePatterns is used.
func New(patterns []string, mask string) *Redactor {
	if len(patterns) == 0 {
		patterns = DefaultSensitivePatterns
	}
	if mask == "" {
		mask = "***"
	}
	upper := make([]string, len(patterns))
	for i, p := range patterns {
		upper[i] = strings.ToUpper(p)
	}
	return &Redactor{patterns: upper, mask: mask}
}

// IsSensitive reports whether the given key matches any sensitive pattern.
func (r *Redactor) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range r.patterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// RedactValue returns the masked string if the key is sensitive,
// otherwise it returns the original value unchanged.
func (r *Redactor) RedactValue(key, value string) string {
	if r.IsSensitive(key) {
		return r.mask
	}
	return value
}

// RedactMap returns a new map with sensitive values replaced by the mask.
func (r *Redactor) RedactMap(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = r.RedactValue(k, v)
	}
	return out
}
