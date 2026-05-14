package deduper_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/comparator"
	"github.com/yourorg/envdiff/internal/deduper"
)

func makeResult(key string, values map[string]string) comparator.Result {
	return comparator.Result{Key: key, Values: values}
}

func keys(results []comparator.Result) []string {
	out := make([]string, len(results))
	for i, r := range results {
		out[i] = r.Key
	}
	return out
}

func TestDedupe_EmptyInput(t *testing.T) {
	got := deduper.Dedupe(nil)
	if len(got) != 0 {
		t.Fatalf("expected empty, got %d results", len(got))
	}
}

func TestDedupe_NoDuplicates(t *testing.T) {
	input := []comparator.Result{
		makeResult("DB_HOST", map[string]string{"prod": "db.prod", "dev": "localhost"}),
		makeResult("DB_PORT", map[string]string{"prod": "5432", "dev": "5432"}),
	}
	got := deduper.Dedupe(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestDedupe_RemovesDuplicateKey(t *testing.T) {
	base := makeResult("API_KEY", map[string]string{"prod": "secret", "dev": "dev-secret"})
	input := []comparator.Result{base, base}
	got := deduper.Dedupe(input)
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0].Key != "API_KEY" {
		t.Errorf("unexpected key %q", got[0].Key)
	}
}

func TestDedupe_PreservesOrder(t *testing.T) {
	input := []comparator.Result{
		makeResult("Z_KEY", map[string]string{"prod": "z"}),
		makeResult("A_KEY", map[string]string{"prod": "a"}),
		makeResult("M_KEY", map[string]string{"prod": "m"}),
	}
	got := deduper.Dedupe(input)
	want := []string{"Z_KEY", "A_KEY", "M_KEY"}
	for i, k := range want {
		if got[i].Key != k {
			t.Errorf("position %d: want %q, got %q", i, k, got[i].Key)
		}
	}
}

func TestDedupeStrict_DifferentValues_BothKept(t *testing.T) {
	r1 := makeResult("DB_HOST", map[string]string{"prod": "db1.prod"})
	r2 := makeResult("DB_HOST", map[string]string{"prod": "db2.prod"})
	got := deduper.DedupeStrict([]comparator.Result{r1, r2})
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestDedupeStrict_SameValues_OnlyFirstKept(t *testing.T) {
	r := makeResult("DB_HOST", map[string]string{"prod": "db.prod", "dev": "localhost"})
	got := deduper.DedupeStrict([]comparator.Result{r, r, r})
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
}

func TestDedupeStrict_EmptyInput(t *testing.T) {
	got := deduper.DedupeStrict([]comparator.Result{})
	if got == nil || len(got) != 0 {
		t.Fatalf("expected empty non-nil slice")
	}
}
