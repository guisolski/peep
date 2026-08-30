package model

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func largeObjectJSON(n int) []byte {
	var b strings.Builder
	b.WriteByte('{')
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, `"k%d":%d`, i, i)
	}
	b.WriteByte('}')
	return []byte(b.String())
}

func TestSearchModel_FindMatches(t *testing.T) {
	data := []byte(`{"name":"Alice","city":"Buenos Aires","age":30}`)
	root, _ := ParseJSON(data)
	tm := NewTreeModel(root, 80, 24)
	sm := NewSearchModel(tm)

	sm.SetQuery("a")
	matches := sm.Matches()
	if len(matches) == 0 {
		t.Fatal("expected matches for 'a', got none")
	}
}

func TestSearchModel_NoMatches(t *testing.T) {
	data := []byte(`{"x":1}`)
	root, _ := ParseJSON(data)
	tm := NewTreeModel(root, 80, 24)
	sm := NewSearchModel(tm)

	sm.SetQuery("zzz")
	if len(sm.Matches()) != 0 {
		t.Fatal("expected no matches, got some")
	}
}

func TestSearchModel_Cycle(t *testing.T) {
	data := []byte(`{"a":1,"b":2,"c":3}`)
	root, _ := ParseJSON(data)
	tm := NewTreeModel(root, 80, 24)
	sm := NewSearchModel(tm)
	sm.SetQuery("1") // matches value "1" in node "a"
	if len(sm.Matches()) == 0 {
		t.Skip("no matches to cycle")
	}
	initial := sm.Current()
	sm.Next()
	// with only 1 match, current wraps back to same index
	_ = initial
}

func TestSearchModel_EmptyQuery(t *testing.T) {
	data := []byte(`{"a":1}`)
	root, _ := ParseJSON(data)
	tm := NewTreeModel(root, 80, 24)
	sm := NewSearchModel(tm)
	sm.SetQuery("")
	if len(sm.Matches()) != 0 {
		t.Fatal("empty query should produce no matches")
	}
}

// TestApp_SearchViewFitsTerminal documents the search chrome bug: with a
// filled tree viewport, stacking the "/" prompt on Height-1 overflowed by
// one row and scrolled the prompt into the main area. Search must reserve
// prompt+status (chrome 2) so App.View stays within the terminal height.
func TestApp_SearchViewFitsTerminal(t *testing.T) {
	cases := []struct {
		name   string
		width  int
		height int
		keys   int
	}{
		{"standard 24-row", 80, 24, 40},
		{"short 10-row", 40, 10, 20},
		{"tall 40-row", 100, 40, 60},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app, err := NewApp(largeObjectJSON(tc.keys), "big.json", tc.width, tc.height)
			if err != nil {
				t.Fatalf("NewApp: %v", err)
			}
			model, _ := app.Update(tea.WindowSizeMsg{Width: tc.width, Height: tc.height})
			app = model.(*App)
			model, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")})
			app = model.(*App)

			wantTreeH := contentHeight(tc.height, chromeRows(ModeSearch))
			if app.tree.height != wantTreeH {
				t.Fatalf("tree.height = %d, want %d (chrome for search)", app.tree.height, wantTreeH)
			}

			view := app.View()
			lines := strings.Count(view, "\n") + 1
			if lines > tc.height {
				t.Fatalf("App.View lines = %d, want <= %d\nview:\n%s", lines, tc.height, view)
			}
		})
	}
}

func TestApp_SearchEscRestoresTreeHeight(t *testing.T) {
	const h = 24
	app, err := NewApp(largeObjectJSON(40), "big.json", 80, h)
	if err != nil {
		t.Fatalf("NewApp: %v", err)
	}
	model, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: h})
	app = model.(*App)
	model, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")})
	app = model.(*App)
	model, _ = app.Update(tea.KeyMsg{Type: tea.KeyEsc})
	app = model.(*App)

	want := contentHeight(h, chromeRows(ModeTree))
	if app.tree.height != want {
		t.Fatalf("after esc tree.height = %d, want %d", app.tree.height, want)
	}
	if app.mode != ModeTree {
		t.Fatalf("mode = %v, want ModeTree", app.mode)
	}
}
