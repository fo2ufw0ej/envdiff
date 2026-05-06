package sorter_test

import (
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/sorter"
)

func makeResults() []comparator.Result {
	return []comparator.Result{
		{Key: "ZEBRA", Values: map[string]string{"prod": "z", "dev": "z2"}},
		{Key: "ALPHA", Values: map[string]string{"prod": "", "dev": "a"}},
		{Key: "MANGO", Values: map[string]string{"prod": "m1", "dev": "m2"}},
		{Key: "BETA", Values: map[string]string{"prod": "", "dev": "b"}},
	}
}

func TestSort_ByKey(t *testing.T) {
	results := makeResults()
	sorted := sorter.Sort(results, sorter.SortByKey)

	expected := []string{"ALPHA", "BETA", "MANGO", "ZEBRA"}
	for i, r := range sorted {
		if r.Key != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, r.Key, expected[i])
		}
	}
}

func TestSort_ByStatus(t *testing.T) {
	results := makeResults()
	sorted := sorter.Sort(results, sorter.SortByStatus)

	// ALPHA and BETA are missing (rank 0), ZEBRA and MANGO are mismatched (rank 1)
	for i := 0; i < 2; i++ {
		val := sorted[i].Values["prod"]
		if val != "" {
			t.Errorf("index %d: expected missing entry, got key %q with prod=%q", i, sorted[i].Key, val)
		}
	}
	for i := 2; i < 4; i++ {
		val := sorted[i].Values["prod"]
		if val == "" {
			t.Errorf("index %d: expected mismatched entry, got key %q", i, sorted[i].Key)
		}
	}
}

func TestSort_ByEnv(t *testing.T) {
	results := []comparator.Result{
		{Key: "Z", Values: map[string]string{"staging": "1", "dev": "2"}},
		{Key: "A", Values: map[string]string{"prod": "3", "dev": "4"}},
	}
	sorted := sorter.Sort(results, sorter.SortByEnv)

	// "dev" < "prod" < "staging", so first env of Z is "dev", first env of A is "dev" too
	// tie-break by key: A < Z
	if sorted[0].Key != "A" {
		t.Errorf("expected A first, got %q", sorted[0].Key)
	}
}

func TestSort_DoesNotMutateOriginal(t *testing.T) {
	results := makeResults()
	originalFirst := results[0].Key

	sorter.Sort(results, sorter.SortByKey)

	if results[0].Key != originalFirst {
		t.Errorf("original slice was mutated: got %q, want %q", results[0].Key, originalFirst)
	}
}
