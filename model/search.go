package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/guisolski/peep/ui"
)

type SearchModel struct {
	tree    *TreeModel
	input   textinput.Model
	matches []int
	current int
	allFlat []*Node
}

func NewSearchModel(tree *TreeModel) *SearchModel {
	ti := textinput.New()
	ti.Placeholder = "search…"
	ti.Focus()
	sm := &SearchModel{tree: tree, input: ti}
	sm.allFlat = tree.FlatAll()
	return sm
}

func (m *SearchModel) SetQuery(q string) {
	m.input.SetValue(q)
	m.recompute()
}

func (m *SearchModel) Query() string  { return m.input.Value() }
func (m *SearchModel) Matches() []int { return m.matches }
func (m *SearchModel) Current() int   { return m.current }

func (m *SearchModel) Next() {
	if len(m.matches) == 0 {
		return
	}
	m.current = (m.current + 1) % len(m.matches)
	m.applyHighlights()
}

func (m *SearchModel) Prev() {
	if len(m.matches) == 0 {
		return
	}
	m.current = (m.current - 1 + len(m.matches)) % len(m.matches)
	m.applyHighlights()
}

func (m *SearchModel) recompute() {
	q := strings.ToLower(m.input.Value())
	m.matches = nil
	for i, n := range m.allFlat {
		if nodeMatchesQuery(n, q) {
			m.matches = append(m.matches, i)
		}
	}
	m.current = 0
	m.applyHighlights()
}

func nodeMatchesQuery(n *Node, q string) bool {
	if q == "" {
		return false
	}
	if strings.Contains(strings.ToLower(n.Key), q) {
		return true
	}
	return strings.Contains(strings.ToLower(n.Summary()), q)
}

func (m *SearchModel) applyHighlights() {
	h := map[int]bool{}
	for _, idx := range m.matches {
		h[idx] = true
	}
	m.tree.SetHighlights(h)
	if len(m.matches) > 0 {
		m.tree.cursor = m.matches[m.current]
		m.tree.clampOffset()
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
	if len(m.matches) > 0 {
		sb.WriteString("  " + ui.Muted.Render(fmt.Sprintf("[%d/%d]", m.current+1, len(m.matches))))
	}
	return sb.String()
}
