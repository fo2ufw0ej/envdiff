package counter_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/comparator"
	"github.com/yourorg/envdiff/internal/counter"
)

func makeResult(key, status string, values map[string]string) comparator.Result {
	return comparator.Result{Key: key, Status: status, Values: values}
}

func TestCount_Empty(t *testing.T) {
	report := counter.Count(nil)
	if len(report) != 0 {
		t.Fatalf("expected empty report, got %d entries", len(report))
	}
}

func TestCount_AllIdentical(t *testing.T) {
	results := []comparator.Result{
		makeResult("KEY", comparator.StatusIdentical, map[string]string{"dev": "val", "prod": "val"}),
	}
	report := counter.Count(results)

	for _, env := range []string{"dev", "prod"} {
		c := report[env]
		if c.Identical != 1 || c.Missing != 0 || c.Mismatched != 0 || c.Total != 1 {
			t.Errorf("env %s: unexpected counts %+v", env, c)
		}
	}
}

func TestCount_MissingKey(t *testing.T) {
	results := []comparator.Result{
		makeResult("SECRET", comparator.StatusMissing, map[string]string{"dev": "abc"}),
	}
	report := counter.Count(results)

	if report["dev"].Missing != 1 {
		t.Errorf("expected 1 missing for dev, got %d", report["dev"].Missing)
	}
	if _, ok := report["prod"]; ok {
		t.Error("prod should not appear in report")
	}
}

func TestCount_MixedStatuses(t *testing.T) {
	results := []comparator.Result{
		makeResult("A", comparator.StatusIdentical, map[string]string{"dev": "x", "prod": "x"}),
		makeResult("B", comparator.StatusMismatched, map[string]string{"dev": "1", "prod": "2"}),
		makeResult("C", comparator.StatusMissing, map[string]string{"dev": "only"}),
	}
	report := counter.Count(results)

	dev := report["dev"]
	if dev.Identical != 1 || dev.Mismatched != 1 || dev.Missing != 1 || dev.Total != 3 {
		t.Errorf("dev counts wrong: %+v", dev)
	}

	prod := report["prod"]
	if prod.Identical != 1 || prod.Mismatched != 1 || prod.Missing != 0 || prod.Total != 2 {
		t.Errorf("prod counts wrong: %+v", prod)
	}
}

func TestCount_EnvNamesSorted(t *testing.T) {
	results := []comparator.Result{
		makeResult("X", comparator.StatusIdentical, map[string]string{"z": "1", "a": "1", "m": "1"}),
	}
	report := counter.Count(results)
	names := report.EnvNames()

	expected := []string{"a", "m", "z"}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("position %d: got %q want %q", i, name, expected[i])
		}
	}
}
