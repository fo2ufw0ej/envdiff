package reporter

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/user/envdiff/internal/comparator"
)

func TestJSONReporter_NoDifferences(t *testing.T) {
	result := comparator.Result{}
	var buf bytes.Buffer
	r := NewJSONReporter(&buf)
	if err := r.Report(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["has_differences"].(bool) {
		t.Error("expected has_differences=false")
	}
}

func TestJSONReporter_MissingKey(t *testing.T) {
	result := comparator.Result{
		Missing: map[string]comparator.MissingEntry{
			"DB_HOST": {
				PresentIn: []string{"production"},
				AbsentIn:  []string{"staging"},
			},
		},
	}
	var buf bytes.Buffer
	r := NewJSONReporter(&buf)
	if err := r.Report(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !out["has_differences"].(bool) {
		t.Error("expected has_differences=true")
	}
	missing := out["missing"].([]interface{})
	if len(missing) != 1 {
		t.Fatalf("expected 1 missing entry, got %d", len(missing))
	}
	entry := missing[0].(map[string]interface{})
	if entry["key"].(string) != "DB_HOST" {
		t.Errorf("expected key DB_HOST, got %s", entry["key"])
	}
}

func TestJSONReporter_MismatchedKey(t *testing.T) {
	result := comparator.Result{
		Mismatched: map[string]map[string]string{
			"LOG_LEVEL": {
				"production": "error",
				"staging":    "debug",
			},
		},
	}
	var buf bytes.Buffer
	r := NewJSONReporter(&buf)
	if err := r.Report(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	mismatched := out["mismatched"].([]interface{})
	if len(mismatched) != 1 {
		t.Fatalf("expected 1 mismatched entry, got %d", len(mismatched))
	}
	entry := mismatched[0].(map[string]interface{})
	if entry["key"].(string) != "LOG_LEVEL" {
		t.Errorf("expected key LOG_LEVEL, got %s", entry["key"])
	}
	values := entry["values"].(map[string]interface{})
	if values["production"].(string) != "error" {
		t.Errorf("unexpected production value: %v", values["production"])
	}
}
