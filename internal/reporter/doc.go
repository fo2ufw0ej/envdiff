// Package reporter provides output formatters for envdiff comparison results.
//
// Two reporter implementations are available:
//
//   - TextReporter: writes a human-readable, coloured diff-style summary to
//     any io.Writer.  Use NewTextReporter to create one.
//
//   - JSONReporter: serialises the full comparison result as an indented JSON
//     document, suitable for machine consumption or downstream tooling.  Use
//     NewJSONReporter to create one.
//
// Both reporters accept a comparator.Result value produced by
// comparator.Compare and write their output to the io.Writer supplied at
// construction time.  They return a non-nil error only when the underlying
// write fails.
package reporter
