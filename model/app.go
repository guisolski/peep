package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	atclip "github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/guisolski/peep/ui"
)

type Mode int

const (
	ModeTree Mode = iota
	ModeGraph
	ModeRaw
	ModeSearch
	ModeFilter
)

type pendingKeyTimeoutMsg struct{ key string }

type App struct {
	mode       Mode
	tree       *TreeModel
	graph      *GraphModel
	raw        *RawModel
	search     *SearchModel
	filter     *FilterModel
	source     string
	width      int
	height     int
	data       []byte
	pendingKey string
	statusMsg  string
	embedded   bool
}

func NewApp(data []byte, source string, width, height int) (*App, error) {
	root, err := ParseJSON(data)
	if err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}
	tree := NewTreeModel(root, width, height)
	return &App{
		mode:   ModeTree,
		tree:   tree,
		graph:  NewGraphModel(root, width, height),
		raw:    NewRawModel(data, width, height),
		search: NewSearchModel(tree),
		filter: NewFilterModel(data, width, height),
		source: source,
		width:  width,
		height: height,
		data:   data,
	}, nil
}

// SetEmbedded marks the app as hosted inside another Bubble Tea program:
// quit keys are ignored and the quit hint is hidden from the status bar.
func (a *App) SetEmbedded(v bool) {
	a.embedded = v
}

// InputActive reports whether a text input (search or filter) is focused,
// so hosts can avoid stealing printable keys while the user types.
func (a *App) InputActive() bool {
	return a.mode == ModeSearch || a.mode == ModeFilter
}

func (a *App) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width, a.height = msg.Width, msg.Height
		a.tree.Update(msg)
		a.graph.Update(msg)
		a.raw.Update(msg)
		a.search.Update(msg)
		a.filter.Update(msg)
		return a, nil

	case pendingKeyTimeoutMsg:
		if a.pendingKey == msg.key {
			a.pendingKey = ""
			switch msg.key {
			case "g":
				a.toggleGraph()
			case "y":
				if n := a.tree.CurrentNode(); n != nil {
					if err := atclip.WriteAll(n.Summary()); err == nil {
						a.statusMsg = "copied value"
					}
				}
			}
		}
		return a, nil

	case tea.KeyMsg:
		// Global quit (standalone only): ctrl+c always, q outside text inputs.
		if !a.embedded {
			if msg.String() == "ctrl+c" {
				return a, tea.Quit
			}
			if msg.String() == "q" && !a.InputActive() {
				return a, tea.Quit
			}
		}

		// Esc exits search/filter back to tree.
		if msg.String() == "esc" && (a.mode == ModeSearch || a.mode == ModeFilter) {
			a.tree.ClearHighlights()
			a.mode = ModeTree
			a.pendingKey = ""
			return a, nil
		}

		// Multi-key sequences and mode switches.
		switch msg.String() {
		case "g":
			if a.pendingKey == "g" {
				a.pendingKey = ""
				a.tree.JumpToTop()
				return a, nil
			}
			a.pendingKey = "g"
			return a, tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
				return pendingKeyTimeoutMsg{"g"}
			})

		case "y":
			if a.mode != ModeTree && a.mode != ModeSearch {
				break
			}
			if a.pendingKey == "y" {
				a.pendingKey = ""
				if n := a.tree.CurrentNode(); n != nil {
					if b, err := json.MarshalIndent(n.ToInterface(), "", "  "); err == nil {
						if err2 := atclip.WriteAll(string(b)); err2 == nil {
							a.statusMsg = "copied subtree"
						}
					}
				}
				return a, nil
			}
			a.pendingKey = "y"
			return a, tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
				return pendingKeyTimeoutMsg{"y"}
			})

		case "Y":
			if a.mode == ModeTree || a.mode == ModeSearch {
				a.pendingKey = ""
				if n := a.tree.CurrentNode(); n != nil {
					if err := atclip.WriteAll(n.Path()); err == nil {
						a.statusMsg = "copied path"
					}
				}
				return a, nil
			}

		case "r":
			a.pendingKey = ""
			if a.mode == ModeRaw {
				a.mode = ModeTree
			} else {
				a.mode = ModeRaw
			}
			return a, nil

		case "/":
			a.pendingKey = ""
			a.mode = ModeSearch
			return a, a.search.Init()

		case ":":
			a.pendingKey = ""
			a.mode = ModeFilter
			return a, a.filter.Init()
		}

		// Route remaining keys to the active sub-model.
		a.pendingKey = ""
		switch a.mode {
		case ModeTree:
			sub, cmd := a.tree.Update(msg)
			a.tree = sub.(*TreeModel)
			return a, cmd
		case ModeGraph:
			sub, cmd := a.graph.Update(msg)
			a.graph = sub.(*GraphModel)
			return a, cmd
		case ModeRaw:
			sub, cmd := a.raw.Update(msg)
			a.raw = sub.(*RawModel)
			return a, cmd
		case ModeSearch:
			sub, cmd := a.search.Update(msg)
			if sm, ok := sub.(*SearchModel); ok {
				a.search = sm
			} else {
				a.mode = ModeTree
			}
			return a, cmd
		case ModeFilter:
			sub, cmd := a.filter.Update(msg)
			if fm, ok := sub.(*FilterModel); ok {
				a.filter = fm
			} else {
				a.mode = ModeTree
			}
			return a, cmd
		}
	}
	return a, nil
}

func (a *App) toggleGraph() {
	if a.mode == ModeGraph {
		a.mode = ModeTree
	} else {
		a.mode = ModeGraph
	}
}

func (a *App) View() string {
	var content string
	switch a.mode {
	case ModeTree:
		content = a.tree.View()
	case ModeGraph:
		content = a.graph.View()
	case ModeRaw:
		content = a.raw.View()
	case ModeSearch:
		content = a.search.View()
	case ModeFilter:
		content = a.filter.View()
	}
	// Pad short content so the status bar stays pinned to the bottom row.
	lines := strings.Count(content, "\n") + 1
	if pad := a.height - 1 - lines; pad > 0 {
		content += strings.Repeat("\n", pad)
	}
	return content + a.statusBar()
}

func (a *App) statusBar() string {
	source := " " + a.source + " "
	path := "."
	if n := a.tree.CurrentNode(); n != nil {
		path = n.Path()
	}
	modePart := ui.ModeTag.Render(modeLabel(a.mode))
	msg := ""
	if a.statusMsg != "" {
		msg = "  ✓ " + a.statusMsg
		a.statusMsg = ""
	}
	hints := "  j/k ↑↓  l/→ expand  h/← collapse  g graph  r raw  / search  : filter"
	if !a.embedded {
		hints += "  q quit"
	}
	bar := strings.Join([]string{source, " ", path, " ", modePart, msg, hints}, "")
	return "\n" + ui.StatusBar.Width(a.width).Render(bar)
}

func modeLabel(m Mode) string {
	switch m {
	case ModeTree:
		return "TREE"
	case ModeGraph:
		return "GRAPH"
	case ModeRaw:
		return "RAW"
	case ModeSearch:
		return "SEARCH"
	case ModeFilter:
		return "FILTER"
	}
	return "?"
}
