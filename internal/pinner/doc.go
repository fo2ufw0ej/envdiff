// Package pinner provides a mechanism to mark specific .env keys as
// "pinned". Pinned keys are treated as authoritative reference values;
// value mismatches for pinned keys are silently suppressed so that
// intentional per-environment overrides do not pollute comparison
// output. Missing-key findings are still reported for pinned keys
// because absence is always noteworthy.
package pinner
