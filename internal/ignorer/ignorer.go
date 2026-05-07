package ignorer

import (
	"bufio"
	"os"
	"strings"
)

// Ignorer holds a set of key patterns to exclude from comparison results.
type Ignorer struct {
	keys map[string]struct{}
	prefixes []string
}

// New creates an Ignorer from an explicit list of keys and prefixes.
func New(keys []string, prefixes []string) *Ignorer {
	km := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		km[strings.ToUpper(strings.TrimSpace(k))] = struct{}{}
	}
	norm := make([]string, 0, len(prefixes))
	for _, p := range prefixes {
		norm = append(norm, strings.ToUpper(strings.TrimSpace(p)))
	}
	return &Ignorer{keys: km, prefixes: norm}
}

// LoadFile reads an ignore file where each non-blank, non-comment line is a
// key name or prefix pattern ending with '*'.
func LoadFile(path string) (*Ignorer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var keys, prefixes []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasSuffix(line, "*") {
			prefixes = append(prefixes, strings.TrimSuffix(line, "*"))
		} else {
			keys = append(keys, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return New(keys, prefixes), nil
}

// ShouldIgnore reports whether the given key should be excluded.
func (ig *Ignorer) ShouldIgnore(key string) bool {
	upper := strings.ToUpper(key)
	if _, ok := ig.keys[upper]; ok {
		return true
	}
	for _, p := range ig.prefixes {
		if strings.HasPrefix(upper, p) {
			return true
		}
	}
	return false
}
