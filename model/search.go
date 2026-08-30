package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/guisolski/peep/ui"
)

type SearchModel struct {
	tree       *TreeModel
	input      textinput.Model
	matchNodes []*Node
	current    int
	allFlat    []*Node
}

func NewSearchModel(tree *TreeModel) *SearchModel {
	ti := textinput.New()
	ti.Placeholder = "ask in English · id:2 · greater than 5"
	ti.Focus()
	sm := &SearchModel{tree: tree, input: ti}
	sm.allFlat = tree.FlatAll()
	return sm
}

func (m *SearchModel) SetQuery(q string) {
	m.input.SetValue(q)
	m.recompute()
}

func (m *SearchModel) Query() string { return m.input.Value() }

// Matches returns indices into the FlatAll snapshot for the current hits.
func (m *SearchModel) Matches() []int {
	out := make([]int, 0, len(m.matchNodes))
	for _, n := range m.matchNodes {
		if i := indexInNodes(m.allFlat, n); i >= 0 {
			out = append(out, i)
		}
	}
	return out
}

func (m *SearchModel) Current() int { return m.current }

func (m *SearchModel) Next() {
	if len(m.matchNodes) == 0 {
		return
	}
	m.current = (m.current + 1) % len(m.matchNodes)
	m.applyHighlights()
}

func (m *SearchModel) Prev() {
	if len(m.matchNodes) == 0 {
		return
	}
	m.current = (m.current - 1 + len(m.matchNodes)) % len(m.matchNodes)
	m.applyHighlights()
}

func (m *SearchModel) recompute() {
	clauses := parseSearchQuery(m.input.Value())
	m.matchNodes = nil
	for _, n := range m.allFlat {
		if nodeMatchesClauses(n, clauses) {
			m.matchNodes = append(m.matchNodes, n)
		}
	}
	m.current = 0
	m.applyHighlights()
}

func (m *SearchModel) applyHighlights() {
	for _, n := range m.matchNodes {
		expandAncestors(n)
	}
	m.tree.rebuild()

	h := map[int]bool{}
	for _, n := range m.matchNodes {
		if i := visibleIndex(m.tree, n); i >= 0 {
			h[i] = true
		}
	}
	m.tree.SetHighlights(h)
	if len(m.matchNodes) > 0 {
		if i := visibleIndex(m.tree, m.matchNodes[m.current]); i >= 0 {
			m.tree.cursor = i
			m.tree.clampOffset()
		}
	}
}

func (m *SearchModel) Init() tea.Cmd { return textinput.Blink }

func (m *SearchModel) Update(msg tea.Msg) (SubModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			m.Next()
			return m, nil
		case "N":
			m.Prev()
			return m, nil
		case "esc":
			m.tree.ClearHighlights()
			return m, nil
		default:
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			m.recompute()
			return m, cmd
		}
	}
	return m, nil
}

func (m *SearchModel) View() string {
	var sb strings.Builder
	sb.WriteString(m.tree.View())
	sb.WriteString(ui.Yellow.Render("/") + " " + m.input.View())
	if len(m.matchNodes) > 0 {
		sb.WriteString("  " + ui.Muted.Render(fmt.Sprintf("[%d/%d]", m.current+1, len(m.matchNodes))))
	}
	return sb.String()
}
