// Package annotator attaches human-readable annotations (hints/suggestions)
// to comparison results based on their status and key patterns.
package annotator

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/comparator"
)

// Annotation holds a message and severity level for a single result entry.
type Annotation struct {
	Key      string
	Env      string
	Severity string // "info", "warn", "error"
	Message  string
}

// Annotator produces annotations for comparison results.
type Annotator struct {
	sensitivePatterns []string
}

// New returns an Annotator with default sensitive key patterns.
func New() *Annotator {
	return &Annotator{
		sensitivePatterns: []string{"secret", "password", "token", "key", "api"},
	}
}

// NewWithPatterns returns an Annotator with custom sensitive key patterns.
func NewWithPatterns(patterns []string) *Annotator {
	return &Annotator{sensitivePatterns: patterns}
}

// Annotate returns a slice of Annotation for the given results.
func (a *Annotator) Annotate(results []comparator.Result) []Annotation {
	var annotations []Annotation
	for _, r := range results {
		for env, status := range r.Statuses {
			switch status {
			case comparator.StatusMissing:
				annotations = append(annotations, Annotation{
					Key:      r.Key,
					Env:      env,
					Severity: "error",
					Message:  fmt.Sprintf("key %q is missing in environment %q", r.Key, env),
				})
			case comparator.StatusMismatch:
				sev := "warn"
				msg := fmt.Sprintf("key %q has a mismatched value in environment %q", r.Key, env)
				if a.isSensitive(r.Key) {
					sev = "error"
					msg += " (sensitive key — verify intentional)"
				}
				annotations = append(annotations, Annotation{
					Key:      r.Key,
					Env:      env,
					Severity: sev,
					Message:  msg,
				})
			}
		}
	}
	return annotations
}

func (a *Annotator) isSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range a.sensitivePatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}
