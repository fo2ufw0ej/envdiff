package tagger_test

import (
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/tagger"
)

func makeResult(key string) comparator.Result {
	return comparator.Result{
		Key:    key,
		Status: comparator.StatusIdentical,
		Values: map[string]string{"prod": "val"},
	}
}

func tagNames(tr tagger.TaggedResult) []string {
	return tr.Tags
}

func TestApply_NoRules(t *testing.T) {
	tgr := tagger.New(nil)
	results := []comparator.Result{makeResult("DB_HOST")}
	out := tgr.Apply(results)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if len(out[0].Tags) != 0 {
		t.Errorf("expected no tags, got %v", out[0].Tags)
	}
}

func TestApply_ExactMatch(t *testing.T) {
	tgr := tagger.New([]tagger.Rule{
		{Pattern: "DB_HOST", Tags: []string{"database"}},
	})
	out := tgr.Apply([]comparator.Result{makeResult("DB_HOST")})
	if len(out[0].Tags) != 1 || out[0].Tags[0] != "database" {
		t.Errorf("unexpected tags: %v", out[0].Tags)
	}
}

func TestApply_PrefixMatch(t *testing.T) {
	tgr := tagger.New([]tagger.Rule{
		{Pattern: "DB_*", Tags: []string{"database"}},
	})
	results := []comparator.Result{
		makeResult("DB_HOST"),
		makeResult("DB_PORT"),
		makeResult("APP_NAME"),
	}
	out := tgr.Apply(results)
	if len(out[0].Tags) == 0 || out[0].Tags[0] != "database" {
		t.Errorf("DB_HOST should be tagged: %v", out[0].Tags)
	}
	if len(out[1].Tags) == 0 || out[1].Tags[0] != "database" {
		t.Errorf("DB_PORT should be tagged: %v", out[1].Tags)
	}
	if len(out[2].Tags) != 0 {
		t.Errorf("APP_NAME should not be tagged: %v", out[2].Tags)
	}
}

func TestApply_MultipleRulesUnionTags(t *testing.T) {
	tgr := tagger.New([]tagger.Rule{
		{Pattern: "DB_*", Tags: []string{"database"}},
		{Pattern: "DB_HOST", Tags: []string{"network", "database"}},
	})
	out := tgr.Apply([]comparator.Result{makeResult("DB_HOST")})
	// should have "database" and "network" without duplicates
	if len(out[0].Tags) != 2 {
		t.Errorf("expected 2 unique tags, got %v", out[0].Tags)
	}
}

func TestApply_CaseInsensitiveExact(t *testing.T) {
	tgr := tagger.New([]tagger.Rule{
		{Pattern: "db_host", Tags: []string{"database"}},
	})
	out := tgr.Apply([]comparator.Result{makeResult("DB_HOST")})
	if len(out[0].Tags) == 0 {
		t.Error("expected case-insensitive match to produce tag")
	}
}

func TestNew_IgnoresInvalidRules(t *testing.T) {
	tgr := tagger.New([]tagger.Rule{
		{Pattern: "", Tags: []string{"orphan"}},
		{Pattern: "DB_HOST", Tags: nil},
		{Pattern: "VALID", Tags: []string{"ok"}},
	})
	out := tgr.Apply([]comparator.Result{makeResult("VALID")})
	if len(out[0].Tags) != 1 || out[0].Tags[0] != "ok" {
		t.Errorf("unexpected tags: %v", out[0].Tags)
	}
}
