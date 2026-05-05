// Package loader provides utilities for loading multiple .env files
// from a directory or a list of file paths, returning a map of
// environment name to parsed key-value pairs.
package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// EnvMap maps an environment name to its key-value pairs.
type EnvMap map[string]map[string]string

// LoadFiles loads the given list of .env file paths and returns an EnvMap.
// The environment name is derived from the file name without its extension.
// For example, ".env.production" becomes "production", and ".env" becomes "default".
func LoadFiles(paths []string) (EnvMap, error) {
	result := make(EnvMap)
	for _, p := range paths {
		name := envName(p)
		kvs, err := parser.ParseFile(p)
		if err != nil {
			return nil, fmt.Errorf("loader: failed to parse %q: %w", p, err)
		}
		result[name] = kvs
	}
	return result, nil
}

// LoadDir scans a directory for files matching the pattern ".env*" and loads
// them all, returning an EnvMap keyed by derived environment name.
func LoadDir(dir string) (EnvMap, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("loader: cannot read directory %q: %w", dir, err)
	}

	var paths []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, ".env") {
			paths = append(paths, filepath.Join(dir, name))
		}
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("loader: no .env files found in %q", dir)
	}

	return LoadFiles(paths)
}

// envName derives a short environment name from a file path.
// ".env"            -> "default"
// ".env.staging"    -> "staging"
// "path/to/.env.qa" -> "qa"
func envName(path string) string {
	base := filepath.Base(path)
	if base == ".env" {
		return "default"
	}
	// strip leading ".env."
	if strings.HasPrefix(base, ".env.") {
		return strings.TrimPrefix(base, ".env.")
	}
	// fallback: use the base name as-is
	return base
}
