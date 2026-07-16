package model

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"peep/ui"
)

// GraphModel implements a Miller-columns layout:
//
//	left panel  │ middle panel (focused node) │ right panel (preview)
//	parent cols │ focused children            │ selected child's children
//
// h/← goes up to parent, l/→ or Enter drills into selected child.
type GraphModel struct {
	root    *Node
	focused *Node // node whose children fill the middle panel
	cursor  int   // selected child index within focused.Children
	width   int
	height  int
}

func NewGraphModel(root *Node, width, height int) *GraphModel {
	return &GraphModel{
		root:    root,
		focused: root,
		width:   width,
		height:  height - 1,
	}
}

func (m *GraphModel) Init() tea.Cmd { return nil }

func (m *GraphModel) Update(msg tea.Msg) (SubModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 1
	case tea.KeyMsg:
		children := m.focusedChildren()
		switch msg.String() {
		case "j", "down":
			if m.cursor < len(children)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "l", "right", "enter":
			if m.cursor < len(children) {
				child := children[m.cursor]
				if child.Type == TypeObject || child.Type == TypeArray {
					m.focused = child
					m.cursor = 0
				}
			}
		case "h", "left":
			if m.focused.Parent != nil {
				prev := m.focused
				m.focused = m.focused.Parent
				// restore cursor to the child we came from
				for i, c := range m.focusedChildren() {
					if c == prev {
						m.cursor = i
						break
					}
				}
			}
		}
	}
	return m, nil
}

// focusedChildren returns the children to display in the middle panel.
func (m *GraphModel) focusedChildren() []*Node {
	if m.focused == nil {
		return nil
	}
	return m.focused.Children
}

func (m *GraphModel) View() string {
	leftW := m.width / 5
	rightW := m.width / 4
	midW := m.width - leftW - rightW

	left := m.renderPanel(m.focused.Parent, leftW, -1, m.focused)
	mid := m.renderPanel(m.focused, midW, m.cursor, nil)

	var previewNode *Node
	children := m.focusedChildren()
	if m.cursor < len(children) {
		previewNode = children[m.cursor]
	}
	right := m.renderPanel(previewNode, rightW, -1, nil)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, mid, right)
}

// renderPanel renders a single column panel. n is the parent node whose
// children are listed. cursor selects one row (-1 = no selection).
// highlight marks one node as the "active" entry (used in left panel).
func (m *GraphModel) renderPanel(n *Node, w, cursor int, highlight *Node) string {
	h := m.height
	style := lipgloss.NewStyle().
		Width(w-2).Height(h).
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(lipgloss.AdaptiveColor{Light: "8", Dark: "8"})

	if n == nil {
		return style.Render("")
	}

	var sb strings.Builder

	// Header: show the node's own key/summary.
	title := n.Key
	if title == "" {
		title = "(root)"
	}
	header := ui.Blue.Render(title) + " " + ui.Muted.Render(n.Summary())
	sb.WriteString(header)
	sb.WriteByte('\n')

	// List children.
	for i, child := range n.Children {
		line := renderGraphRow(child)
		selected := i == cursor
		isHighlighted := child == highlight

		switch {
		case selected:
			plain := ansi.Strip(line)
			padded := plain + strings.Repeat(" ", max(0, w-2-len([]rune(plain))))
			line = ui.Selected.Render(padded)
		case isHighlighted:
			line = ui.Blue.Render(ansi.Strip(line))
		}
		sb.WriteString(line)
		sb.WriteByte('\n')
	}

	return style.Render(sb.String())
}

func renderGraphRow(n *Node) string {
	key := ""
	if n.Key != "" {
		key = ui.Blue.Render(n.Key) + ": "
	}
	switch n.Type {
	case TypeObject, TypeArray:
		return key + ui.Muted.Render(n.Summary()) + ui.Muted.Render(" ▶")
	case TypeString:
		return key + ui.Green.Render(fmt.Sprintf("%q", n.StrVal))
	case TypeNumber:
		return key + ui.Peach.Render(n.Summary())
	case TypeBool:
		return key + ui.Red.Render(n.Summary())
	case TypeNull:
		return key + ui.Muted.Render("null")
	}
	return key
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

