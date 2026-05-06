// Package exporter provides functionality to export comparison results
// to various output formats such as shell-sourceable env files or JSON patches.
package exporter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/comparator"
)

// Format represents an export output format.
type Format string

const (
	FormatEnv  Format = "env"
	FormatPatch Format = "patch"
)

// Exporter writes missing or mismatched keys from a reference environment
// into a target-compatible format.
type Exporter struct {
	format    Format
	reference string
}

// New creates a new Exporter for the given format and reference environment name.
func New(format Format, reference string) *Exporter {
	return &Exporter{format: format, reference: reference}
}

// Export writes the export output to w based on the comparison results.
// Only keys missing or mismatched relative to the reference env are included.
func (e *Exporter) Export(w io.Writer, results []comparator.Result) error {
	switch e.format {
	case FormatEnv:
		return e.writeEnv(w, results)
	case FormatPatch:
		return e.writePatch(w, results)
	default:
		return fmt.Errorf("unsupported export format: %s", e.format)
	}
}

func (e *Exporter) writeEnv(w io.Writer, results []comparator.Result) error {
	keys := sortedKeys(results)
	for _, key := range keys {
		for _, r := range results {
			if r.Key != key {
				continue
			}
			val, ok := r.Values[e.reference]
			if !ok || val == "" {
				continue
			}
			if _, err := fmt.Fprintf(w, "%s=%s\n", key, shellQuote(val)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Exporter) writePatch(w io.Writer, results []comparator.Result) error {
	keys := sortedKeys(results)
	for _, key := range keys {
		for _, r := range results {
			if r.Key != key {
				continue
			}
			val, ok := r.Values[e.reference]
			if !ok {
				continue
			}
			if _, err := fmt.Fprintf(w, "+ %s=%s\n", key, val); err != nil {
				return err
			}
		}
	}
	return nil
}

func sortedKeys(results []comparator.Result) []string {
	seen := map[string]struct{}{}
	for _, r := range results {
		seen[r.Key] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func shellQuote(s string) string {
	if !strings.ContainsAny(s, " \t\n'\"\\$`") {
		return s
	}
	return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
}
