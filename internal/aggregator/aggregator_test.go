package aggregator_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/aggregator"
	"github.com/your-org/envdiff/internal/comparator"
)

func makeResult(key, status string, vals map[string]string) comparator.Result {
	return comparator.Result{Key: key, Status: status, Values: vals}
}

func TestAggregate_Empty(t *testing.T) {
	r := aggregator.Aggregate(nil)
	if r.TotalKeys != 0 || r.TotalIssues != 0 || len(r.Envs) != 0 {
		t.Fatalf("expected empty report, got %+v", r)
	}
}

func TestAggregate_AllIdentical(t *testing.T) {
	results := []comparator.Result{
		makeResult("PORT", comparator.StatusIdentical, map[string]string{"dev": "8080", "prod": "8080"}),
		makeResult("HOST", comparator.StatusIdentical, map[string]string{"dev": "localhost", "prod": "localhost"}),
	}
	r := aggregator.Aggregate(results)
	if r.TotalKeys != 2 {
		t.Fatalf("expected 2 total keys, got %d", r.TotalKeys)
	}
	if r.TotalIssues != 0 {
		t.Fatalf("expected 0 issues, got %d", r.TotalIssues)
	}
	if len(r.Envs) != 2 {
		t.Fatalf("expected 2 envs, got %d", len(r.Envs))
	}
}

func TestAggregate_MismatchCounted(t *testing.T) {
	results := []comparator.Result{
		makeResult("DB_URL", comparator.StatusMismatch, map[string]string{"dev": "localhost", "prod": "db.prod"}),
	}
	r := aggregator.Aggregate(results)
	if r.TotalIssues == 0 {
		t.Fatal("expected at least one issue")
	}
	for _, e := range r.Envs {
		if e.Mismatched != 1 {
			t.Errorf("env %s: expected 1 mismatched, got %d", e.Env, e.Mismatched)
		}
	}
}

func TestAggregate_SortedEnvs(t *testing.T) {
	results := []comparator.Result{
		makeResult("X", comparator.StatusIdentical, map[string]string{"z_env": "1", "a_env": "1"}),
	}
	r := aggregator.Aggregate(results)
	if len(r.Envs) < 2 {
		t.Fatal("expected 2 envs")
	}
	if r.Envs[0].Env != "a_env" {
		t.Errorf("expected a_env first, got %s", r.Envs[0].Env)
	}
}
