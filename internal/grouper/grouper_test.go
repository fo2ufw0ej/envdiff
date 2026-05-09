package grouper_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/comparator"
	"github.com/yourorg/envdiff/internal/grouper"
)

func makeResult(key string, status comparator.Status, values map[string]string) comparator.Result {
	return comparator.Result{Key: key, Status: status, Values: values}
}

func groupNames(groups []grouper.Group) []string {
	names := make([]string, len(groups))
	for i, g := range groups {
		names[i] = g.Name
	}
	return names
}

func TestGroupBy_Prefix(t *testing.T) {
	results := []comparator.Result{
		makeResult("DB_HOST", comparator.StatusMissing, nil),
		makeResult("DB_PORT", comparator.StatusMatch, nil),
		makeResult("APP_ENV", comparator.StatusMismatch, nil),
		makeResult("NOUNDERSCORE", comparator.StatusMatch, nil),
	}

	groups := grouper.GroupBy(results, grouper.ByPrefix)
	names := groupNames(groups)

	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
	if names[0] != "APP" || names[1] != "DB" || names[2] != "NOUNDERSCORE" {
		t.Errorf("unexpected group names: %v", names)
	}
	if len(groups[1].Results) != 2 {
		t.Errorf("expected 2 results in DB group, got %d", len(groups[1].Results))
	}
}

func TestGroupBy_Status(t *testing.T) {
	results := []comparator.Result{
		makeResult("KEY_A", comparator.StatusMissing, nil),
		makeResult("KEY_B", comparator.StatusMismatch, nil),
		makeResult("KEY_C", comparator.StatusMissing, nil),
	}

	groups := grouper.GroupBy(results, grouper.ByStatus)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	names := groupNames(groups)
	if names[0] != string(comparator.StatusMismatch) || names[1] != string(comparator.StatusMissing) {
		t.Errorf("unexpected group names: %v", names)
	}
	if len(groups[1].Results) != 2 {
		t.Errorf("expected 2 missing results, got %d", len(groups[1].Results))
	}
}

func TestGroupBy_Env(t *testing.T) {
	results := []comparator.Result{
		makeResult("KEY_A", comparator.StatusMissing, map[string]string{"production": "val"}),
		makeResult("KEY_B", comparator.StatusMismatch, map[string]string{"staging": "val"}),
		makeResult("KEY_C", comparator.StatusMissing, map[string]string{"production": "val"}),
	}

	groups := grouper.GroupBy(results, grouper.ByEnv)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	names := groupNames(groups)
	if names[0] != "production" || names[1] != "staging" {
		t.Errorf("unexpected group names: %v", names)
	}
}

func TestGroupBy_Empty(t *testing.T) {
	groups := grouper.GroupBy(nil, grouper.ByPrefix)
	if len(groups) != 0 {
		t.Errorf("expected no groups for empty input, got %d", len(groups))
	}
}
