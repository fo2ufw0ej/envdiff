package merger_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/merger"
)

func env(name string, vals map[string]string) merger.Env {
	return merger.Env{Name: name, Values: vals}
}

func TestMerge_SingleEnv(t *testing.T) {
	envs := []merger.Env{
		env("prod", map[string]string{"FOO": "bar", "PORT": "8080"}),
	}
	res := merger.Merge(envs, merger.StrategyFirst)
	if res.Values["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", res.Values["FOO"])
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
}

func TestMerge_NoConflict(t *testing.T) {
	envs := []merger.Env{
		env("dev", map[string]string{"FOO": "same"}),
		env("prod", map[string]string{"FOO": "same"}),
	}
	res := merger.Merge(envs, merger.StrategyFirst)
	if res.Values["FOO"] != "same" {
		t.Errorf("expected FOO=same, got %s", res.Values["FOO"])
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
}

func TestMerge_StrategyFirst(t *testing.T) {
	envs := []merger.Env{
		env("dev", map[string]string{"FOO": "dev-value"}),
		env("prod", map[string]string{"FOO": "prod-value"}),
	}
	res := merger.Merge(envs, merger.StrategyFirst)
	if res.Values["FOO"] != "dev-value" {
		t.Errorf("StrategyFirst: expected dev-value, got %s", res.Values["FOO"])
	}
	if _, ok := res.Conflicts["FOO"]; !ok {
		t.Error("expected conflict for FOO")
	}
}

func TestMerge_StrategyLast(t *testing.T) {
	envs := []merger.Env{
		env("dev", map[string]string{"FOO": "dev-value"}),
		env("prod", map[string]string{"FOO": "prod-value"}),
	}
	res := merger.Merge(envs, merger.StrategyLast)
	if res.Values["FOO"] != "prod-value" {
		t.Errorf("StrategyLast: expected prod-value, got %s", res.Values["FOO"])
	}
}

func TestMerge_UnionKeys(t *testing.T) {
	envs := []merger.Env{
		env("dev", map[string]string{"ONLY_DEV": "1"}),
		env("prod", map[string]string{"ONLY_PROD": "2"}),
	}
	res := merger.Merge(envs, merger.StrategyFirst)
	if res.Values["ONLY_DEV"] != "1" {
		t.Errorf("expected ONLY_DEV=1")
	}
	if res.Values["ONLY_PROD"] != "2" {
		t.Errorf("expected ONLY_PROD=2")
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("unique keys should not conflict, got %v", res.Conflicts)
	}
}

func TestMerge_ConflictLists(t *testing.T) {
	envs := []merger.Env{
		env("a", map[string]string{"KEY": "x"}),
		env("b", map[string]string{"KEY": "y"}),
		env("c", map[string]string{"KEY": "z"}),
	}
	res := merger.Merge(envs, merger.StrategyFirst)
	names := res.Conflicts["KEY"]
	if len(names) != 3 {
		t.Errorf("expected 3 conflict sources, got %d: %v", len(names), names)
	}
}
