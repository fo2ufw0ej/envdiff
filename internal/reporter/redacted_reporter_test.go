package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/redactor"
	"github.com/user/envdiff/internal/reporter"
)

func TestRedactedReporter_MasksSensitiveValues(t *testing.T) {
	results := []comparator.Result{
		{
			Key:    "DB_PASSWORD",
			Status: comparator.StatusMismatch,
			Values: map[string]string{
				"dev":  "devpass",
				"prod": "prodpass",
			},
		},
		{
			Key:    "APP_NAME",
			Status: comparator.StatusIdentical,
			Values: map[string]string{
				"dev":  "myapp",
				"prod": "myapp",
			},
		},
	}

	r := redactor.New(nil, "***")
	inner := reporter.NewTextReporter()
	rr := reporter.NewRedactedReporter(inner, r)

	var buf bytes.Buffer
	if err := rr.Report(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "devpass") || strings.Contains(out, "prodpass") {
		t.Error("sensitive values should have been redacted from output")
	}
	if !strings.Contains(out, "***") {
		t.Error("expected mask *** to appear in output")
	}
	if !strings.Contains(out, "myapp") {
		t.Error("non-sensitive value myapp should remain in output")
	}
}

func TestRedactedReporter_NilRedactorUsesDefaults(t *testing.T) {
	results := []comparator.Result{
		{
			Key:    "API_KEY",
			Status: comparator.StatusMismatch,
			Values: map[string]string{"dev": "key1", "prod": "key2"},
		},
	}
	inner := reporter.NewTextReporter()
	rr := reporter.NewRedactedReporter(inner, nil)

	var buf bytes.Buffer
	if err := rr.Report(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "key1") || strings.Contains(out, "key2") {
		t.Error("expected API_KEY values to be redacted")
	}
}

func TestRedactedReporter_MissingKeyRedacted(t *testing.T) {
	results := []comparator.Result{
		{
			Key:    "SECRET_TOKEN",
			Status: comparator.StatusMissing,
			Values: map[string]string{"dev": "tok123", "prod": ""},
		},
	}
	r := redactor.New(nil, "[hidden]")
	inner := reporter.NewTextReporter()
	rr := reporter.NewRedactedReporter(inner, r)

	var buf bytes.Buffer
	if err := rr.Report(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "tok123") {
		t.Error("expected tok123 to be hidden")
	}
}
