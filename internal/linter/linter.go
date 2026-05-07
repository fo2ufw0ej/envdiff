package linter

import (
	"fmt"
	"strings"
)

// Severity indicates how serious a lint finding is.
type Severity string

const (
	SeverityWarn  Severity = "warn"
	SeverityError Severity = "error"
)

// Finding represents a single lint issue found in an env map.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

// Rule is a function that inspects a key/value pair and returns findings.
type Rule func(key, value string) []Finding

// Linter runs a set of rules over an env map.
type Linter struct {
	rules []Rule
}

// New creates a Linter with the default built-in rules.
func New() *Linter {
	return &Linter{
		rules: []Rule{
			RuleNoEmptyValue,
			RuleNoWhitespaceInKey,
			RuleUpperCaseKey,
		},
	}
}

// NewWithRules creates a Linter with a custom set of rules.
func NewWithRules(rules ...Rule) *Linter {
	return &Linter{rules: rules}
}

// Lint runs all rules over the provided env map and returns all findings.
func (l *Linter) Lint(env map[string]string) []Finding {
	var findings []Finding
	for k, v := range env {
		for _, rule := range l.rules {
			findings = append(findings, rule(k, v)...)
		}
	}
	return findings
}

// RuleNoEmptyValue warns when a key has an empty value.
func RuleNoEmptyValue(key, value string) []Finding {
	if strings.TrimSpace(value) == "" {
		return []Finding{{Key: key, Message: "value is empty", Severity: SeverityWarn}}
	}
	return nil
}

// RuleNoWhitespaceInKey errors when a key contains whitespace.
func RuleNoWhitespaceInKey(key, _ string) []Finding {
	if strings.ContainsAny(key, " \t") {
		return []Finding{{
			Key:      key,
			Message:  fmt.Sprintf("key %q contains whitespace", key),
			Severity: SeverityError,
		}}
	}
	return nil
}

// RuleUpperCaseKey warns when a key contains lower-case letters.
func RuleUpperCaseKey(key, _ string) []Finding {
	if key != strings.ToUpper(key) {
		return []Finding{{
			Key:      key,
			Message:  fmt.Sprintf("key %q is not upper-case", key),
			Severity: SeverityWarn,
		}}
	}
	return nil
}
