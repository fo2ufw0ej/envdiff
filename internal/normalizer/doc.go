// Package normalizer transforms comparator results before reporting or
// further analysis. It supports key normalisation (upper-casing, prefix
// stripping) and value normalisation (whitespace trimming) through a
// functional-option API.
//
// Example:
//
//	results := comparator.Compare(envs)
//	clean := normalizer.Normalize(results,
//		normalizer.WithUpperKeys(),
//		normalizer.WithTrimValues(),
//		normalizer.WithStripPrefix("APP_"),
//	)
package normalizer
