package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envdiff/internal/comparator"
)

// Snapshot captures the state of a comparison at a point in time.
type Snapshot struct {
	CreatedAt time.Time                  `json:"created_at"`
	Label     string                     `json:"label"`
	Results   []comparator.Result        `json:"results"`
}

// Save writes a snapshot to the given file path as JSON.
func Save(path, label string, results []comparator.Result) error {
	s := Snapshot{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Results:   results,
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read %s: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &s, nil
}

// Compare returns entries whose status differs between two snapshots.
func Compare(old, new *Snapshot) []Delta {
	oldIndex := indexResults(old.Results)
	var deltas []Delta
	for _, r := range new.Results {
		key := entryKey(r)
		if prev, ok := oldIndex[key]; ok {
			if prev.Status != r.Status {
				deltas = append(deltas, Delta{Old: prev, New: r})
			}
		} else {
			deltas = append(deltas, Delta{New: r, Added: true})
		}
	}
	return deltas
}

// Delta describes a change between two snapshots for a single result entry.
type Delta struct {
	Old   comparator.Result
	New   comparator.Result
	Added bool
}

func indexResults(results []comparator.Result) map[string]comparator.Result {
	m := make(map[string]comparator.Result, len(results))
	for _, r := range results {
		m[entryKey(r)] = r
	}
	return m
}

func entryKey(r comparator.Result) string {
	return r.Key + "::" + r.Env
}
