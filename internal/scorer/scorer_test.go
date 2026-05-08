package scorer_test

import (
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/scorer"
)

func makeResult(key, status string) comparator.Result {
	return comparator.Result{
		Key:    key,
		Status: status,
		Values: map[string]string{},
	}
}

func TestCompute_Empty(t *testing.T) {
	s := scorer.Compute(nil)
	if s.TotalKeys != 0 {
		t.Errorf("expected 0 total keys, got %d", s.TotalKeys)
	}
	if s.HealthPercent != 100 {
		t.Errorf("expected 100%% health for empty input, got %.2f", s.HealthPercent)
	}
}

func TestCompute_AllIdentical(t *testing.T) {
	results := []comparator.Result{
		makeResult("A", comparator.StatusIdentical),
		makeResult("B", comparator.StatusIdentical),
	}
	s := scorer.Compute(results)
	if s.TotalKeys != 2 {
		t.Errorf("expected 2 total, got %d", s.TotalKeys)
	}
	if s.IdenticalKeys != 2 {
		t.Errorf("expected 2 identical, got %d", s.IdenticalKeys)
	}
	if s.HealthPercent != 100 {
		t.Errorf("expected 100%% health, got %.2f", s.HealthPercent)
	}
}

func TestCompute_MixedStatuses(t *testing.T) {
	results := []comparator.Result{
		makeResult("A", comparator.StatusIdentical),
		makeResult("B", comparator.StatusMissing),
		makeResult("C", comparator.StatusMismatched),
		makeResult("D", comparator.StatusMissing),
	}
	s := scorer.Compute(results)
	if s.TotalKeys != 4 {
		t.Errorf("expected 4 total, got %d", s.TotalKeys)
	}
	if s.IdenticalKeys != 1 {
		t.Errorf("expected 1 identical, got %d", s.IdenticalKeys)
	}
	if s.MissingKeys != 2 {
		t.Errorf("expected 2 missing, got %d", s.MissingKeys)
	}
	if s.MismatchedKeys != 1 {
		t.Errorf("expected 1 mismatched, got %d", s.MismatchedKeys)
	}
	want := 25.0
	if s.HealthPercent != want {
		t.Errorf("expected %.2f%% health, got %.2f", want, s.HealthPercent)
	}
}

func TestCompute_AllMissing(t *testing.T) {
	results := []comparator.Result{
		makeResult("X", comparator.StatusMissing),
		makeResult("Y", comparator.StatusMissing),
	}
	s := scorer.Compute(results)
	if s.HealthPercent != 0 {
		t.Errorf("expected 0%% health, got %.2f", s.HealthPercent)
	}
}
