package promoter_test

import (
	"testing"

	"github.com/user/envdiff/internal/promoter"
)

func resultMap(results []promoter.Result) map[string]promoter.Result {
	m := make(map[string]promoter.Result, len(results))
	for _, r := range results {
		m[r.Key] = r
	}
	return m
}

func TestPromote_NilSrcReturnsError(t *testing.T) {
	_, _, err := promoter.Promote(nil, map[string]string{}, promoter.SkipExisting)
	if err == nil {
		t.Fatal("expected error for nil src, got nil")
	}
}

func TestPromote_NilDstTreatedAsEmpty(t *testing.T) {
	src := map[string]string{"KEY": "val"}
	out, results, err := promoter.Promote(src, nil, promoter.SkipExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %q", out["KEY"])
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestPromote_SkipExisting(t *testing.T) {
	src := map[string]string{"A": "new", "B": "bval"}
	dst := map[string]string{"A": "old"}

	out, results, err := promoter.Promote(src, dst, promoter.SkipExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "old" {
		t.Errorf("expected A to remain 'old', got %q", out["A"])
	}
	if out["B"] != "bval" {
		t.Errorf("expected B=bval, got %q", out["B"])
	}

	rm := resultMap(results)
	if !rm["A"].Skipped {
		t.Error("expected A to be skipped")
	}
	if rm["B"].Skipped {
		t.Error("expected B not to be skipped")
	}
}

func TestPromote_OverwriteExisting(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}

	out, results, err := promoter.Promote(src, dst, promoter.OverwriteExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "new" {
		t.Errorf("expected A=new, got %q", out["A"])
	}

	rm := resultMap(results)
	if rm["A"].Skipped {
		t.Error("expected A not to be skipped")
	}
	if rm["A"].OldValue != "old" {
		t.Errorf("expected OldValue=old, got %q", rm["A"].OldValue)
	}
}

func TestPromote_DoesNotMutateDst(t *testing.T) {
	src := map[string]string{"X": "1"}
	dst := map[string]string{"Y": "2"}

	_, _, _ = promoter.Promote(src, dst, promoter.OverwriteExisting)
	if _, ok := dst["X"]; ok {
		t.Error("Promote must not mutate the dst map")
	}
}
