// Package renamer provides utilities for renaming keys across comparator
// results. This is useful when migrating environment variable names between
// releases without losing the context of existing comparison data.
//
// Usage:
//
//	rules := []renamer.Rule{
//		{OldKey: "DB_URL", NewKey: "DATABASE_URL"},
//	}
//	updated, ruleResults, err := renamer.Apply(results, rules)
//
// Each Rule maps an OldKey to a NewKey. Matching is case-insensitive so that
// environment files with mixed casing are handled gracefully. Duplicate OldKey
// entries within the rule slice are rejected with an error.
package renamer
