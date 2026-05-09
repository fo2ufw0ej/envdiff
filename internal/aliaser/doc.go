// Package aliaser rewrites comparator result keys according to a set of
// user-defined alias rules.
//
// This is useful when the same logical configuration value is stored under
// different key names across environments (e.g. DATABASE_URL in production
// and DB_URL in development). By aliasing the keys before comparison or
// reporting, envdiff can treat them as equivalent.
//
// Usage:
//
//	a, err := aliaser.New([]aliaser.Rule{
//		{From: "DATABASE_URL", To: "DB_URL"},
//	})
//	results = a.Apply(results)
package aliaser
