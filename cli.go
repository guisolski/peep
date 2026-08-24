package main

import (
	"fmt"
	"io"

	"github.com/guisolski/peep/model"
)

// headlessMode resolves which single headless mode (if any) --query,
// --schema, or --llm requested, validating that at most one was set.
func headlessMode(query string, schema, llm bool) (mode, expr string, err error) {
	count := 0
	if query != "" {
		count++
		mode, expr = "query", query
	}
	if schema {
		count++
		mode = "schema"
	}
	if llm {
		count++
		mode = "llm"
	}
	if count > 1 {
		return "", "", fmt.Errorf("--query, --schema, and --llm are mutually exclusive")
	}
	return mode, expr, nil
}

// runHeadless runs the resolved headless mode against data and writes its
// output to stdout/stderr, returning the process exit code. Composable
// with the TUI's own filter: "peep --query '.items[]' file.json | peep
// --schema" already gives you "schema of a filtered subset" for free.
func runHeadless(data []byte, mode, expr string, stdout, stderr io.Writer) int {
	switch mode {
	case "query":
		return runQuery(data, expr, stdout, stderr)
	case "schema":
		return printExport(data, (*model.Node).Schema, stdout, stderr)
	case "llm":
		return printExport(data, (*model.Node).CompactYAML, stdout, stderr)
	}
	return 0
}

func printExport(data []byte, render func(*model.Node) string, stdout, stderr io.Writer) int {
	root, err := model.ParseJSON(data)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "peep: %v\n", err)
		return 1
	}
	_, _ = fmt.Fprintln(stdout, render(root))
	return 0
}

func runQuery(data []byte, expr string, stdout, stderr io.Writer) int {
	results, err := model.EvalAllJQ(data, expr)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "peep: %v\n", err)
		return 1
	}
	exit := 0
	for _, r := range results {
		if r.Err != nil {
			_, _ = fmt.Fprintf(stderr, "peep: %v\n", r.Err)
			exit = 1
			continue
		}
		_, _ = fmt.Fprintln(stdout, string(r.JSON))
	}
	return exit
}
