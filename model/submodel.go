package model

import tea "github.com/charmbracelet/bubbletea"

// SubModel is implemented by every view mode.
type SubModel interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (SubModel, tea.Cmd)
	View() string
}
