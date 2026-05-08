package scorer

import (
	"github.com/user/envdiff/internal/comparator"
)

// Score holds the health metrics for a set of compared environments.
type Score struct {
	// Total number of unique keys across all environments.
	TotalKeys int
	// Number of keys that are identical in every environment.
	IdenticalKeys int
	// Number of keys missing from at least one environment.
	MissingKeys int
	// Number of keys present everywhere but with differing values.
	MismatchedKeys int
	// HealthPercent is a 0–100 score: 100 means all keys identical.
	HealthPercent float64
}

// Compute derives a Score from a slice of comparator results.
func Compute(results []comparator.Result) Score {
	if len(results) == 0 {
		return Score{HealthPercent: 100}
	}

	var total, identical, missing, mismatched int

	for _, r := range results {
		total++
		switch r.Status {
		case comparator.StatusIdentical:
			identical++
		case comparator.StatusMissing:
			missing++
		case comparator.StatusMismatched:
			mismatched++
		}
	}

	var health float64
	if total > 0 {
		health = float64(identical) / float64(total) * 100
	}

	return Score{
		TotalKeys:      total,
		IdenticalKeys:  identical,
		MissingKeys:    missing,
		MismatchedKeys: mismatched,
		HealthPercent:  health,
	}
}
