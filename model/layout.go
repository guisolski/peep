package model

// contentHeight returns rows available for scrollable content given the
// terminal height and chrome rows stacked below the viewport (status bar,
// prompt line, etc.). Always at least 1 so viewports stay usable.
func contentHeight(termHeight, chromeRows int) int {
	h := termHeight - chromeRows
	if h < 1 {
		return 1
	}
	return h
}

// chromeRows returns how many rows below the viewport a mode needs.
// Search and filter stack a prompt above the status bar (2); other modes
// only reserve the status bar (1).
func chromeRows(mode Mode) int {
	switch mode {
	case ModeSearch, ModeFilter:
		return 2
	default:
		return 1
	}
}
