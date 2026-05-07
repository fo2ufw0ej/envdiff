package profiler_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/profiler"
)

func TestProfile_Empty(t *testing.T) {
	r := profiler.Profile(nil)
	if r.TotalKeys != 0 || r.TotalEnvs != 0 {
		t.Fatalf("expected empty report, got %+v", r)
	}
}

func TestProfile_SingleEnv(t *testing.T) {
	envs := map[string]map[string]string{
		"production": {"HOST": "prod.example.com", "PORT": "443"},
	}
	r := profiler.Profile(envs)
	if r.TotalKeys != 2 {
		t.Fatalf("expected 2 keys, got %d", r.TotalKeys)
	}
	if r.TotalEnvs != 1 {
		t.Fatalf("expected 1 env, got %d", r.TotalEnvs)
	}
	for _, s := range r.KeyStats {
		if s.EnvsPresent != 1 {
			t.Errorf("key %q: expected EnvsPresent=1, got %d", s.Key, s.EnvsPresent)
		}
		if s.UniqueValues != 1 {
			t.Errorf("key %q: expected UniqueValues=1, got %d", s.Key, s.UniqueValues)
		}
	}
}

func TestProfile_MissingKey(t *testing.T) {
	envs := map[string]map[string]string{
		"staging":    {"HOST": "staging.example.com"},
		"production": {"HOST": "prod.example.com", "SECRET": "abc"},
	}
	r := profiler.Profile(envs)
	if r.TotalKeys != 2 {
		t.Fatalf("expected 2 keys, got %d", r.TotalKeys)
	}
	for _, s := range r.KeyStats {
		if s.Key == "SECRET" && s.EnvsPresent != 1 {
			t.Errorf("SECRET should be present in 1 env, got %d", s.EnvsPresent)
		}
		if s.Key == "HOST" && s.EnvsPresent != 2 {
			t.Errorf("HOST should be present in 2 envs, got %d", s.EnvsPresent)
		}
	}
}

func TestProfile_UniqueValues(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"LOG_LEVEL": "debug"},
		"prod": {"LOG_LEVEL": "error"},
		"ci":   {"LOG_LEVEL": "debug"},
	}
	r := profiler.Profile(envs)
	if r.TotalKeys != 1 {
		t.Fatalf("expected 1 key, got %d", r.TotalKeys)
	}
	s := r.KeyStats[0]
	if s.UniqueValues != 2 {
		t.Errorf("expected 2 unique values, got %d", s.UniqueValues)
	}
}

func TestProfile_SortedKeys(t *testing.T) {
	envs := map[string]map[string]string{
		"dev": {"ZEBRA": "1", "ALPHA": "2", "MIDDLE": "3"},
	}
	r := profiler.Profile(envs)
	keys := make([]string, len(r.KeyStats))
	for i, s := range r.KeyStats {
		keys[i] = s.Key
	}
	expected := []string{"ALPHA", "MIDDLE", "ZEBRA"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("position %d: expected %q, got %q", i, k, keys[i])
		}
	}
}
