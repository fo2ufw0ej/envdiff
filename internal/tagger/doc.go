// Package tagger attaches free-form string tags to comparison results based
// on configurable rules. Rules may match keys exactly or by prefix (using a
// trailing "*" wildcard). Multiple rules may fire for a single key; their
// tags are merged into a deduplicated slice.
//
// Typical usage:
//
//	tr := tagger.New([]tagger.Rule{
//		{Pattern: "DB_*",   Tags: []string{"database"}},
//		{Pattern: "AWS_*",  Tags: []string{"cloud", "sensitive"}},
//	})
//	tagged := tr.Apply(results)
package tagger
