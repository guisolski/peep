package model

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"peep/ui"
)

type TreeModel struct {
	root       *Node
	flat       []*Node
	cursor     int
	offset     int
	height     int
	width      int
	highlights map[int]bool
}

func NewTreeModel(root *Node, width, height int) *TreeModel {
	m := &TreeModel{
		root:       root,
		width:      width,
		height:     height,
		highlights: map[int]bool{},
	}
	collapseAll(root)
	root.Collapsed = false // keep root expanded
	m.rebuild()
	return m
}

func collapseAll(n *Node) {
	if n.Type == TypeObject || n.Type == TypeArray {
		n.Collapsed = true
		for _, c := range n.Children {
			collapseAll(c)
		}
	}
}

func (m *TreeModel) rebuild() {
	m.flat = nil
	flattenVisible(m.root, &m.flat)
}

func flattenVisible(n *Node, out *[]*Node) {
	*out = append(*out, n)
	if (n.Type == TypeObject || n.Type == TypeArray) && !n.Collapsed {
		for _, child := range n.Children {
			flattenVisible(child, out)
		}
	}
}

func (m *TreeModel) Init() tea.Cmd { return nil }

func (m *TreeModel) Update(msg tea.Msg) (SubModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 1
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.MoveCursor(1)
		case "k", "up":
			m.MoveCursor(-1)
		case "ctrl+d":
			m.MoveCursor(m.height / 2)
		case "ctrl+u":
			m.MoveCursor(-(m.height / 2))
		case "G":
			m.cursor = len(m.flat) - 1
			m.clampOffset()
		case "l", "right":
			m.Expand()
		case "enter":
			m.Toggle()
		case "h", "left":
			m.Collapse()
		}
	}
	return m, nil
}

func (m *TreeModel) MoveCursor(delta int) {
	m.cursor += delta
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(m.flat) {
		m.cursor = len(m.flat) - 1
	}
	m.clampOffset()
}

func (m *TreeModel) clampOffset() {
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+m.height {
		m.offset = m.cursor - m.height + 1
	}
	if m.offset < 0 {
		m.offset = 0
	}
}

func (m *TreeModel) Expand() {
	if m.cursor >= len(m.flat) {
		return
	}
	n := m.flat[m.cursor]
	if n.Type == TypeObject || n.Type == TypeArray {
		n.Collapsed = false
		m.rebuild()
	}
}

func (m *TreeModel) Collapse() {
	if m.cursor >= len(m.flat) {
		return
	}
	n := m.flat[m.cursor]
	if n.Type == TypeObject || n.Type == TypeArray {
		n.Collapsed = true
		m.rebuild()
	}
}

func (m *TreeModel) Toggle() {
	if m.cursor >= len(m.flat) {
		return
	}
	n := m.flat[m.cursor]
	if n.Type == TypeObject || n.Type == TypeArray {
		n.Collapsed = !n.Collapsed
		m.rebuild()
	}
}

func (m *TreeModel) JumpToTop() {
	m.cursor = 0
	m.offset = 0
}

func (m *TreeModel) JumpToBottom() {
	m.cursor = len(m.flat) - 1
	m.clampOffset()
}

func (m *TreeModel) CurrentNode() *Node {
	if m.cursor < len(m.flat) {
		return m.flat[m.cursor]
	}
	return nil
}

func (m *TreeModel) SetHighlights(h map[int]bool) {
	m.highlights = h
}

func (m *TreeModel) ClearHighlights() {
	m.highlights = map[int]bool{}
}

// FlatAll returns all nodes regardless of collapse state (used by search).
func (m *TreeModel) FlatAll() []*Node {
	var all []*Node
	flattenAll(m.root, &all)
	return all
}

func flattenAll(n *Node, out *[]*Node) {
	*out = append(*out, n)
	for _, c := range n.Children {
		flattenAll(c, out)
	}
}

func (m *TreeModel) View() string {
	var sb strings.Builder
	end := m.offset + m.height
	if end > len(m.flat) {
		end = len(m.flat)
	}
	for i := m.offset; i < end; i++ {
		line := m.renderNode(m.flat[i], i == m.cursor, m.highlights[i])
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (m *TreeModel) renderNode(n *Node, selected, highlighted bool) string {
	connector := nodeConnector(n)

	collapseIcon := "  "
	if (n.Type == TypeObject || n.Type == TypeArray) && len(n.Children) > 0 {
		if n.Collapsed {
			collapseIcon = "▶ "
		} else {
			collapseIcon = "▼ "
		}
	}

	var keyPart, valPart string
	if n.Key != "" {
		keyPart = ui.Blue.Render(n.Key) + ": "
	}

	switch n.Type {
	case TypeObject, TypeArray:
		if n.Collapsed || len(n.Children) == 0 {
			valPart = ui.Muted.Render(n.Summary())
		} else {
			if n.Type == TypeObject {
				valPart = ui.Muted.Render("{")
			} else {
				valPart = ui.Muted.Render("[")
			}
		}
	case TypeString:
		valPart = ui.Green.Render(fmt.Sprintf("%q", n.StrVal))
	case TypeNumber:
		valPart = ui.Peach.Render(n.Summary())
	case TypeBool:
		valPart = ui.Red.Render(n.Summary())
	case TypeNull:
		valPart = ui.Muted.Render("null")
	}

	line := connector + collapseIcon + keyPart + valPart

	switch {
	case selected:
		line = padRight(line, m.width)
		line = ui.Selected.Render(ansi.Strip(line))
	case highlighted:
		line = ui.Yellow.Render(ansi.Strip(line))
	}
	return line
}

// nodeConnector builds the Unicode tree connector prefix for n (e.g. "│   ├── ").
func nodeConnector(n *Node) string {
	if n.Parent == nil {
		return "" // root: no connector
	}
	// Collect ancestors between root and n's parent.
	var ancestors []*Node
	for cur := n.Parent; cur.Parent != nil; cur = cur.Parent {
		ancestors = append(ancestors, cur)
	}
	var sb strings.Builder
	// Write ancestor guides from outermost inward.
	for i := len(ancestors) - 1; i >= 0; i-- {
		if isLastChild(ancestors[i]) {
			sb.WriteString("    ")
		} else {
			sb.WriteString("│   ")
		}
	}
	// Write own branch character.
	if isLastChild(n) {
		sb.WriteString("└── ")
	} else {
		sb.WriteString("├── ")
	}
	return sb.String()
}

func isLastChild(n *Node) bool {
	if n.Parent == nil {
		return true
	}
	siblings := n.Parent.Children
	return len(siblings) > 0 && siblings[len(siblings)-1] == n
}

func padRight(s string, width int) string {
	w := lipgloss.Width(s)
	if width > w {
		return s + strings.Repeat(" ", width-w)
	}
	return s
}
