package ui

import "github.com/charmbracelet/lipgloss"

var (
	Blue   = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "4", Dark: "12"})
	Green  = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "2", Dark: "10"})
	Peach  = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "3", Dark: "11"})
	Red    = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "1", Dark: "9"})
	Muted  = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "8", Dark: "8"})
	Yellow = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "3", Dark: "11"})
	Text   = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "15"})

	Selected = lipgloss.NewStyle().
			Background(lipgloss.AdaptiveColor{Light: "253", Dark: "237"})

	Box = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "4", Dark: "12"}).
		Padding(0, 1)

	BoxFocused = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "3", Dark: "11"}).
			Padding(0, 1)

	StatusBar = lipgloss.NewStyle().
			Background(lipgloss.AdaptiveColor{Light: "254", Dark: "235"}).
			Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "15"})

	ModeTag = lipgloss.NewStyle().
		Background(lipgloss.AdaptiveColor{Light: "4", Dark: "12"}).
		Foreground(lipgloss.AdaptiveColor{Light: "15", Dark: "0"}).
		Bold(true).
		Padding(0, 1)
)
