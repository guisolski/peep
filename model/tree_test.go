package model

import (
	"testing"
)

func TestTreeModel_flatten(t *testing.T) {
	data := []byte(`{"a":1,"b":{"c":2}}`)
	root, err := ParseJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	tm := NewTreeModel(root, 80, 24)
	// flat: root + "a" + "b" (b is collapsed by default) = 3
	if len(tm.flat) != 3 {
		t.Fatalf("flat length: got %d, want 3", len(tm.flat))
	}
}

func TestTreeModel_expand(t *testing.T) {
	data := []byte(`{"a":1,"b":{"c":2}}`)
	root, err := ParseJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	tm := NewTreeModel(root, 80, 24)
	// cursor=2 is "b" object; expand it
	tm.cursor = 2
	tm.Expand()
	if len(tm.flat) != 4 {
		t.Fatalf("flat after expand: got %d, want 4", len(tm.flat))
	}
}

func TestTreeModel_collapse(t *testing.T) {
	data := []byte(`{"a":1,"b":{"c":2}}`)
	root, _ := ParseJSON(data)
	tm := NewTreeModel(root, 80, 24)
	tm.cursor = 2
	tm.Expand()
	tm.cursor = 2
	tm.Collapse()
	if len(tm.flat) != 3 {
		t.Fatalf("flat after collapse: got %d, want 3", len(tm.flat))
	}
}

func TestTreeModel_moveCursor(t *testing.T) {
	data := []byte(`{"a":1,"b":2,"c":3}`)
	root, _ := ParseJSON(data)
	tm := NewTreeModel(root, 80, 24)
	// flat: root, a, b, c = 4 nodes
	tm.MoveCursor(2)
	if tm.cursor != 2 {
		t.Fatalf("cursor: got %d, want 2", tm.cursor)
	}
	tm.MoveCursor(10) // clamp to last
	if tm.cursor != 3 {
		t.Fatalf("cursor after clamp: got %d, want 3", tm.cursor)
	}
	tm.MoveCursor(-100) // clamp to 0
	if tm.cursor != 0 {
		t.Fatalf("cursor after neg clamp: got %d, want 0", tm.cursor)
	}
}

func TestTreeModel_JumpToTop(t *testing.T) {
	data := []byte(`[1,2,3]`)
	root, _ := ParseJSON(data)
	tm := NewTreeModel(root, 80, 24)
	tm.MoveCursor(3)
	tm.JumpToTop()
	if tm.cursor != 0 {
		t.Fatalf("JumpToTop: cursor %d, want 0", tm.cursor)
	}
}

func TestTreeModel_CurrentNode(t *testing.T) {
	data := []byte(`{"x":99}`)
	root, _ := ParseJSON(data)
	tm := NewTreeModel(root, 80, 24)
	tm.MoveCursor(1) // "x" node
	n := tm.CurrentNode()
	if n == nil || n.Key != "x" {
		t.Fatalf("CurrentNode: got %v", n)
	}
}

func TestTreeModel_SetHighlights(t *testing.T) {
	data := []byte(`{"a":1,"b":2}`)
	root, _ := ParseJSON(data)
	tm := NewTreeModel(root, 80, 24)
	tm.SetHighlights(map[int]bool{1: true})
	if !tm.highlights[1] {
		t.Fatal("highlight 1 not set")
	}
	tm.ClearHighlights()
	if len(tm.highlights) != 0 {
		t.Fatal("highlights not cleared")
	}
}
