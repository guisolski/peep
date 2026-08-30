package model

import (
	"reflect"
	"testing"
)

func TestNormalizeSearchQuery(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"which have id 2", "Which have id 2?", "id 2"},
		{"what has id 2", "what has id 2", "id 2"},
		{"greater than phrase", "value greater than 5", "value > 5"},
		{"more than", "id more than 10", "id > 10"},
		{"less than", "age less than 3", "age < 3"},
		{"at least", "score at least 7", "score >= 7"},
		{"preserves non-english", "quais id 2", "quais id 2"},
		{"filler only", "which have the?", ""},
		{"already compact", "id > 5", "id > 5"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeSearchQuery(tc.in)
			if got != tc.want {
				t.Fatalf("normalizeSearchQuery(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestParseSearchQuery(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []searchClause
	}{
		{"empty", "", nil},
		{"single term", "alice", []searchClause{{kind: clauseTerm, text: "alice"}}},
		{"and terms", "id 2", []searchClause{
			{kind: clauseTerm, text: "id"},
			{kind: clauseTerm, text: "2"},
		}},
		{"field colon", "id:2", []searchClause{
			{kind: clauseField, key: "id", op: cmpContains, value: "2"},
		}},
		{"field colon spaced", "id: 2", []searchClause{
			{kind: clauseField, key: "id", op: cmpContains, value: "2"},
		}},
		{"field gt compact", "id>5", []searchClause{
			{kind: clauseField, key: "id", op: cmpGt, value: "5"},
		}},
		{"field gt spaced", "id > 5", []searchClause{
			{kind: clauseField, key: "id", op: cmpGt, value: "5"},
		}},
		{"english nl eq", "which have id 2?", []searchClause{
			{kind: clauseTerm, text: "id"},
			{kind: clauseTerm, text: "2"},
		}},
		{"english nl gt", "value greater than 5", []searchClause{
			{kind: clauseField, key: "value", op: cmpGt, value: "5"},
		}},
		{"quoted phrase", `"via place"`, []searchClause{
			{kind: clauseTerm, text: "via place"},
		}},
		{"gte", "score >= 7", []searchClause{
			{kind: clauseField, key: "score", op: cmpGte, value: "7"},
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseSearchQuery(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("parseSearchQuery(%q)\n got %#v\nwant %#v", tc.in, got, tc.want)
			}
		})
	}
}

func TestNodeMatchesClauses(t *testing.T) {
	num := func(key, raw string, v float64) *Node {
		return &Node{Type: TypeNumber, Key: key, NumRaw: raw, NumVal: v}
	}
	str := func(key, val string) *Node {
		return &Node{Type: TypeString, Key: key, StrVal: val}
	}

	cases := []struct {
		name    string
		node    *Node
		clauses []searchClause
		want    bool
	}{
		{
			"empty clauses",
			num("id", "2", 2),
			nil,
			false,
		},
		{
			"field contains id 2",
			num("id", "2", 2),
			[]searchClause{{kind: clauseField, key: "id", op: cmpContains, value: "2"}},
			true,
		},
		{
			"field wrong value",
			num("id", "2", 2),
			[]searchClause{{kind: clauseField, key: "id", op: cmpContains, value: "9"}},
			false,
		},
		{
			"terms and on same node",
			num("id", "2", 2),
			[]searchClause{
				{kind: clauseTerm, text: "id"},
				{kind: clauseTerm, text: "2"},
			},
			true,
		},
		{
			"gt hit",
			num("id", "10", 10),
			[]searchClause{{kind: clauseField, key: "id", op: cmpGt, value: "5"}},
			true,
		},
		{
			"gt miss",
			num("id", "3", 3),
			[]searchClause{{kind: clauseField, key: "id", op: cmpGt, value: "5"}},
			false,
		},
		{
			"gt on string false",
			str("id", "10"),
			[]searchClause{{kind: clauseField, key: "id", op: cmpGt, value: "5"}},
			false,
		},
		{
			"lt hit",
			num("age", "2", 2),
			[]searchClause{{kind: clauseField, key: "age", op: cmpLt, value: "5"}},
			true,
		},
		{
			"term in string value",
			str("title", "reprehenderit"),
			[]searchClause{{kind: clauseTerm, text: "repre"}},
			true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := nodeMatchesClauses(tc.node, tc.clauses)
			if got != tc.want {
				t.Fatalf("nodeMatchesClauses = %v, want %v", got, tc.want)
			}
		})
	}
}
