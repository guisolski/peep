package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/guisolski/peep/input"
	"github.com/guisolski/peep/model"
)

// version is overridden at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	paste := flag.Bool("paste", false, "read JSON from system clipboard")
	versionFlag := flag.Bool("version", false, "print version and exit")
	query := flag.String("query", "", "run a jq expression and print each result, without opening the TUI")
	schema := flag.Bool("schema", false, "print a compact schema (types + example values) and exit, without opening the TUI")
	llm := flag.Bool("llm", false, "print a compact, LLM-friendly serialization and exit, without opening the TUI")
	flag.Parse()

	if *versionFlag {
		fmt.Println("peep " + version)
		return
	}

	// Resolved before loading input: a flag-usage mistake should fail fast
	// rather than risk blocking on the interactive-paste stdin prompt.
	mode, expr, err := headlessMode(*query, *schema, *llm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "peep: %v\n", err)
		os.Exit(2)
	}

	result, err := input.Load(flag.Args(), *paste)
	if err != nil {
		fmt.Fprintf(os.Stderr, "peep: %v\n", err)
		os.Exit(1)
	}

	if mode != "" {
		os.Exit(runHeadless(result.Data, mode, expr, os.Stdout, os.Stderr))
	}

	app, err := model.NewApp(result.Data, result.Name, 80, 24)
	if err != nil {
		fmt.Fprintf(os.Stderr, "peep: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "peep: %v\n", err)
		os.Exit(1)
	}
}
