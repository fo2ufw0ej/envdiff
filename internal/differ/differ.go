// Package differ provides utilities for computing a line-level diff
// between two environment variable values, making it easier to spot
// subtle differences such as trailing spaces or casing changes.
package differ

import (
	"fmt"
	"strings"
)

// LineDiff represents a single line in a unified-style diff.
type LineDiff struct {
	Op   rune   // ' ' for context, '-' for removed, '+' for added
	Text string // the value fragment
}

// ValueDiff holds the diff result between two scalar values.
type ValueDiff struct {
	Key   string
	Left  string
	Right string
	Lines []LineDiff
}

// DiffValues computes a simple line-by-line diff between two values
// associated with the given key. For single-line values the diff
// degenerates to a removal + addition pair.
func DiffValues(key, left, right string) ValueDiff {
	if left == right {
		return ValueDiff{Key: key, Left: left, Right: right,
			Lines: []LineDiff{{Op: ' ', Text: left}}}
	}

	leftLines := splitLines(left)
	rightLines := splitLines(right)

	lines := lcs(leftLines, rightLines)
	return ValueDiff{Key: key, Left: left, Right: right, Lines: lines}
}

// Format returns a human-readable representation of the diff.
func (d ValueDiff) Format() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "key: %s\n", d.Key)
	for _, l := range d.Lines {
		fmt.Fprintf(&sb, "%c %s\n", l.Op, l.Text)
	}
	return sb.String()
}

func splitLines(s string) []string {
	if s == "" {
		return []string{""}
	}
	return strings.Split(s, "\n")
}

// lcs builds a diff using a simple longest-common-subsequence approach.
func lcs(a, b []string) []LineDiff {
	n, m := len(a), len(b)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if a[i] == b[j] {
				dp[i][j] = dp[i+1][j+1] + 1
			} else if dp[i+1][j] >= dp[i][j+1] {
				dp[i][j] = dp[i+1][j]
			} else {
				dp[i][j] = dp[i][j+1]
			}
		}
	}

	var result []LineDiff
	i, j := 0, 0
	for i < n && j < m {
		if a[i] == b[j] {
			result = append(result, LineDiff{Op: ' ', Text: a[i]})
			i++
			j++
		} else if dp[i+1][j] >= dp[i][j+1] {
			result = append(result, LineDiff{Op: '-', Text: a[i]})
			i++
		} else {
			result = append(result, LineDiff{Op: '+', Text: b[j]})
			j++
		}
	}
	for ; i < n; i++ {
		result = append(result, LineDiff{Op: '-', Text: a[i]})
	}
	for ; j < m; j++ {
		result = append(result, LineDiff{Op: '+', Text: b[j]})
	}
	return result
}
