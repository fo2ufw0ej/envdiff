package patcher_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/patcher"
)

func actionMap(results []patcher.Result) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		m[r.Key] = r.Action
	}
	return m
}

func TestPatch_AddMissing(t *testing.T) {
	dst := map[string]string{"A": "1", "B": "2"}
	src := map[string]string{"B": "99", "C": "3"}

	out, results := patcher.Patch(dst, src, patcher.StrategyAddMissing)

	if out["A"] != "1" {
		t.Errorf("A should be unchanged, got %q", out["A"])
	}
	if out["B"] != "2" {
		t.Errorf("B should not be overwritten, got %q", out["B"])
	}
	if out["C"] != "3" {
		t.Errorf("C should be added, got %q", out["C"])
	}

	am := actionMap(results)
	if am["B"] != "skipped" {
		t.Errorf("B action should be skipped, got %q", am["B"])
	}
	if am["C"] != "added" {
		t.Errorf("C action should be added, got %q", am["C"])
	}
}

func TestPatch_Overwrite(t *testing.T) {
	dst := map[string]string{"A": "1", "B": "2"}
	src := map[string]string{"B": "99", "C": "3"}

	out, results := patcher.Patch(dst, src, patcher.StrategyOverwrite)

	if out["B"] != "99" {
		t.Errorf("B should be overwritten, got %q", out["B"])
	}
	am := actionMap(results)
	if am["B"] != "overwritten" {
		t.Errorf("B action should be overwritten, got %q", am["B"])
	}
}

func TestPatch_DoesNotMutateDst(t *testing.T) {
	dst := map[string]string{"X": "original"}
	src := map[string]string{"X": "changed"}

	patcher.Patch(dst, src, patcher.StrategyOverwrite)

	if dst["X"] != "original" {
		t.Error("Patch must not mutate the dst map")
	}
}

func TestWriteFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.env")

	data := map[string]string{
		"FOO": "bar",
		"HAS SPACE": "val",
		"Z": "last",
	}
	if err := patcher.WriteFile(path, data); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	content := string(b)
	if !strings.Contains(content, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got:\n%s", content)
	}
	if !strings.Contains(content, "Z=last") {
		t.Errorf("expected Z=last in output, got:\n%s", content)
	}
}
