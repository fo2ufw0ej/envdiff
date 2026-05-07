package ignorer_test

import (
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/ignorer"
)

// applyIgnorer filters comparator results using an Ignorer.
func applyIgnorer(results []comparator.Result, ig *ignorer.Ignorer) []comparator.Result {
	out := results[:0:0]
	for _, r := range results {
		if !ig.ShouldIgnore(r.Key) {
			out = append(out, r)
		}
	}
	return out
}

func TestIgnorer_WithComparatorResults(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {
			"HOST":       "prod.example.com",
			"SECRET_KEY": "abc123",
			"AWS_REGION": "us-east-1",
		},
		"staging": {
			"HOST":       "staging.example.com",
			"SECRET_KEY": "xyz789",
		},
	}

	results := comparator.Compare(envs)

	ig := ignorer.New([]string{"SECRET_KEY"}, []string{"AWS_"})
	filtered := applyIgnorer(results, ig)

	for _, r := range filtered {
		if ig.ShouldIgnore(r.Key) {
			t.Errorf("key %q should have been filtered out", r.Key)
		}
	}

	keys := make(map[string]bool)
	for _, r := range filtered {
		keys[r.Key] = true
	}
	if keys["SECRET_KEY"] {
		t.Error("SECRET_KEY should be absent after ignoring")
	}
	if keys["AWS_REGION"] {
		t.Error("AWS_REGION should be absent after ignoring")
	}
	if !keys["HOST"] {
		t.Error("HOST should remain after ignoring")
	}
}
