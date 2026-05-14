package scoper_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/comparator"
	"github.com/yourorg/envdiff/internal/scoper"
)

func makeResult(key, status string, vals map[string]string) comparator.Result {
	return comparator.Result{Key: key, Status: status, Values: vals}
}

func TestApply_NoScope_ReturnsAll(t *testing.T) {
	s := scoper.New(nil)
	input := []comparator.Result{
		makeResult("KEY", "identical", map[string]string{"dev": "v", "prod": "v"}),
	}
	out := s.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
}

func TestApply_FiltersToScope(t *testing.T) {
	s := scoper.New([]string{"dev"})
	input := []comparator.Result{
		makeResult("KEY", "mismatch", map[string]string{"dev": "a", "prod": "b"}),
	}
	out := s.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if _, ok := out[0].Values["prod"]; ok {
		t.Error("prod should have been removed from scope")
	}
	if out[0].Values["dev"] != "a" {
		t.Errorf("expected dev=a, got %s", out[0].Values["dev"])
	}
}

func TestApply_DropsResultWithNoMatchingEnv(t *testing.T) {
	s := scoper.New([]string{"staging"})
	input := []comparator.Result{
		makeResult("KEY", "missing", map[string]string{"dev": "a", "prod": "b"}),
	}
	out := s.Apply(input)
	if len(out) != 0 {
		t.Fatalf("expected 0 results, got %d", len(out))
	}
}

func TestApply_MultipleEnvsInScope(t *testing.T) {
	s := scoper.New([]string{"dev", "staging"})
	input := []comparator.Result{
		makeResult("A", "identical", map[string]string{"dev": "1", "staging": "1", "prod": "1"}),
		makeResult("B", "missing", map[string]string{"prod": "x"}),
	}
	out := s.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 result (B dropped), got %d", len(out))
	}
	if _, ok := out[0].Values["prod"]; ok {
		t.Error("prod should not appear in scoped result")
	}
}

func TestEnvs_ReturnsSorted(t *testing.T) {
	s := scoper.New([]string{"prod", "dev", "staging"})
	envs := s.Envs()
	expected := []string{"dev", "prod", "staging"}
	if len(envs) != len(expected) {
		t.Fatalf("expected %d envs, got %d", len(expected), len(envs))
	}
	for i, e := range expected {
		if envs[i] != e {
			t.Errorf("pos %d: expected %s, got %s", i, e, envs[i])
		}
	}
}

func TestEnvs_EmptyScope(t *testing.T) {
	s := scoper.New([]string{})
	if s.Envs() != nil {
		t.Error("expected nil for empty scope")
	}
}
