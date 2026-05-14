package masker_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/masker"
)

func TestMask_FullStyle(t *testing.T) {
	out := masker.Mask("secret123", masker.Options{Style: masker.StyleFull})
	if out != "*********" {
		t.Errorf("expected all stars, got %q", out)
	}
}

func TestMask_FullStyle_CustomChar(t *testing.T) {
	out := masker.Mask("abc", masker.Options{Style: masker.StyleFull, MaskChar: '#'})
	if out != "###" {
		t.Errorf("got %q", out)
	}
}

func TestMask_PartialStyle_Normal(t *testing.T) {
	// "password" => "pa****rd"
	out := masker.Mask("password", masker.Options{Style: masker.StylePartial, RevealChars: 2})
	if out != "pa****rd" {
		t.Errorf("got %q", out)
	}
}

func TestMask_PartialStyle_TooShort(t *testing.T) {
	// 2*2 >= 3, so fully masked
	out := masker.Mask("abc", masker.Options{Style: masker.StylePartial, RevealChars: 2})
	if out != "***" {
		t.Errorf("got %q", out)
	}
}

func TestMask_PrefixStyle_Normal(t *testing.T) {
	out := masker.Mask("AKIAIOSFODNN7EXAMPLE", masker.Options{Style: masker.StylePrefix, PrefixLen: 4})
	expected := "AKIA****************"
	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}

func TestMask_PrefixStyle_PrefixLongerThanValue(t *testing.T) {
	out := masker.Mask("hi", masker.Options{Style: masker.StylePrefix, PrefixLen: 10})
	if out != "**" {
		t.Errorf("got %q", out)
	}
}

func TestMask_EmptyValue(t *testing.T) {
	out := masker.Mask("", masker.Options{Style: masker.StyleFull})
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestMask_DefaultMaskChar(t *testing.T) {
	out := masker.Mask("xy", masker.Options{Style: masker.StyleFull})
	if out != "**" {
		t.Errorf("got %q", out)
	}
}

func TestMask_PartialStyle_DefaultRevealChars(t *testing.T) {
	// RevealChars defaults to 2; "hello" => "he*lo"
	out := masker.Mask("hello", masker.Options{Style: masker.StylePartial})
	if out != "he*lo" {
		t.Errorf("got %q", out)
	}
}
