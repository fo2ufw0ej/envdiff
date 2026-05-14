// Package archiver saves and retrieves historical snapshots of comparison
// results, allowing users to track how their .env files evolve over time.
package archiver

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/user/envdiff/internal/comparator"
)

// Entry is a single archived run.
type Entry struct {
	Timestamp time.Time            `json:"timestamp"`
	Label     string               `json:"label,omitempty"`
	Results   []comparator.Result  `json:"results"`
}

// Archive holds a collection of historical entries persisted to a directory.
type Archive struct {
	dir string
}

// New returns an Archive that stores entries in dir.
// The directory is created if it does not exist.
func New(dir string) (*Archive, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("archiver: create dir: %w", err)
	}
	return &Archive{dir: dir}, nil
}

// Save writes results to a timestamped JSON file inside the archive directory.
func (a *Archive) Save(results []comparator.Result, label string) (string, error) {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Results:   results,
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return "", fmt.Errorf("archiver: marshal: %w", err)
	}
	name := fmt.Sprintf("%d.json", entry.Timestamp.UnixNano())
	path := filepath.Join(a.dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", fmt.Errorf("archiver: write: %w", err)
	}
	return path, nil
}

// List returns all archived entries sorted from oldest to newest.
func (a *Archive) List() ([]Entry, error) {
	glob := filepath.Join(a.dir, "*.json")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, fmt.Errorf("archiver: glob: %w", err)
	}
	sort.Strings(matches)
	var entries []Entry
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			return nil, fmt.Errorf("archiver: read %s: %w", m, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("archiver: unmarshal %s: %w", m, err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// Latest returns the most recently saved entry, or an error if none exist.
func (a *Archive) Latest() (Entry, error) {
	entries, err := a.List()
	if err != nil {
		return Entry{}, err
	}
	if len(entries) == 0 {
		return Entry{}, fmt.Errorf("archiver: no entries found")
	}
	return entries[len(entries)-1], nil
}
