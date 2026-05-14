package normalizer_test

import (
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/normalizer"
)

func makeResult(key, status string, values map[string]string) comparator.Result {
	return comparator.Result{Key: key, Status: status, Values: values}
}

func TestNormalize_NoOptions(t *testing.T) {
	input := []comparator.Result{
		makeResult("db_host", "identical", map[string]string{"prod": "localhost"}),
	}
	out := normalizer.Normalize(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Key != "db_host" {
		t.Errorf("unexpected key %q", out[0].Key)
	}
}

func TestNormalize_UpperKeys(t *testing.T) {
	input := []comparator.Result{
		makeResult("db_host", "identical", map[string]string{"prod": "localhost"}),
		makeResult("Api_Key", "missing", map[string]string{"dev": "abc"}),
	}
	out := normalizer.Normalize(input, normalizer.WithUpperKeys())
	if out[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", out[0].Key)
	}
	if out[1].Key != "API_KEY" {
		t.Errorf("expected API_KEY, got %q", out[1].Key)
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	input := []comparator.Result{
		makeResult("HOST", "mismatch", map[string]string{
			"prod": "  localhost  ",
			"dev":  "\t127.0.0.1\n",
		}),
	}
	out := normalizer.Normalize(input, normalizer.WithTrimValues())
	if out[0].Values["prod"] != "localhost" {
		t.Errorf("expected trimmed value, got %q", out[0].Values["prod"])
	}
	if out[0].Values["dev"] != "127.0.0.1" {
		t.Errorf("expected trimmed value, got %q", out[0].Values["dev"])
	}
}

func TestNormalize_StripPrefix(t *testing.T) {
	input := []comparator.Result{
		makeResult("APP_DB_HOST", "identical", map[string]string{"prod": "db"}),
		makeResult("APP_SECRET", "missing", map[string]string{"prod": "x"}),
		makeResult("OTHER_KEY", "identical", map[string]string{"prod": "y"}),
	}
	out := normalizer.Normalize(input, normalizer.WithStripPrefix("APP_"))
	if out[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", out[0].Key)
	}
	if out[1].Key != "SECRET" {
		t.Errorf("expected SECRET, got %q", out[1].Key)
	}
	if out[2].Key != "OTHER_KEY" {
		t.Errorf("expected OTHER_KEY unchanged, got %q", out[2].Key)
	}
}

func TestNormalize_DoesNotMutateInput(t *testing.T) {
	input := []comparator.Result{
		makeResult("my_key", "identical", map[string]string{"prod": "  v  "}),
	}
	_ = normalizer.Normalize(input, normalizer.WithUpperKeys(), normalizer.WithTrimValues())
	if input[0].Key != "my_key" {
		t.Error("original key was mutated")
	}
	if input[0].Values["prod"] != "  v  " {
		t.Error("original value was mutated")
	}
}
