package pivotter_test

import (
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/pivotter"
)

func makeResult(key string, values map[string]string) comparator.Result {
	return comparator.Result{Key: key, Values: values}
}

func TestPivot_Empty(t *testing.T) {
	table := pivotter.Pivot(nil)
	if len(table.Rows) != 0 {
		t.Fatalf("expected 0 rows, got %d", len(table.Rows))
	}
	if len(table.Envs) != 0 {
		t.Fatalf("expected 0 envs, got %d", len(table.Envs))
	}
}

func TestPivot_SingleEnv(t *testing.T) {
	results := []comparator.Result{
		makeResult("PORT", map[string]string{"prod": "8080"}),
		makeResult("HOST", map[string]string{"prod": "localhost"}),
	}
	table := pivotter.Pivot(results)

	if len(table.Envs) != 1 || table.Envs[0] != "prod" {
		t.Fatalf("unexpected envs: %v", table.Envs)
	}
	if len(table.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(table.Rows))
	}
	// Rows must be sorted by key.
	if table.Rows[0].Key != "HOST" || table.Rows[1].Key != "PORT" {
		t.Fatalf("rows not sorted: %v %v", table.Rows[0].Key, table.Rows[1].Key)
	}
}

func TestPivot_AbsentSentinel(t *testing.T) {
	results := []comparator.Result{
		makeResult("DEBUG", map[string]string{"dev": "true"}),
		makeResult("DEBUG", map[string]string{"prod": "false"}),
	}
	table := pivotter.Pivot(results)

	if len(table.Envs) != 2 {
		t.Fatalf("expected 2 envs, got %d", len(table.Envs))
	}
	row := table.Rows[0]
	if row.Values["dev"] != "true" {
		t.Errorf("dev value wrong: %s", row.Values["dev"])
	}
	if row.Values["prod"] != "false" {
		t.Errorf("prod value wrong: %s", row.Values["prod"])
	}
}

func TestPivot_MissingKeyGetsAbsent(t *testing.T) {
	results := []comparator.Result{
		makeResult("ONLY_IN_DEV", map[string]string{"dev": "yes"}),
	}
	// Simulate a second env that never defines ONLY_IN_DEV.
	results = append(results, makeResult("SHARED", map[string]string{"dev": "a", "prod": "a"}))

	table := pivotter.Pivot(results)

	for _, row := range table.Rows {
		if row.Key == "ONLY_IN_DEV" {
			if row.Values["prod"] != pivotter.Absent {
				t.Errorf("expected Absent for prod, got %q", row.Values["prod"])
			}
			return
		}
	}
	t.Fatal("ONLY_IN_DEV row not found")
}

func TestPivot_EnvsSorted(t *testing.T) {
	results := []comparator.Result{
		makeResult("K", map[string]string{"staging": "1", "dev": "2", "prod": "3"}),
	}
	table := pivotter.Pivot(results)
	want := []string{"dev", "prod", "staging"}
	for i, e := range table.Envs {
		if e != want[i] {
			t.Errorf("env[%d]: got %q want %q", i, e, want[i])
		}
	}
}
