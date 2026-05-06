package validator_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/validator"
)

func issueMessages(issues []validator.Issue) []string {
	out := make([]string, len(issues))
	for i, iss := range issues {
		out[i] = iss.String()
	}
	return out
}

func TestValidate_CleanMap(t *testing.T) {
	env := map[string]string{
		"APP_ENV":    "production",
		"DB_HOST":    "localhost",
		"PORT":       "8080",
	}
	issues := validator.Validate(env)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got: %v", issueMessages(issues))
	}
}

func TestValidate_EmptyKey(t *testing.T) {
	env := map[string]string{"": "value"}
	issues := validator.Validate(env)
	if len(issues) == 0 {
		t.Fatal("expected an issue for empty key")
	}
	if !strings.Contains(issues[0].Message, "must not be empty") {
		t.Errorf("unexpected message: %s", issues[0].Message)
	}
}

func TestValidate_InvalidKeyChars(t *testing.T) {
	env := map[string]string{"BAD-KEY": "val"}
	issues := validator.Validate(env)
	if len(issues) == 0 {
		t.Fatal("expected an issue for invalid character in key")
	}
	if !strings.Contains(issues[0].Message, "invalid character") {
		t.Errorf("unexpected message: %s", issues[0].Message)
	}
}

func TestValidate_CaseInsensitiveDuplicate(t *testing.T) {
	// Go maps cannot have truly duplicate keys, so we test the detection
	// by calling Validate twice with overlapping lower-case keys.
	env := map[string]string{
		"db_host": "localhost",
		"DB_HOST": "remotehost",
	}
	issues := validator.Validate(env)
	if len(issues) == 0 {
		t.Fatal("expected a duplicate-key issue")
	}
	found := false
	for _, iss := range issues {
		if strings.Contains(iss.Message, "duplicate") {
			found = true
		}
	}
	if !found {
		t.Errorf("no duplicate issue found; got: %v", issueMessages(issues))
	}
}

func TestValidate_BlankValue(t *testing.T) {
	env := map[string]string{"EMPTY_VAL": "   "}
	issues := validator.Validate(env)
	if len(issues) == 0 {
		t.Fatal("expected a blank-value issue")
	}
	if !strings.Contains(issues[0].Message, "blank") {
		t.Errorf("unexpected message: %s", issues[0].Message)
	}
}

func TestIssue_String(t *testing.T) {
	iss := validator.Issue{Key: "FOO", Message: "something wrong"}
	got := iss.String()
	if !strings.Contains(got, "FOO") || !strings.Contains(got, "something wrong") {
		t.Errorf("unexpected String() output: %s", got)
	}
}
