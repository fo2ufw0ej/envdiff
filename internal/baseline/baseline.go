// Package baseline provides functionality to save and compare env comparison
// results against a previously stored baseline snapshot.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/comparator"
)

// Baseline holds a persisted snapshot of comparison results.
type Baseline struct {
	Entries []comparator.Result `json:"entries"`
}

// Save writes the given results to a JSON file at the specified path.
func Save(path string, results []comparator.Result) error {
	b := Baseline{Entries: results}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write %s: %w", path, err)
	}
	return nil
}

// Load reads a baseline snapshot from the given JSON file.
func Load(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: read %s: %w", path, err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &b, nil
}

// Diff returns entries that are new or changed relative to the baseline.
// An entry is considered new if its key+env combination did not exist before,
// and changed if its Status differs from the baseline entry.
func Diff(base *Baseline, current []comparator.Result) []comparator.Result {
	index := make(map[string]comparator.Result, len(base.Entries))
	for _, e := range base.Entries {
		index[entryKey(e)] = e
	}

	var delta []comparator.Result
	for _, e := range current {
		prev, found := index[entryKey(e)]
		if !found || prev.Status != e.Status {
			delta = append(delta, e)
		}
	}
	return delta
}

func entryKey(r comparator.Result) string {
	return r.Key + "\x00" + r.Env
}
