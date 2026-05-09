package templater

import (
	"fmt"
	"sort"
	"strings"
)

// Template represents a .env template where values are replaced with
// placeholder descriptions derived from the key names.
type Template struct {
	Placeholder func(key string) string
}

// Result holds the rendered template output for a single environment.
type Result struct {
	EnvName  string
	Lines    []string
	Missing  []string
}

// New returns a Template with the default placeholder strategy.
// Keys are rendered as <KEY_NAME> placeholders.
func New() *Template {
	return &Template{
		Placeholder: func(key string) string {
			return fmt.Sprintf("<%s>", strings.ToUpper(key))
		},
	}
}

// NewWithPlaceholder returns a Template with a custom placeholder function.
func NewWithPlaceholder(fn func(key string) string) *Template {
	if fn == nil {
		return New()
	}
	return &Template{Placeholder: fn}
}

// Render produces a template from a reference key set and the values of a
// specific environment. Keys present in reference but absent in env are
// recorded in Result.Missing.
func (t *Template) Render(envName string, reference []string, env map[string]string) Result {
	keys := make([]string, len(reference))
	copy(keys, reference)
	sort.Strings(keys)

	result := Result{EnvName: envName}
	for _, key := range keys {
		if _, ok := env[key]; !ok {
			result.Missing = append(result.Missing, key)
			result.Lines = append(result.Lines, fmt.Sprintf("%s=%s", key, t.Placeholder(key)))
			continue
		}
		result.Lines = append(result.Lines, fmt.Sprintf("%s=%s", key, t.Placeholder(key)))
	}
	return result
}

// RenderAll renders a template for every environment in envs, using the union
// of all keys as the reference set.
func (t *Template) RenderAll(envs map[string]map[string]string) []Result {
	keySet := map[string]struct{}{}
	for _, env := range envs {
		for k := range env {
			keySet[k] = struct{}{}
		}
	}
	reference := make([]string, 0, len(keySet))
	for k := range keySet {
		reference = append(reference, k)
	}

	names := make([]string, 0, len(envs))
	for name := range envs {
		names = append(names, name)
	}
	sort.Strings(names)

	results := make([]Result, 0, len(names))
	for _, name := range names {
		results = append(results, t.Render(name, reference, envs[name]))
	}
	return results
}
