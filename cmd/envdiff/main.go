package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/filter"
	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/reporter"
)

func main() {
	files := flag.String("files", "", "Comma-separated list of .env files to compare")
	dir := flag.String("dir", "", "Directory containing .env files to compare")
	format := flag.String("format", "text", "Output format: text, json, markdown, csv")
	onlyMissing := flag.Bool("only-missing", false, "Show only keys missing in at least one environment")
	onlyMismatched := flag.Bool("only-mismatched", false, "Show only keys with differing values")
	keyPrefix := flag.String("key-prefix", "", "Filter keys by prefix (case-insensitive)")
	flag.Parse()

	var envs map[string]map[string]string
	var err error

	switch {
	case *files != "":
		envs, err = loadCommaSeparated(*files)
	case *dir != "":
		envs, err = loader.LoadDir(*dir)
	default:
		fmt.Fprintln(os.Stderr, "error: provide --files or --dir")
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading files: %v\n", err)
		os.Exit(1)
	}

	result := comparator.Compare(envs)

	result.Entries = filter.Apply(result.Entries, filter.Options{
		OnlyMissing:    *onlyMissing,
		OnlyMismatched: *onlyMismatched,
		KeyPrefix:      *keyPrefix,
	})

	if err := writeReport(os.Stdout, *format, result); err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		os.Exit(1)
	}
}

func writeReport(w io.Writer, format string, result comparator.Result) error {
	switch format {
	case "json":
		return reporter.NewJSONReporter(w).Report(result)
	case "markdown":
		return reporter.NewMarkdownReporter(w).Report(result)
	case "csv":
		return reporter.NewCSVReporter(w).Report(result)
	default:
		return reporter.NewTextReporter(w).Report(result)
	}
}

func loadCommaSeparated(s string) (map[string]map[string]string, error) {
	paths := splitComma(s)
	return loader.LoadFiles(paths)
}

func splitComma(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
