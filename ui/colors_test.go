package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// TestStyles_Render is a smoke test: every exported style should render
// non-empty output for non-empty input without panicking.
func TestStyles_Render(t *testing.T) {
	styles := map[string]lipgloss.Style{
		"Blue":       Blue,
		"Green":      Green,
		"Peach":      Peach,
		"Red":        Red,
		"Muted":      Muted,
		"Yellow":     Yellow,
		"Text":       Text,
		"Selected":   Selected,
		"Box":        Box,
		"BoxFocused": BoxFocused,
		"StatusBar":  StatusBar,
		"ModeTag":    ModeTag,
	}
	for name, style := range styles {
		if got := style.Render("x"); got == "" {
			t.Errorf("%s.Render(%q): got empty string", name, "x")
		}
	}
}
