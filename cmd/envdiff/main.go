// Command envdiff compares .env files across environments and highlights
// missing or mismatched keys.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/reporter"
)

func main() {
	var (
		dir     = flag.String("dir", "", "directory containing .env files to compare")
		files   = flag.String("files", "", "comma-separated list of .env files to compare")
		outFmt  = flag.String("format", "text", "output format: text (default)")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  envdiff -dir ./envs\n")
		fmt.Fprintf(os.Stderr, "  envdiff -files .env.development,.env.production\n")
	}
	flag.Parse()

	if *dir == "" && *files == "" {
		fmt.Fprintln(os.Stderr, "error: provide -dir or -files")
		flag.Usage()
		os.Exit(1)
	}

	var envs map[string]map[string]string
	var err error

	if *dir != "" {
		envs, err = loader.LoadDir(*dir)
	} else {
		envs, err = loadCommaSeparated(*files)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading files: %v\n", err)
		os.Exit(1)
	}

	if len(envs) == 0 {
		fmt.Fprintln(os.Stderr, "no .env files found")
		os.Exit(1)
	}

	results := comparator.Compare(envs)

	var rep reporter.Reporter
	switch *outFmt {
	default:
		rep = reporter.NewTextReporter(os.Stdout)
	}

	if err := rep.Report(results); err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		os.Exit(1)
	}

	if len(results) > 0 {
		os.Exit(2)
	}
}

// loadCommaSeparated parses a comma-separated list of file paths and delegates
// to loader.LoadFiles.
func loadCommaSeparated(raw string) (map[string]map[string]string, error) {
	var paths []string
	for _, p := range splitComma(raw) {
		if p != "" {
			paths = append(paths, p)
		}
	}
	return loader.LoadFiles(paths)
}

func splitComma(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	out = append(out, s[start:])
	return out
}
