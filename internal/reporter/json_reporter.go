package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/comparator"
)

// JSONReporter writes comparison results as a JSON document.
type JSONReporter struct {
	w io.Writer
}

// NewJSONReporter returns a JSONReporter that writes to w.
func NewJSONReporter(w io.Writer) *JSONReporter {
	return &JSONReporter{w: w}
}

type jsonReport struct {
	HasDifferences bool            `json:"has_differences"`
	Missing        []jsonMissing   `json:"missing,omitempty"`
	Mismatched     []jsonMismatch  `json:"mismatched,omitempty"`
}

type jsonMissing struct {
	Key     string   `json:"key"`
	PresentIn []string `json:"present_in"`
	AbsentIn  []string `json:"absent_in"`
}

type jsonMismatch struct {
	Key    string            `json:"key"`
	Values map[string]string `json:"values"`
}

// Report serialises the comparison result to JSON and writes it to the reporter's writer.
func (r *JSONReporter) Report(result comparator.Result) error {
	report := jsonReport{
		HasDifferences: result.HasDifferences(),
	}

	// Collect missing keys in deterministic order.
	keys := make([]string, 0, len(result.Missing))
	for k := range result.Missing {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		entry := result.Missing[k]
		report.Missing = append(report.Missing, jsonMissing{
			Key:       k,
			PresentIn: sorted(entry.PresentIn),
			AbsentIn:  sorted(entry.AbsentIn),
		})
	}

	// Collect mismatched keys in deterministic order.
	mkeys := make([]string, 0, len(result.Mismatched))
	for k := range result.Mismatched {
		mkeys = append(mkeys, k)
	}
	sort.Strings(mkeys)

	for _, k := range mkeys {
		report.Mismatched = append(report.Mismatched, jsonMismatch{
			Key:    k,
			Values: result.Mismatched[k],
		})
	}

	enc := json.NewEncoder(r.w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		return fmt.Errorf("json reporter: %w", err)
	}
	return nil
}

func sorted(in []string) []string {
	out := make([]string, len(in))
	copy(out, in)
	sort.Strings(out)
	return out
}
