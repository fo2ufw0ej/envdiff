package aliaser_test

import (
	"testing"

	"github.com/user/envdiff/internal/aliaser"
	"github.com/user/envdiff/internal/comparator"
)

func makeResult(key, status string) comparator.Result {
	return comparator.Result{
		Key:    key,
		Status: status,
		Values: map[string]string{"dev": "val"},
	}
}

func TestNew_InvalidRule_EmptyFrom(t *testing.T) {
	_, err := aliaser.New([]aliaser.Rule{{From: "", To: "NEW_KEY"}})
	if err == nil {
		t.Fatal("expected error for empty From field")
	}
}

func TestNew_InvalidRule_EmptyTo(t *testing.T) {
	_, err := aliaser.New([]aliaser.Rule{{From: "OLD_KEY", To: ""}})
	if err == nil {
		t.Fatal("expected error for empty To field")
	}
}

func TestApply_NoRules(t *testing.T) {
	a, _ := aliaser.New(nil)
	input := []comparator.Result{makeResult("FOO", "identical")}
	out := a.Apply(input)
	if out[0].Key != "FOO" {
		t.Fatalf("expected FOO, got %s", out[0].Key)
	}
}

func TestApply_RenamesMatchingKey(t *testing.T) {
	a, _ := aliaser.New([]aliaser.Rule{{From: "DATABASE_URL", To: "DB_URL"}})
	input := []comparator.Result{makeResult("DATABASE_URL", "missing")}
	out := a.Apply(input)
	if out[0].Key != "DB_URL" {
		t.Fatalf("expected DB_URL, got %s", out[0].Key)
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	a, _ := aliaser.New([]aliaser.Rule{{From: "database_url", To: "DB_URL"}})
	input := []comparator.Result{makeResult("DATABASE_URL", "missing")}
	out := a.Apply(input)
	if out[0].Key != "DB_URL" {
		t.Fatalf("expected DB_URL, got %s", out[0].Key)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	a, _ := aliaser.New([]aliaser.Rule{{From: "OLD", To: "NEW"}})
	input := []comparator.Result{makeResult("OLD", "missing")}
	a.Apply(input)
	if input[0].Key != "OLD" {
		t.Fatal("Apply must not mutate the input slice")
	}
}

func TestApply_FirstRuleWins(t *testing.T) {
	a, _ := aliaser.New([]aliaser.Rule{
		{From: "OLD", To: "FIRST"},
		{From: "OLD", To: "SECOND"},
	})
	input := []comparator.Result{makeResult("OLD", "missing")}
	out := a.Apply(input)
	if out[0].Key != "FIRST" {
		t.Fatalf("expected FIRST, got %s", out[0].Key)
	}
}
