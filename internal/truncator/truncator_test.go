package truncator_test

import (
	"strings"
	"testing"

	"github.com/your-org/envdiff/internal/comparator"
	"github.com/your-org/envdiff/internal/truncator"
)

func makeResult(key, status string, vals map[string]string) comparator.Result {
	return comparator.Result{Key: key, Status: status, Values: vals}
}

func TestTruncate_ShortValues_Unchanged(t *testing.T) {
	input := []comparator.Result{
		makeResult("KEY", "identical", map[string]string{"prod": "hello", "staging": "hello"}),
	}
	out := truncator.Truncate(input, truncator.Options{MaxLen: 64})
	if out[0].Values["prod"] != "hello" {
		t.Fatalf("expected 'hello', got %q", out[0].Values["prod"])
	}
}

func TestTruncate_LongValue_Truncated(t *testing.T) {
	long := strings.Repeat("a", 100)
	input := []comparator.Result{
		makeResult("KEY", "mismatch", map[string]string{"prod": long}),
	}
	out := truncator.Truncate(input, truncator.Options{MaxLen: 10, Suffix: "..."})
	got := out[0].Values["prod"]
	if len([]rune(got)) != 13 { // 10 + len("...")
		t.Fatalf("unexpected length %d: %q", len([]rune(got)), got)
	}
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("expected suffix '...', got %q", got)
	}
}

func TestTruncate_DefaultOptions(t *testing.T) {
	long := strings.Repeat("x", 80)
	input := []comparator.Result{
		makeResult("K", "missing", map[string]string{"dev": long}),
	}
	out := truncator.Truncate(input, truncator.Options{})
	got := out[0].Values["dev"]
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("expected default suffix, got %q", got)
	}
	if len([]rune(got)) != 67 { // 64 + 3
		t.Fatalf("unexpected length %d", len([]rune(got)))
	}
}

func TestTruncate_DoesNotMutateOriginal(t *testing.T) {
	long := strings.Repeat("z", 100)
	input := []comparator.Result{
		makeResult("K", "mismatch", map[string]string{"prod": long}),
	}
	truncator.Truncate(input, truncator.Options{MaxLen: 5})
	if input[0].Values["prod"] != long {
		t.Fatal("original result was mutated")
	}
}

func TestTruncate_NilValues(t *testing.T) {
	input := []comparator.Result{
		makeResult("K", "missing", nil),
	}
	out := truncator.Truncate(input, truncator.Options{})
	if out[0].Values != nil {
		t.Fatal("expected nil values map")
	}
}

func TestTruncate_CustomSuffix(t *testing.T) {
	long := strings.Repeat("b", 20)
	input := []comparator.Result{
		makeResult("K", "mismatch", map[string]string{"e": long}),
	}
	out := truncator.Truncate(input, truncator.Options{MaxLen: 5, Suffix: "[…]"})
	got := out[0].Values["e"]
	if !strings.HasSuffix(got, "[…]") {
		t.Fatalf("expected custom suffix, got %q", got)
	}
}
