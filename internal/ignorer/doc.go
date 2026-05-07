// Package ignorer provides key-based exclusion for envdiff comparisons.
//
// An Ignorer can be constructed from an explicit list of key names and
// prefix patterns, or loaded from a plain-text ignore file (similar in
// spirit to .gitignore).  Lines beginning with '#' are treated as
// comments; lines ending with '*' are treated as prefix patterns;
// all other non-blank lines are treated as exact key names.
// Matching is always case-insensitive.
package ignorer
