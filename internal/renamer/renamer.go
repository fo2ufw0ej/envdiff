package renamer

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/comparator"
)

// Rule describes a key rename mapping from OldKey to NewKey.
type Rule struct {
	OldKey string
	NewKey string
}

// Result holds the outcome of applying a rename rule to a set of comparison results.
type Result struct {
	Rule    Rule
	Matched bool
	Updated []comparator.Result
}

// Apply renames keys in the provided comparator results according to the given
// rules. Keys that match OldKey are rewritten to NewKey in both the Key field
// and the per-environment Values map. Rules are applied in order; if a key
// matches multiple rules only the first match is used.
func Apply(results []comparator.Result, rules []Rule) ([]comparator.Result, []Result, error) {
	if err := validateRules(rules); err != nil {
		return nil, nil, err
	}

	ruleResults := make([]Result, len(rules))
	for i, r := range rules {
		ruleResults[i] = Result{Rule: r}
	}

	out := make([]comparator.Result, 0, len(results))
	for _, res := range results {
		applied := false
		for i, rule := range rules {
			if strings.EqualFold(res.Key, rule.OldKey) {
				renamed := comparator.Result{
					Key:    rule.NewKey,
					Status: res.Status,
					Values: res.Values,
				}
				out = append(out, renamed)
				ruleResults[i].Matched = true
				ruleResults[i].Updated = append(ruleResults[i].Updated, renamed)
				applied = true
				break
			}
		}
		if !applied {
			out = append(out, res)
		}
	}
	return out, ruleResults, nil
}

func validateRules(rules []Rule) error {
	seen := make(map[string]struct{}, len(rules))
	for _, r := range rules {
		if r.OldKey == "" {
			return fmt.Errorf("renamer: OldKey must not be empty")
		}
		if r.NewKey == "" {
			return fmt.Errorf("renamer: NewKey must not be empty for rule %q", r.OldKey)
		}
		key := strings.ToLower(r.OldKey)
		if _, dup := seen[key]; dup {
			return fmt.Errorf("renamer: duplicate OldKey %q", r.OldKey)
		}
		seen[key] = struct{}{}
	}
	return nil
}
