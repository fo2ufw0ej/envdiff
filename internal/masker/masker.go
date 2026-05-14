// Package masker provides utilities for partially masking env values
// so that sensitive data can be displayed in reports without full exposure.
package masker

import "strings"

// Style controls how a value is masked.
type Style int

const (
	// StyleFull replaces the entire value with asterisks.
	StyleFull Style = iota
	// StylePartial reveals the first and last characters, masking the middle.
	StylePartial
	// StylePrefix reveals only the first N characters.
	StylePrefix
)

// Options configures masking behaviour.
type Options struct {
	Style       Style
	PrefixLen   int // used by StylePrefix; defaults to 3
	RevealChars int // used by StylePartial; defaults to 2
	MaskChar    rune
}

// defaultOptions returns sensible defaults.
func defaultOptions(o Options) Options {
	if o.MaskChar == 0 {
		o.MaskChar = '*'
	}
	if o.PrefixLen <= 0 {
		o.PrefixLen = 3
	}
	if o.RevealChars <= 0 {
		o.RevealChars = 2
	}
	return o
}

// Mask applies the given Options to value and returns the masked string.
// Empty values are returned unchanged.
func Mask(value string, o Options) string {
	if value == "" {
		return value
	}
	o = defaultOptions(o)
	runes := []rune(value)
	n := len(runes)
	mask := strings.Repeat(string(o.MaskChar), n)

	switch o.Style {
	case StylePartial:
		reveal := o.RevealChars
		if 2*reveal >= n {
			// Value too short to partially reveal — mask fully.
			return mask
		}
		midLen := n - 2*reveal
		mid := strings.Repeat(string(o.MaskChar), midLen)
		return string(runes[:reveal]) + mid + string(runes[n-reveal:])
	case StylePrefix:
		pfx := o.PrefixLen
		if pfx >= n {
			return mask
		}
		return string(runes[:pfx]) + strings.Repeat(string(o.MaskChar), n-pfx)
	default: // StyleFull
		return mask
	}
}
