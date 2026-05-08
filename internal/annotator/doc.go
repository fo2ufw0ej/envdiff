// Package annotator attaches human-readable annotations to comparator results.
//
// Each annotation carries a severity level ("info", "warn", or "error") and a
// descriptive message explaining why the key was flagged. Keys that match
// sensitive patterns (e.g. containing "secret", "token", "password") are
// escalated to "error" severity when they are mismatched across environments.
//
// Usage:
//
//	a := annotator.New()
//	annotations := a.Annotate(results)
//	for _, ann := range annotations {
//		fmt.Printf("[%s] %s\n", ann.Severity, ann.Message)
//	}
package annotator
