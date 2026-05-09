package renamer_test

import (
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/renamer"
)

func makeResult(key, status string, values map[string]string) comparator.Result {
	return comparator.Result{Key: key, Status: status, Values: values}
}

func TestApply_NoRules(t *testing.T) {
	input := []comparator.Result{makeResult("DB_HOST", "identical", nil)}
	out, rr, err := renamer.Apply(input, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out[0].Key != "DB_HOST" {
		t.Fatalf("expected unchanged result, got %+v", out)
	}
	if len(rr) != 0 {
		t.Fatalf("expected no rule results")
	}
}

func TestApply_RenamesMatchingKey(t *testing.T) {
	input := []comparator.Result{
		makeResult("OLD_KEY", "missing", map[string]string{"prod": ""}),
		makeResult("KEEP", "identical", map[string]string{"prod": "val"}),
	}
	rules := []renamer.Rule{{OldKey: "OLD_KEY", NewKey: "NEW_KEY"}}
	out, rr, err := renamer.Apply(input, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Key != "NEW_KEY" {
		t.Errorf("expected NEW_KEY, got %q", out[0].Key)
	}
	if out[1].Key != "KEEP" {
		t.Errorf("expected KEEP unchanged, got %q", out[1].Key)
	}
	if !rr[0].Matched {
		t.Errorf("expected rule to be marked as matched")
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	input := []comparator.Result{makeResult("db_host", "identical", nil)}
	rules := []renamer.Rule{{OldKey: "DB_HOST", NewKey: "DATABASE_HOST"}}
	out, _, err := renamer.Apply(input, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Key != "DATABASE_HOST" {
		t.Errorf("expected DATABASE_HOST, got %q", out[0].Key)
	}
}

func TestApply_FirstRuleWins(t *testing.T) {
	input := []comparator.Result{makeResult("FOO", "identical", nil)}
	rules := []renamer.Rule{
		{OldKey: "FOO", NewKey: "BAR"},
		{OldKey: "FOO", NewKey: "BAZ"},
	}
	_, _, err := renamer.Apply(input, rules)
	if err == nil {
		t.Fatal("expected duplicate key error")
	}
}

func TestApply_EmptyOldKeyError(t *testing.T) {
	_, _, err := renamer.Apply(nil, []renamer.Rule{{OldKey: "", NewKey: "X"}})
	if err == nil {
		t.Fatal("expected error for empty OldKey")
	}
}

func TestApply_UnmatchedRuleNotMarked(t *testing.T) {
	input := []comparator.Result{makeResult("OTHER", "identical", nil)}
	rules := []renamer.Rule{{OldKey: "MISSING", NewKey: "RENAMED"}}
	_, rr, err := renamer.Apply(input, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rr[0].Matched {
		t.Errorf("expected rule to be unmatched")
	}
}
