package censor_test

import (
	"testing"

	"github.com/user/envdiff/internal/censor"
	"github.com/user/envdiff/internal/comparator"
)

func makeResult(key, status string, values map[string]string) comparator.Result {
	return comparator.Result{Key: key, Status: status, Values: values}
}

func TestApply_NoKeys_ReturnsUnchanged(t *testing.T) {
	c := censor.New(nil)
	input := []comparator.Result{
		makeResult("DB_HOST", "identical", map[string]string{"dev": "localhost", "prod": "localhost"}),
	}
	out := c.Apply(input)
	if out[0].Values["dev"] != "localhost" {
		t.Fatalf("expected localhost, got %s", out[0].Values["dev"])
	}
}

func TestApply_CensorsMatchingKey(t *testing.T) {
	c := censor.New([]string{"DB_PASSWORD"})
	input := []comparator.Result{
		makeResult("DB_PASSWORD", "mismatched", map[string]string{"dev": "secret", "prod": "hunter2"}),
	}
	out := c.Apply(input)
	for env, val := range out[0].Values {
		if val != "[CENSORED]" {
			t.Errorf("env %s: expected [CENSORED], got %s", env, val)
		}
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	c := censor.New([]string{"api_key"})
	input := []comparator.Result{
		makeResult("API_KEY", "identical", map[string]string{"dev": "abc123"}),
	}
	out := c.Apply(input)
	if out[0].Values["dev"] != "[CENSORED]" {
		t.Fatalf("expected [CENSORED], got %s", out[0].Values["dev"])
	}
}

func TestApply_CustomPlaceholder(t *testing.T) {
	c := censor.NewWithPlaceholder([]string{"TOKEN"}, "***")
	input := []comparator.Result{
		makeResult("TOKEN", "identical", map[string]string{"dev": "tok_live_xyz"}),
	}
	out := c.Apply(input)
	if out[0].Values["dev"] != "***" {
		t.Fatalf("expected ***, got %s", out[0].Values["dev"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	c := censor.New([]string{"SECRET"})
	original := []comparator.Result{
		makeResult("SECRET", "identical", map[string]string{"dev": "original"}),
	}
	c.Apply(original)
	if original[0].Values["dev"] != "original" {
		t.Fatal("original result was mutated")
	}
}

func TestApply_EmptyValueNotCensored(t *testing.T) {
	// A missing key has an empty string value; we leave it empty so status
	// information is preserved faithfully.
	c := censor.New([]string{"MISSING_KEY"})
	input := []comparator.Result{
		makeResult("MISSING_KEY", "missing", map[string]string{"dev": "val", "prod": ""}),
	}
	out := c.Apply(input)
	if out[0].Values["prod"] != "" {
		t.Fatalf("expected empty string for absent env, got %q", out[0].Values["prod"])
	}
	if out[0].Values["dev"] != "[CENSORED]" {
		t.Fatalf("expected [CENSORED] for present env, got %q", out[0].Values["dev"])
	}
}
