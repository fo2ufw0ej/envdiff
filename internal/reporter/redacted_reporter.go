package reporter

import (
	"io"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/redactor"
)

// RedactedReporter wraps any Reporter and masks sensitive values in the
// comparison results before delegating to the inner reporter.
type RedactedReporter struct {
	inner   Reporter
	redactor *redactor.Redactor
}

// NewRedactedReporter returns a Reporter that redacts sensitive values from
// results before passing them to inner.
func NewRedactedReporter(inner Reporter, r *redactor.Redactor) *RedactedReporter {
	if r == nil {
		r = redactor.New(nil, "")
	}
	return &RedactedReporter{inner: inner, redactor: r}
}

// Report redacts sensitive values from each result entry and then delegates
// to the wrapped reporter.
func (rr *RedactedReporter) Report(w io.Writer, results []comparator.Result) error {
	redacted := make([]comparator.Result, len(results))
	for i, res := range results {
		redacted[i] = redactResult(res, rr.redactor)
	}
	return rr.inner.Report(w, redacted)
}

func redactResult(res comparator.Result, r *redactor.Redactor) comparator.Result {
	if !r.IsSensitive(res.Key) {
		return res
	}
	masked := make(map[string]string, len(res.Values))
	for env, val := range res.Values {
		if val != "" {
			masked[env] = r.RedactValue(res.Key, val)
		} else {
			masked[env] = val
		}
	}
	return comparator.Result{
		Key:    res.Key,
		Status: res.Status,
		Values: masked,
	}
}
