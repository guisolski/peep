package input

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/atotto/clipboard"
)

type Source string

const (
	SourceFile        Source = "file"
	SourceStdin       Source = "stdin"
	SourceClipboard   Source = "clipboard"
	SourceInteractive Source = "paste"
)

type Result struct {
	Data   []byte
	Source Source
	Name   string
}

// Load detects the JSON source in priority order:
// explicit file arg → piped stdin → --paste flag → interactive paste.
func Load(args []string, pasteFlag bool) (*Result, error) {
	if len(args) > 0 {
		data, err := os.ReadFile(args[0])
		if err != nil {
			return nil, fmt.Errorf("read file: %w", err)
		}
		return &Result{Data: data, Source: SourceFile, Name: args[0]}, nil
	}

	if !isTerminal(os.Stdin) {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("read stdin: %w", err)
		}
		return &Result{Data: data, Source: SourceStdin, Name: "stdin"}, nil
	}

	if pasteFlag {
		text, err := clipboard.ReadAll()
		if err != nil {
			return nil, fmt.Errorf("read clipboard: %w", err)
		}
		return &Result{Data: []byte(text), Source: SourceClipboard, Name: "clipboard"}, nil
	}

	fmt.Fprintln(os.Stderr, "Paste JSON and press Ctrl+D:")
	data, err := io.ReadAll(bufio.NewReader(os.Stdin))
	if err != nil {
		return nil, fmt.Errorf("read interactive: %w", err)
	}
	return &Result{Data: data, Source: SourceInteractive, Name: "paste"}, nil
}

func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
