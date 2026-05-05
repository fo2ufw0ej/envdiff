package comparator_test

import (
	"sort"
	"testing"

	"github.com/yourorg/envdiff/internal/comparator"
)

func sortedKeys(m map[string][]string) map[string][]string {
	for k := range m {
		sort.Strings(m[k])
	}
	return m
}

func TestCompare_NoEnvs(t *testing.T) {
	result := comparator.Compare(nil)
	if len(result.MissingIn) != 0 || len(result.Mismatched) != 0 {
		t.Error("expected empty result for nil input")
	}
}

func TestCompare_MissingKeys(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"DB_HOST": "localhost", "DB_PORT": "5432"},
		"prod": {"DB_HOST": "prod.db"},
	}

	result := comparator.Compare(envs)
	sortedKeys(result.MissingIn)

	if missing, ok := result.MissingIn["prod"]; !ok || len(missing) != 1 || missing[0] != "DB_PORT" {
		t.Errorf("expected prod to be missing DB_PORT, got %v", result.MissingIn["prod"])
	}
	if _, ok := result.MissingIn["dev"]; ok && len(result.MissingIn["dev"]) > 0 {
		t.Errorf("dev should not be missing any keys, got %v", result.MissingIn["dev"])
	}
}

func TestCompare_MismatchedValues(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"LOG_LEVEL": "debug"},
		"prod": {"LOG_LEVEL": "error"},
	}

	result := comparator.Compare(envs)

	if len(result.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatched key, got %d", len(result.Mismatched))
	}
	if result.Mismatched[0].Key != "LOG_LEVEL" {
		t.Errorf("expected LOG_LEVEL mismatch, got %s", result.Mismatched[0].Key)
	}
}

func TestCompare_IdenticalEnvs(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"APP_ENV": "dev", "PORT": "8080"},
		"prod": {"APP_ENV": "dev", "PORT": "8080"},
	}

	result := comparator.Compare(envs)

	if len(result.Mismatched) != 0 {
		t.Errorf("expected no mismatches, got %v", result.Mismatched)
	}
	for env, keys := range result.MissingIn {
		if len(keys) > 0 {
			t.Errorf("expected no missing keys for %s, got %v", env, keys)
		}
	}
}
