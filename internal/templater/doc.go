// Package templater generates .env template files from a set of environments.
//
// A template replaces every value with a placeholder (e.g. <DB_HOST>) so that
// the file can be safely committed to version control or shared with new team
// members without exposing real credentials.
//
// Usage:
//
//	tmpl := templater.New()
//	result := tmpl.Render("production", referenceKeys, envMap)
//	for _, line := range result.Lines {
//		fmt.Println(line)
//	}
//
// Keys that are present in the reference set but absent from the supplied
// environment are still rendered with a placeholder and are recorded in
// Result.Missing for downstream reporting.
package templater
