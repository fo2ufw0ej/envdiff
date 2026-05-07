// Package linter provides a rule-based linter for environment variable maps.
//
// It ships with built-in rules that flag common problems such as empty values,
// keys containing whitespace, and keys that are not fully upper-case.  Custom
// rules can be supplied via NewWithRules to extend or replace the defaults.
//
// Example usage:
//
//	l := linter.New()
//	findings := l.Lint(envMap)
//	for _, f := range findings {
//		fmt.Printf("[%s] %s: %s\n", f.Severity, f.Key, f.Message)
//	}
package linter
