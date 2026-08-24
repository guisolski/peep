package model

import (
	"encoding/json/v2"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/guisolski/peep/ui"
)

type FilterModel struct {
	raw    *RawModel
	data   []byte
	input  textinput.Model
	errMsg string
}

func NewFilterModel(data []byte, width, height int) *FilterModel {
	ti := textinput.New()
	ti.Placeholder = "jq filter…"
	ti.Focus()
	return &FilterModel{
		raw:   NewRawModel(data, width, height-2),
		data:  data,
		input: ti,
	}
}

// Eval runs a jq expression against the stored JSON and returns the first
// result, marshaled. Used for the live filter preview; see EvalAllJQ for a
// variant that returns every output a query yields.
func (m *FilterModel) Eval(expr string) ([]byte, error) {
	q, v, err := parseJQ(m.data, expr)
	if err != nil {
		return nil, err
	}
	iter := q.Run(v)
	val, ok := iter.Next()
	if !ok {
		return []byte("null"), nil
	}
	if err, ok := val.(error); ok {
		return nil, err
	}
	return json.Marshal(val)
}

func (m *FilterModel) apply() {
	expr := m.input.Value()
	if expr == "" {
		m.raw.SetContent(m.data)
		m.errMsg = ""
		return
	}
	result, err := m.Eval(expr)
	if err != nil {
		m.errMsg = err.Error()
		return
	}
	m.errMsg = ""
	m.raw.SetContent(result)
}

func (m *FilterModel) Init() tea.Cmd { return textinput.Blink }

func (m *FilterModel) Update(msg tea.Msg) (SubModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			return m, nil
		}
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		m.apply()
		return m, cmd
	case tea.WindowSizeMsg:
		m.raw.vp.Width = msg.Width
		m.raw.vp.Height = msg.Height - 2
	}
	return m, nil
}

func (m *FilterModel) View() string {
	var sb strings.Builder
	sb.WriteString(m.raw.View())
	sb.WriteByte('\n')
	sb.WriteString(ui.Yellow.Render(":") + " " + m.input.View())
	if m.errMsg != "" {
		sb.WriteString("  " + ui.Red.Render(m.errMsg))
	}
	return sb.String()
}
