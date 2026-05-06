package reporter

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/comparator"
)

// csvReporter writes comparison results as CSV output.
type csvReporter struct {
	w io.Writer
}

// NewCSVReporter returns a Reporter that formats results as CSV.
// Columns: key, status, env_name, value
func NewCSVReporter(w io.Writer) Reporter {
	return &csvReporter{w: w}
}

func (r *csvReporter) Report(results []comparator.Result) error {
	cw := csv.NewWriter(r.w)

	// Write header
	if err := cw.Write([]string{"key", "status", "env", "value"}); err != nil {
		return fmt.Errorf("csv reporter: write header: %w", err)
	}

	if len(results) == 0 {
		cw.Flush()
		return cw.Error()
	}

	for _, res := range results {
		switch res.Status {
		case comparator.StatusMissing:
			// One row per env that is missing the key
			envNames := sortedEnvNames(res.Values)
			for _, env := range envNames {
				val, present := res.Values[env]
				status := "missing"
				if present {
					status = "present"
				}
				if err := cw.Write([]string{res.Key, status, env, val}); err != nil {
					return fmt.Errorf("csv reporter: write row: %w", err)
				}
			}
		case comparator.StatusMismatch:
			envNames := sortedEnvNames(res.Values)
			for _, env := range envNames {
				if err := cw.Write([]string{res.Key, "mismatch", env, res.Values[env]}); err != nil {
					return fmt.Errorf("csv reporter: write row: %w", err)
				}
			}
		default:
			envNames := sortedEnvNames(res.Values)
			for _, env := range envNames {
				if err := cw.Write([]string{res.Key, "ok", env, res.Values[env]}); err != nil {
					return fmt.Errorf("csv reporter: write row: %w", err)
				}
			}
		}
	}

	cw.Flush()
	return cw.Error()
}

func sortedEnvNames(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
