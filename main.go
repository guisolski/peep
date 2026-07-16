package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/guisolski/peep/input"
	"github.com/guisolski/peep/model"
)

func main() {
	paste := flag.Bool("paste", false, "read JSON from system clipboard")
	flag.Parse()

	result, err := input.Load(flag.Args(), *paste)
	if err != nil {
		fmt.Fprintf(os.Stderr, "peep: %v\n", err)
		os.Exit(1)
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
