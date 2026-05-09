// Package patcher merges key-value pairs from a source env map into a
// destination env map and can persist the result to disk.
//
// Two strategies are supported:
//
//   - StrategyAddMissing – only keys that are absent from the destination are
//     inserted; existing values are never touched.
//
//   - StrategyOverwrite – missing keys are inserted AND keys whose values
//     differ from the source are updated to the source value.
//
// Every key processed by Patch is recorded in a Result slice so callers can
// audit or report exactly what changed.
package patcher
