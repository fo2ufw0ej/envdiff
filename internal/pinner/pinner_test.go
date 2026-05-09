package pinner_test

import (
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/pinner"
)

func makeResult(key, status string) comparator.Result {
	return comparator.Result{
		Key:    key,
		Status: status,
		Values: map[string]string{"prod": "val1", "staging": "val2"},
	}
}

func TestIsPinned_CaseInsensitive(t *testing.T) {
	p := pinner.New([]string{"DB_HOST", "api_key"})
	if !p.IsPinned("db_host") {
		t.Error("expected db_host to be pinned")
	}
	if !p.IsPinned("API_KEY") {
		t.Error("expected API_KEY to be pinned")
	}
	if p.IsPinned("SECRET") {
		t.Error("expected SECRET to not be pinned")
	}
}

func TestApply_RemovesMismatchForPinnedKey(t *testing.T) {
	p := pinner.New([]string{"DB_HOST"})
	input := []comparator.Result{
		makeResult("DB_HOST", comparator.StatusMismatch),
		makeResult("APP_ENV", comparator.StatusMismatch),
	}
	out := p.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Key != "APP_ENV" {
		t.Errorf("expected APP_ENV, got %s", out[0].Key)
	}
}

func TestApply_PreservesMissingForPinnedKey(t *testing.T) {
	p := pinner.New([]string{"DB_HOST"})
	input := []comparator.Result{
		makeResult("DB_HOST", comparator.StatusMissing),
	}
	out := p.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected missing entry to be preserved, got %d results", len(out))
	}
}

func TestApply_NoPins(t *testing.T) {
	p := pinner.New(nil)
	input := []comparator.Result{
		makeResult("DB_HOST", comparator.StatusMismatch),
		makeResult("APP_ENV", comparator.StatusIdentical),
	}
	out := p.Apply(input)
	if len(out) != len(input) {
		t.Errorf("expected %d results, got %d", len(input), len(out))
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	p := pinner.New([]string{"SECRET"})
	input := []comparator.Result{
		makeResult("SECRET", comparator.StatusMismatch),
		makeResult("PORT", comparator.StatusMismatch),
	}
	copy := make([]comparator.Result, len(input))
	for i, r := range input {
		copy[i] = r
	}
	p.Apply(input)
	for i, r := range input {
		if r.Key != copy[i].Key || r.Status != copy[i].Status {
			t.Error("Apply mutated original slice")
		}
	}
}
