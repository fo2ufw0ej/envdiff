// Package profiler provides statistical analysis of one or more parsed
// environment maps.
//
// Given a collection of env maps (keyed by environment name), Profile returns
// a Report describing:
//
//   - The total number of distinct keys across all environments.
//   - The total number of environments analysed.
//   - Per-key statistics: in how many environments the key appears, how many
//     distinct values it carries, and the maximum value length observed.
//
// This is useful for understanding coverage gaps and spotting keys whose
// values vary unexpectedly across environments.
package profiler
