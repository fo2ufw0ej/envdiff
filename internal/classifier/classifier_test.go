package classifier_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/classifier"
	"github.com/yourorg/envdiff/internal/comparator"
)

func TestClassify_DefaultRules(t *testing.T) {
	c := classifier.New()
	cases := []struct {
		key  string
		want classifier.Category
	}{
		{"DATABASE_URL", classifier.CategoryDatabase},
		{"POSTGRES_HOST", classifier.CategoryDatabase},
		{"REDIS_ADDR", classifier.CategoryDatabase},
		{"JWT_SECRET", classifier.CategoryAuth},
		{"API_KEY", classifier.CategoryAuth},
		{"AUTH_TOKEN", classifier.CategoryAuth},
		{"SERVER_HOST", classifier.CategoryNetwork},
		{"GRPC_PORT", classifier.CategoryNetwork},
		{"LOG_LEVEL", classifier.CategoryObserv},
		{"SENTRY_DSN", classifier.CategoryDatabase}, // DSN matches database first
		{"FEATURE_DARK_MODE", classifier.CategoryFeature},
		{"ENABLE_BETA", classifier.CategoryFeature},
		{"APP_NAME", classifier.CategoryUnknown},
	}
	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			got := c.Classify(tc.key)
			if got != tc.want {
				t.Errorf("Classify(%q) = %q, want %q", tc.key, got, tc.want)
			}
		})
	}
}

func TestClassify_CustomRules(t *testing.T) {
	c := classifier.NewWithRules([]classifier.Rule{
		{Keywords: []string{"custom"}, Category: "custom_cat"},
	})
	if got := c.Classify("MY_CUSTOM_VAR"); got != "custom_cat" {
		t.Errorf("expected custom_cat, got %q", got)
	}
	if got := c.Classify("DATABASE_URL"); got != classifier.CategoryUnknown {
		t.Errorf("expected unknown with custom rules, got %q", got)
	}
}

func TestClassifyResults_AddsTag(t *testing.T) {
	c := classifier.New()
	input := []comparator.Result{
		{Key: "DATABASE_URL", Status: comparator.StatusMissing},
		{Key: "APP_NAME", Status: comparator.StatusIdentical},
	}
	out := c.ClassifyResults(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
	if out[0].Tags["category"] != string(classifier.CategoryDatabase) {
		t.Errorf("expected database, got %q", out[0].Tags["category"])
	}
	if out[1].Tags["category"] != string(classifier.CategoryUnknown) {
		t.Errorf("expected unknown, got %q", out[1].Tags["category"])
	}
}

func TestClassifyResults_DoesNotMutateOriginal(t *testing.T) {
	c := classifier.New()
	input := []comparator.Result{
		{Key: "API_KEY", Status: comparator.StatusMismatch},
	}
	_ = c.ClassifyResults(input)
	if input[0].Tags != nil {
		t.Error("original result should not be mutated")
	}
}
