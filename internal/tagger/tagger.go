// Package tagger attaches free-form string tags to comparison results based
// on user-supplied rules. Each rule maps a key pattern (exact or prefix) to
// one or more tags. Tagged results can be filtered, grouped, or reported on
// independently of the built-in status field.
package tagger

import (
	"strings"

	"github.com/user/envdiff/internal/comparator"
)

// Rule associates a key pattern with the tags that should be applied when the
// pattern matches a result's key.
type Rule struct {
	// Pattern is matched as an exact key name or, when it ends with "*", as a
	// prefix (the trailing asterisk is stripped before comparison).
	Pattern string
	// Tags are the labels attached to matching results.
	Tags []string
}

// Tagger holds a compiled set of tagging rules.
type Tagger struct {
	rules []Rule
}

// New returns a Tagger for the given rules. Rules with an empty Pattern or no
// Tags are silently ignored.
func New(rules []Rule) *Tagger {
	valid := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if r.Pattern != "" && len(r.Tags) > 0 {
			valid = append(valid, r)
		}
	}
	return &Tagger{rules: valid}
}

// TaggedResult wraps a comparator.Result with the tags matched by the Tagger's
// rules.
type TaggedResult struct {
	comparator.Result
	Tags []string
}

// Apply iterates over results and returns TaggedResults. Results that match no
// rule receive an empty Tags slice.
func (t *Tagger) Apply(results []comparator.Result) []TaggedResult {
	out := make([]TaggedResult, 0, len(results))
	for _, r := range results {
		out = append(out, TaggedResult{
			Result: r,
			Tags:   t.tagsFor(r.Key),
		})
	}
	return out
}

// tagsFor returns the union of all tags from rules that match key.
func (t *Tagger) tagsFor(key string) []string {
	seen := map[string]struct{}{}
	var tags []string
	for _, r := range t.rules {
		if matches(r.Pattern, key) {
			for _, tag := range r.Tags {
				if _, ok := seen[tag]; !ok {
					seen[tag] = struct{}{}
					tags = append(tags, tag)
				}
			}
		}
	}
	return tags
}

// matches returns true when pattern (exact or prefix with trailing "*")
// matches key.
func matches(pattern, key string) bool {
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(strings.ToLower(key), strings.ToLower(strings.TrimSuffix(pattern, "*")))
	}
	return strings.EqualFold(pattern, key)
}
