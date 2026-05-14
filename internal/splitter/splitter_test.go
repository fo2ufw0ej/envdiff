package splitter_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/comparator"
	"github.com/yourorg/envdiff/internal/splitter"
)

func makeResult(key string, values map[string]string, status comparator.Status) comparator.Result {
	return comparator.Result{Key: key, Values: values, Status: status}
}

func TestSplit_Empty(t *testing.T) {
	er := splitter.Split(nil)
	if len(er) != 0 {
		t.Fatalf("expected empty EnvResults, got %d buckets", len(er))
	}
}

func TestSplit_SingleEnv(t *testing.T) {
	results := []comparator.Result{
		makeResult("PORT", map[string]string{"prod": "8080"}, comparator.StatusMissing),
	}
	er := splitter.Split(results)
	if len(er["prod"]) != 1 {
		t.Fatalf("expected 1 result for prod, got %d", len(er["prod"]))
	}
}

func TestSplit_MultipleEnvs(t *testing.T) {
	results := []comparator.Result{
		makeResult("HOST", map[string]string{"dev": "localhost", "prod": "example.com"}, comparator.StatusMismatch),
		makeResult("PORT", map[string]string{"dev": "3000"}, comparator.StatusMissing),
	}
	er := splitter.Split(results)

	if len(er["dev"]) != 2 {
		t.Errorf("expected 2 results for dev, got %d", len(er["dev"]))
	}
	if len(er["prod"]) != 1 {
		t.Errorf("expected 1 result for prod, got %d", len(er["prod"]))
	}
}

func TestSplit_GlobalBucket(t *testing.T) {
	results := []comparator.Result{
		{Key: "ORPHAN", Values: map[string]string{}, Status: comparator.StatusIdentical},
	}
	er := splitter.Split(results)
	if len(er["_global"]) != 1 {
		t.Fatalf("expected 1 result in _global, got %d", len(er["_global"]))
	}
}

func TestEnvs_ExcludesGlobal(t *testing.T) {
	results := []comparator.Result{
		makeResult("A", map[string]string{"staging": "1"}, comparator.StatusMissing),
		{Key: "B", Values: map[string]string{}, Status: comparator.StatusIdentical},
	}
	er := splitter.Split(results)
	envs := splitter.Envs(er)

	if len(envs) != 1 || envs[0] != "staging" {
		t.Errorf("expected [staging], got %v", envs)
	}
}

func TestEnvs_Sorted(t *testing.T) {
	results := []comparator.Result{
		makeResult("X", map[string]string{"prod": "a", "dev": "b", "staging": "c"}, comparator.StatusMismatch),
	}
	er := splitter.Split(results)
	envs := splitter.Envs(er)

	want := []string{"dev", "prod", "staging"}
	for i, v := range want {
		if envs[i] != v {
			t.Errorf("envs[%d] = %q, want %q", i, envs[i], v)
		}
	}
}
