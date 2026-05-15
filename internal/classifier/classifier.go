package classifier

import (
	"strings"

	"github.com/yourorg/envdiff/internal/comparator"
)

// Category represents a semantic classification of an env key.
type Category string

const (
	CategoryDatabase  Category = "database"
	CategoryAuth      Category = "auth"
	CategoryNetwork   Category = "network"
	CategoryObserv    Category = "observability"
	CategoryFeature   Category = "feature"
	CategoryUnknown   Category = "unknown"
)

// Rule maps a set of key substrings to a Category.
type Rule struct {
	Keywords []string
	Category Category
}

var defaultRules = []Rule{
	{Keywords: []string{"db", "database", "postgres", "mysql", "mongo", "redis", "dsn"}, Category: CategoryDatabase},
	{Keywords: []string{"secret", "token", "auth", "password", "passwd", "api_key", "jwt"}, Category: CategoryAuth},
	{Keywords: []string{"host", "port", "url", "addr", "endpoint", "proxy"}, Category: CategoryNetwork},
	{Keywords: []string{"log", "trace", "metric", "sentry", "datadog", "otel", "debug"}, Category: CategoryObserv},
	{Keywords: []string{"feature", "flag", "enable", "disable", "toggle"}, Category: CategoryFeature},
}

// Classifier assigns categories to env keys.
type Classifier struct {
	rules []Rule
}

// New returns a Classifier using the default built-in rules.
func New() *Classifier {
	return &Classifier{rules: defaultRules}
}

// NewWithRules returns a Classifier using the provided rules.
func NewWithRules(rules []Rule) *Classifier {
	return &Classifier{rules: rules}
}

// Classify returns the Category for a given key.
func (c *Classifier) Classify(key string) Category {
	lower := strings.ToLower(key)
	for _, rule := range c.rules {
		for _, kw := range rule.Keywords {
			if strings.Contains(lower, kw) {
				return rule.Category
			}
		}
	}
	return CategoryUnknown
}

// ClassifyResults annotates each comparator.Result with a category tag.
// The category is stored in the Tags map under the key "category".
func (c *Classifier) ClassifyResults(results []comparator.Result) []comparator.Result {
	out := make([]comparator.Result, len(results))
	for i, r := range results {
		copy := r
		if copy.Tags == nil {
			copy.Tags = make(map[string]string)
		}
		copy.Tags["category"] = string(c.Classify(r.Key))
		out[i] = copy
	}
	return out
}
