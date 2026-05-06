package summary_test

import (
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/summary"
)

func makeResult(key, status string, values map[string]string, missingIn []string) comparator.Result {
	return comparator.Result{
		Key:       key,
		Status:    status,
		Values:    values,
		MissingIn: missingIn,
	}
}

func TestCompute_Empty(t *testing.T) {
	stats := summary.Compute(nil)
	if stats.TotalKeys != 0 || stats.HasDifferences() {
		t.Errorf("expected zero stats, got %+v", stats)
	}
}

func TestCompute_AllIdentical(t *testing.T) {
	results := []comparator.Result{
		makeResult("PORT", comparator.StatusOK, map[string]string{"dev": "8080", "prod": "8080"}, nil),
		makeResult("HOST", comparator.StatusOK, map[string]string{"dev": "localhost", "prod": "localhost"}, nil),
	}
	stats := summary.Compute(results)
	if stats.TotalKeys != 2 {
		t.Errorf("expected 2 total keys, got %d", stats.TotalKeys)
	}
	if stats.Identical != 2 {
		t.Errorf("expected 2 identical, got %d", stats.Identical)
	}
	if stats.HasDifferences() {
		t.Error("expected no differences")
	}
	if stats.EnvCount != 2 {
		t.Errorf("expected 2 envs, got %d", stats.EnvCount)
	}
}

func TestCompute_MissingAndMismatched(t *testing.T) {
	results := []comparator.Result{
		makeResult("DB_URL", comparator.StatusMissing, map[string]string{"dev": "postgres://"}, []string{"prod"}),
		makeResult("LOG_LEVEL", comparator.StatusMismatch, map[string]string{"dev": "debug", "prod": "warn"}, nil),
		makeResult("PORT", comparator.StatusOK, map[string]string{"dev": "3000", "prod": "3000"}, nil),
	}
	stats := summary.Compute(results)
	if stats.TotalKeys != 3 {
		t.Errorf("expected 3 total keys, got %d", stats.TotalKeys)
	}
	if stats.MissingKeys != 1 {
		t.Errorf("expected 1 missing, got %d", stats.MissingKeys)
	}
	if stats.Mismatched != 1 {
		t.Errorf("expected 1 mismatched, got %d", stats.Mismatched)
	}
	if stats.Identical != 1 {
		t.Errorf("expected 1 identical, got %d", stats.Identical)
	}
	if !stats.HasDifferences() {
		t.Error("expected differences")
	}
	if len(stats.AffectedEnvs) == 0 {
		t.Error("expected affected envs to be populated")
	}
}
