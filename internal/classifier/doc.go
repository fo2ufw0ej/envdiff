// Package classifier assigns semantic categories to environment variable keys
// based on configurable keyword rules.
//
// Built-in categories include database, auth, network, observability, and feature.
// Custom rules can be supplied via NewWithRules.
//
// Example:
//
//	c := classifier.New()
//	cat := c.Classify("DATABASE_URL") // => "database"
//
//	// Annotate comparator results with a "category" tag:
//	annotated := c.ClassifyResults(results)
package classifier
