package differ

import (
	"strings"
	"testing"
)

func TestDiffValues_Identical(t *testing.T) {
	d := DiffValues("KEY", "hello", "hello")
	if len(d.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(d.Lines))
	}
	if d.Lines[0].Op != ' ' {
		t.Errorf("expected context op, got %c", d.Lines[0].Op)
	}
}

func TestDiffValues_SingleLineChange(t *testing.T) {
	d := DiffValues("KEY", "foo", "bar")
	ops := opsString(d.Lines)
	if ops != "-+" {
		t.Errorf("expected '-+', got %q", ops)
	}
	if d.Lines[0].Text != "foo" {
		t.Errorf("expected removed text 'foo', got %q", d.Lines[0].Text)
	}
	if d.Lines[1].Text != "bar" {
		t.Errorf("expected added text 'bar', got %q", d.Lines[1].Text)
	}
}

func TestDiffValues_MultiLinePartialMatch(t *testing.T) {
	left := "line1\nline2\nline3"
	right := "line1\nchanged\nline3"
	d := DiffValues("MULTI", left, right)
	ops := opsString(d.Lines)
	// line1 context, line2 removed, changed added, line3 context
	if ops != " -+ " {
		t.Errorf("unexpected ops %q", ops)
	}
}

func TestDiffValues_EmptyLeft(t *testing.T) {
	d := DiffValues("KEY", "", "value")
	ops := opsString(d.Lines)
	if !strings.Contains(ops, "+") {
		t.Errorf("expected addition op, got %q", ops)
	}
}

func TestDiffValues_EmptyRight(t *testing.T) {
	d := DiffValues("KEY", "value", "")
	ops := opsString(d.Lines)
	if !strings.Contains(ops, "-") {
		t.Errorf("expected removal op, got %q", ops)
	}
}

func TestValueDiff_Format(t *testing.T) {
	d := DiffValues("MY_KEY", "old", "new")
	out := d.Format()
	if !strings.Contains(out, "key: MY_KEY") {
		t.Errorf("format missing key header: %q", out)
	}
	if !strings.Contains(out, "- old") {
		t.Errorf("format missing removal line: %q", out)
	}
	if !strings.Contains(out, "+ new") {
		t.Errorf("format missing addition line: %q", out)
	}
}

func TestDiffValues_KeyPreserved(t *testing.T) {
	d := DiffValues("DB_URL", "a", "b")
	if d.Key != "DB_URL" {
		t.Errorf("expected key DB_URL, got %q", d.Key)
	}
	if d.Left != "a" || d.Right != "b" {
		t.Errorf("left/right not preserved")
	}
}

func opsString(lines []LineDiff) string {
	var sb strings.Builder
	for _, l := range lines {
		sb.WriteRune(l.Op)
	}
	return sb.String()
}
