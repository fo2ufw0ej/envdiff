// Package trimmer provides utilities for reducing a slice of comparison
// results based on key allowlists, blocklists, and value-length constraints.
//
// Typical usage:
//
//	results := comparator.Compare(envs)
//	trimmed := trimmer.Trim(results, trimmer.Options{
//		BlockKeys:   []string{"INTERNAL_TOKEN", "CI_JOB_TOKEN"},
//		MaxValueLen: 256,
//	})
//
// Trimming is applied before reporting so that sensitive or irrelevant keys
// do not appear in output.
package trimmer
