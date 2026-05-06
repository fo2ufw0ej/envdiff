// Package differ computes value-level diffs between environment variable
// entries found in two environments. It uses a longest-common-subsequence
// algorithm to produce unified-style line diffs, which is especially
// useful when values span multiple lines (e.g. PEM certificates or
// multi-line JSON blobs stored in env vars).
//
// Basic usage:
//
//	diff := differ.DiffValues("DATABASE_URL", leftVal, rightVal)
//	fmt.Print(diff.Format())
package differ
