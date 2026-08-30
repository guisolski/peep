package model

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestApp_LazySubModels(t *testing.T) {
	app, err := NewApp([]byte(`{"a":1}`), "test.json", 80, 24)
	if err != nil {
		t.Fatalf("NewApp error: %v", err)
	}

	if app.raw != nil {
		t.Fatal("raw: want nil before any key, got non-nil")
	}
	if app.search != nil {
		t.Fatal("search: want nil before any key, got non-nil")
	}
	if app.filter != nil {
		t.Fatal("filter: want nil before any key, got non-nil")
	}

	cases := []struct {
		name string
		key  string
		want func(*App) bool
	}{
		{"r opens raw", "r", func(a *App) bool { return a.raw != nil }},
		{"/ opens search", "/", func(a *App) bool { return a.search != nil }},
		{":  opens filter", ":", func(a *App) bool { return a.filter != nil }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app, err := NewApp([]byte(`{"a":1}`), "test.json", 80, 24)
			if err != nil {
				t.Fatalf("NewApp error: %v", err)
			}
			model, _ := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tc.key)})
			app = model.(*App)
			if !tc.want(app) {
				t.Fatalf("key %q: expected sub-model to be constructed", tc.key)
			}
		})
	}
}

// TestApp_WindowSizeBeforeAnyKey documents the lazy-construction gotcha:
// a WindowSizeMsg (which Bubble Tea sends before any key) must not force
// eager construction of the lazily-built sub-models.
func TestApp_WindowSizeBeforeAnyKey(t *testing.T) {
	app, err := NewApp([]byte(`{"a":1}`), "test.json", 80, 24)
	if err != nil {
		t.Fatalf("NewApp error: %v", err)
	}
	model, _ := app.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	app = model.(*App)
	if app.raw != nil || app.search != nil || app.filter != nil {
		t.Fatal("WindowSizeMsg should not eagerly construct lazy sub-models")
	}
}

// TestApp_SearchTypingIgnoresModeKeys documents that mode-switch keys must
// insert into the search query while the input is focused — otherwise typing
// "id:" jumps into filter mode on the colon.
func TestApp_SearchTypingIgnoresModeKeys(t *testing.T) {
	cases := []struct {
		name string
		key  string
	}{
		{"colon stays in search", ":"},
		{"slash stays in search", "/"},
		{"r stays in search", "r"},
		{"g stays in search", "g"},
		{"q stays in search", "q"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app, err := NewApp([]byte(`{"id":1}`), "test.json", 80, 24)
			if err != nil {
				t.Fatalf("NewApp: %v", err)
			}
			model, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
			app = model.(*App)
			model, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")})
			app = model.(*App)
			model, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tc.key)})
			app = model.(*App)

			if app.mode != ModeSearch {
				t.Fatalf("mode = %v, want ModeSearch after key %q", app.mode, tc.key)
			}
			if got := app.searchModel().Query(); got != tc.key {
				t.Fatalf("query = %q, want %q", got, tc.key)
			}
		})
	}
}
