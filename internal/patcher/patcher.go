package patcher

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Strategy controls how missing keys are patched into a target env map.
type Strategy int

const (
	// StrategyAddMissing adds only keys absent from the target.
	StrategyAddMissing Strategy = iota
	// StrategyOverwrite adds missing keys and overwrites mismatched values.
	StrategyOverwrite
)

// Result describes a single patch action applied to a key.
type Result struct {
	Key    string
	OldVal string // empty string when key was absent
	NewVal string
	Action string // "added" | "overwritten" | "skipped"
}

// Patch applies values from src into dst according to the chosen strategy.
// It returns the patched map and a log of every action taken.
func Patch(dst, src map[string]string, strategy Strategy) (map[string]string, []Result) {
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	keys := sortedKeys(src)
	var results []Result

	for _, k := range keys {
		srcVal := src[k]
		dstVal, exists := out[k]

		switch {
		case !exists:
			out[k] = srcVal
			results = append(results, Result{Key: k, OldVal: "", NewVal: srcVal, Action: "added"})
		case strategy == StrategyOverwrite && dstVal != srcVal:
			out[k] = srcVal
			results = append(results, Result{Key: k, OldVal: dstVal, NewVal: srcVal, Action: "overwritten"})
		default:
			results = append(results, Result{Key: k, OldVal: dstVal, NewVal: dstVal, Action: "skipped"})
		}
	}

	return out, results
}

// WriteFile serialises patched as KEY=VALUE lines into path, creating or
// truncating the file. Values containing spaces are double-quoted.
func WriteFile(path string, patched map[string]string) error {
	keys := sortedKeys(patched)
	var sb strings.Builder
	for _, k := range keys {
		v := patched[k]
		if strings.ContainsAny(v, " \t") {
			v = fmt.Sprintf("%q", v)
		}
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(v)
		sb.WriteByte('\n')
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
