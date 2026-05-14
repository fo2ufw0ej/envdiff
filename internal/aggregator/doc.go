// Package aggregator collects and summarises comparator results across
// all environments, producing per-environment statistics (missing,
// mismatched, identical counts) as well as totals across the whole
// comparison run.
//
// Usage:
//
//	results := comparator.Compare(envs)
//	report  := aggregator.Aggregate(results)
//	for _, e := range report.Envs {
//		fmt.Printf("%s: %d missing, %d mismatched\n", e.Env, e.Missing, e.Mismatched)
//	}
package aggregator
