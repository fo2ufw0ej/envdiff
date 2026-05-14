package trimmer_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/comparator"
	"github.com/your-org/envdiff/internal/trimmer"
)

func makeResult(key string, values map[string]string) comparator.Result {
	return comparator.Result{
		Key:    key,
		Values: values,
	}
}

func resultKeys(results []comparator.Result) []string {
	keys := make([]string, len(results))
	for i, r := range results {
		keys[i] = r.Key
	}
	return keys
}

func TestTrim_NoOptions_ReturnsAll(t *testing.T) {
	input := []comparator.Result{
		makeResult("FOO", map[string]string{"prod": "a"}),
		makeResult("BAR", map[string]string{"prod": "b"}),
	}
	out := trimmer.Trim(input, trimmer.Options{})
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestTrim_AllowKeys_FiltersOthers(t *testing.T) {
	input := []comparator.Result{
		makeResult("FOO", map[string]string{"prod": "a"}),
		makeResult("BAR", map[string]string{"prod": "b"}),
		makeResult("BAZ", map[string]string{"prod": "c"}),
	}
	out := trimmer.Trim(input, trimmer.Options{AllowKeys: []string{"FOO", "BAZ"}})
	keys := resultKeys(out)
	if len(keys) != 2 || keys[0] != "FOO" || keys[1] != "BAZ" {
		t.Fatalf("unexpected keys: %v", keys)
	}
}

func TestTrim_AllowKeys_CaseInsensitive(t *testing.T) {
	input := []comparator.Result{
		makeResult("foo", map[string]string{"prod": "a"}),
	}
	out := trimmer.Trim(input, trimmer.Options{AllowKeys: []string{"FOO"}})
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
}

func TestTrim_BlockKeys_RemovesMatched(t *testing.T) {
	input := []comparator.Result{
		makeResult("SECRET", map[string]string{"prod": "x"}),
		makeResult("PORT", map[string]string{"prod": "8080"}),
	}
	out := trimmer.Trim(input, trimmer.Options{BlockKeys: []string{"SECRET"}})
	if len(out) != 1 || out[0].Key != "PORT" {
		t.Fatalf("unexpected result: %v", out)
	}
}

func TestTrim_MaxValueLen_DropsLongValues(t *testing.T) {
	input := []comparator.Result{
		makeResult("SHORT", map[string]string{"prod": "hi"}),
		makeResult("LONG", map[string]string{"prod": "this-is-a-very-long-value"}),
	}
	out := trimmer.Trim(input, trimmer.Options{MaxValueLen: 5})
	if len(out) != 1 || out[0].Key != "SHORT" {
		t.Fatalf("expected only SHORT, got %v", resultKeys(out))
	}
}

func TestTrim_DoesNotMutateInput(t *testing.T) {
	input := []comparator.Result{
		makeResult("A", map[string]string{"prod": "1"}),
		makeResult("B", map[string]string{"prod": "2"}),
	}
	copy := append([]comparator.Result(nil), input...)
	trimmer.Trim(input, trimmer.Options{BlockKeys: []string{"A"}})
	if len(input) != len(copy) {
		t.Fatal("input was mutated")
	}
}
