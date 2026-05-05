// Package comparator implements the core diffing logic for envdiff.
//
// It accepts parsed environment maps (produced by the parser package) and
// computes two categories of discrepancies:
//
//  1. Missing keys — a key present in at least one environment but absent
//     from one or more others.
//
//  2. Mismatched values — a key present in all environments but whose value
//     differs between at least two of them.
//
// Usage:
//
//	envs := map[string]map[string]string{
//	    "dev":  parser.ParseFile("dev.env"),
//	    "prod": parser.ParseFile("prod.env"),
//	}
//	result := comparator.Compare(envs)
//
The Result type exposes MissingIn and Mismatched fields that downstream
// reporters can use to render human-readable output.
package comparator
