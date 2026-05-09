package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/aliaser"
	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/reporter"
)

func TestAliasReporter_NoRules(t *testing.T) {
	inner := reporter.NewTextReporter()
	ar, err := reporter.NewAliasReporter(inner, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var buf bytes.Buffer
	results := []comparator.Result{}
	if err := ar.Report(&buf, results); err != nil {
		t.Fatalf("Report error: %v", err)
	}
	if strings.Contains(buf.String(), "Active key aliases") {
		t.Error("expected no alias preamble when rules are empty")
	}
}

func TestAliasReporter_PrintsPreamble(t *testing.T) {
	inner := reporter.NewTextReporter()
	rules := []aliaser.Rule{{From: "DATABASE_URL", To: "DB_URL"}}
	ar, _ := reporter.NewAliasReporter(inner, rules)
	var buf bytes.Buffer
	if err := ar.Report(&buf, []comparator.Result{}); err != nil {
		t.Fatalf("Report error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "DATABASE_URL -> DB_URL") {
		t.Errorf("expected preamble with alias rule, got:\n%s", got)
	}
}

func TestAliasReporter_RewritesKey(t *testing.T) {
	inner := reporter.NewTextReporter()
	rules := []aliaser.Rule{{From: "OLD_KEY", To: "NEW_KEY"}}
	ar, _ := reporter.NewAliasReporter(inner, rules)
	results := []comparator.Result{
		{Key: "OLD_KEY", Status: "missing", Values: map[string]string{"dev": "val"}},
	}
	var buf bytes.Buffer
	if err := ar.Report(&buf, results); err != nil {
		t.Fatalf("Report error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "NEW_KEY") {
		t.Errorf("expected rewritten key NEW_KEY in output, got:\n%s", got)
	}
	if strings.Contains(got, "OLD_KEY") {
		t.Errorf("expected OLD_KEY to be replaced, got:\n%s", got)
	}
}

func TestAliasReporter_InvalidRule(t *testing.T) {
	inner := reporter.NewTextReporter()
	_, err := reporter.NewAliasReporter(inner, []aliaser.Rule{{From: "", To: "X"}})
	if err == nil {
		t.Fatal("expected error for invalid alias rule")
	}
}
