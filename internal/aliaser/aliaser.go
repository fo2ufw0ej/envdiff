package aliaser

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/comparator"
)

// Rule maps an old key name to a new key name.
type Rule struct {
	From string
	To   string
}

// Aliaser rewrites comparator results so that keys known under different
// names across environments are treated as the same key.
type Aliaser struct {
	rules []Rule
}

// New returns an Aliaser configured with the provided rules.
// Rules are validated; an error is returned if any rule is malformed.
func New(rules []Rule) (*Aliaser, error) {
	if err := validateRules(rules); err != nil {
		return nil, err
	}
	return &Aliaser{rules: rules}, nil
}

// Apply rewrites the Key field of each result whose key matches a rule's
// From value (case-insensitive). The original key is preserved in a
// copy; the slice itself is not mutated.
func (a *Aliaser) Apply(results []comparator.Result) []comparator.Result {
	out := make([]comparator.Result, len(results))
	for i, r := range results {
		if alias, ok := a.resolve(r.Key); ok {
			r.Key = alias
		}
		out[i] = r
	}
	return out
}

// resolve returns the alias for key if a matching rule exists.
func (a *Aliaser) resolve(key string) (string, bool) {
	lower := strings.ToLower(key)
	for _, rule := range a.rules {
		if strings.ToLower(rule.From) == lower {
			return rule.To, true
		}
	}
	return "", false
}

func validateRules(rules []Rule) error {
	for i, r := range rules {
		if strings.TrimSpace(r.From) == "" {
			return fmt.Errorf("aliaser: rule %d has empty From field", i)
		}
		if strings.TrimSpace(r.To) == "" {
			return fmt.Errorf("aliaser: rule %d has empty To field", i)
		}
	}
	return nil
}
