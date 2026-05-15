package digester_test

import (
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/digester"
)

func makeResult(key string, status comparator.Status, values map[string]string) comparator.Result {
	return comparator.Result{Key: key, Status: status, Values: values}
}

func TestDigest_EmptyInput(t *testing.T) {
	d := digester.Digest(nil)
	if d == "" {
		t.Fatal("expected non-empty digest for empty input")
	}
}

func TestDigest_Deterministic(t *testing.T) {
	results := []comparator.Result{
		makeResult("DB_HOST", comparator.StatusMatch, map[string]string{"prod": "db.prod", "staging": "db.staging"}),
		makeResult("API_KEY", comparator.StatusMissing, map[string]string{"prod": "secret"}),
	}

	d1 := digester.Digest(results)
	d2 := digester.Digest(results)
	if d1 != d2 {
		t.Fatalf("digest not deterministic: %q vs %q", d1, d2)
	}
}

func TestDigest_OrderIndependent(t *testing.T) {
	a := []comparator.Result{
		makeResult("FOO", comparator.StatusMatch, map[string]string{"prod": "1", "dev": "1"}),
		makeResult("BAR", comparator.StatusMismatch, map[string]string{"prod": "x", "dev": "y"}),
	}
	b := []comparator.Result{
		makeResult("BAR", comparator.StatusMismatch, map[string]string{"dev": "y", "prod": "x"}),
		makeResult("FOO", comparator.StatusMatch, map[string]string{"dev": "1", "prod": "1"}),
	}

	if digester.Digest(a) != digester.Digest(b) {
		t.Fatal("digest should be order-independent")
	}
}

func TestDigest_ChangesOnValueDiff(t *testing.T) {
	a := []comparator.Result{
		makeResult("HOST", comparator.StatusMatch, map[string]string{"prod": "a"}),
	}
	b := []comparator.Result{
		makeResult("HOST", comparator.StatusMatch, map[string]string{"prod": "b"}),
	}

	if digester.Digest(a) == digester.Digest(b) {
		t.Fatal("digest should differ when values change")
	}
}

func TestEqual_SameResults(t *testing.T) {
	results := []comparator.Result{
		makeResult("PORT", comparator.StatusMatch, map[string]string{"prod": "8080"}),
	}
	if !digester.Equal(results, results) {
		t.Fatal("Equal should return true for identical result sets")
	}
}

func TestEqual_DifferentResults(t *testing.T) {
	a := []comparator.Result{makeResult("X", comparator.StatusMatch, map[string]string{"prod": "1"})}
	b := []comparator.Result{makeResult("X", comparator.StatusMissing, map[string]string{"prod": "1"})}
	if digester.Equal(a, b) {
		t.Fatal("Equal should return false for different result sets")
	}
}
