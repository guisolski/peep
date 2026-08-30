package model

import "testing"

func TestContentHeight(t *testing.T) {
	cases := []struct {
		name       string
		termHeight int
		chromeRows int
		want       int
	}{
		{"status only", 24, 1, 23},
		{"prompt and status", 24, 2, 22},
		{"exact fit clamps", 2, 2, 1},
		{"undersized clamps", 1, 2, 1},
		{"zero term clamps", 0, 1, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := contentHeight(tc.termHeight, tc.chromeRows)
			if got != tc.want {
				t.Fatalf("contentHeight(%d, %d) = %d, want %d",
					tc.termHeight, tc.chromeRows, got, tc.want)
			}
		})
	}
}

func TestChromeRows(t *testing.T) {
	cases := []struct {
		name string
		mode Mode
		want int
	}{
		{"tree", ModeTree, 1},
		{"graph", ModeGraph, 1},
		{"raw", ModeRaw, 1},
		{"search", ModeSearch, 2},
		{"filter", ModeFilter, 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := chromeRows(tc.mode)
			if got != tc.want {
				t.Fatalf("chromeRows(%v) = %d, want %d", tc.mode, got, tc.want)
			}
		})
	}
}
