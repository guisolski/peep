package model

import (
	"encoding/json/jsontext"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type RawModel struct {
	vp viewport.Model
}

func NewRawModel(data []byte, width, height int) *RawModel {
	pretty := jsontext.Value(data).Clone()
	if err := pretty.Indent(jsontext.WithIndent("  ")); err != nil {
		pretty = data
	}
	vp := viewport.New(width, height-1)
	vp.SetContent(string(pretty))
	return &RawModel{vp: vp}
}

func (m *RawModel) SetContent(data []byte) {
	pretty := jsontext.Value(data).Clone()
	if err := pretty.Indent(jsontext.WithIndent("  ")); err != nil {
		pretty = data
	}
	m.vp.SetContent(string(pretty))
}

func (m *RawModel) Init() tea.Cmd { return nil }

func (m *RawModel) Update(msg tea.Msg) (SubModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.vp.Width = msg.Width
		m.vp.Height = msg.Height - 1
	}
	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m *RawModel) View() string {
	return m.vp.View()
}
