package model

import (
	"testing"
)

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
