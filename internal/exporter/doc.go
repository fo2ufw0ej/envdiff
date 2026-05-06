// Package exporter provides utilities to export envdiff comparison results
// into actionable output formats.
//
// Supported formats:
//
//   - env:   Produces a shell-sourceable .env file containing keys from a
//            reference environment, suitable for bootstrapping a missing env.
//
//   - patch: Produces a diff-style patch listing keys prefixed with '+' that
//            are present in the reference but absent or different elsewhere.
//
// Example usage:
//
//	ex := exporter.New(exporter.FormatEnv, "production")
//	if err := ex.Export(os.Stdout, results); err != nil {
//		log.Fatal(err)
//	}
package exporter
