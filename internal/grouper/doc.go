// Package grouper partitions a slice of comparator.Result values into named
// groups along a chosen dimension.
//
// Three grouping strategies are supported:
//
//   - ByPrefix  — groups by the first underscore-delimited segment of the key
//     (e.g. "DB_HOST" and "DB_PORT" both fall into the "DB" group).
//   - ByStatus  — groups by comparison status (match, missing, mismatch).
//   - ByEnv     — groups by the first environment name present in the result's
//     Values map.
//
// Groups are always returned in alphabetical order so that callers receive
// deterministic output regardless of the input ordering.
package grouper
