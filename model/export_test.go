package model

import (
	"strings"
	"testing"
)

func TestNode_Schema(t *testing.T) {
	cases := []struct {
		name string
		json string
		want string
	}{
		{
			"nested object",
			`{"user":{"name":"Alice","age":30}}`,
			"{\n  user: {\n    name: string = \"Alice\"\n    age: number = 30\n  }\n}",
		},
		{
			"array of scalars",
			`{"tags":["a","b","c"]}`,
			"{\n  tags: array<string> (3 items)\n}",
		},
		{
			"array of objects",
			`{"items":[{"id":1},{"id":2}]}`,
			"{\n  items: array<object> (2 items) {\n    id: number = 1\n  }\n}",
		},
		{
			"mixed-type array",
			`{"mixed":[1,"a"]}`,
			"{\n  mixed: array<mixed: number|string> (2 items)\n}",
		},
		{
			"empty containers",
			`{"empty_obj":{},"empty_arr":[]}`,
			"{\n  empty_obj: {}\n  empty_arr: []\n}",
		},
		{
			"null scalar root",
			`null`,
			"null",
		},
		{
			"big integer field preserves exact literal",
			`{"id":9223372036854775807}`,
			"{\n  id: number = 9223372036854775807\n}",
		},
		{
			"float trailing zeros preserved",
			`{"price":10.50000}`,
			"{\n  price: number = 10.50000\n}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, err := ParseJSON([]byte(tc.json))
			if err != nil {
				t.Fatalf("parse %q: %v", tc.json, err)
			}
			if got := root.Schema(); got != tc.want {
				t.Fatalf("Schema(%q):\ngot:  %q\nwant: %q", tc.json, got, tc.want)
			}
		})
	}
}

func TestNode_Schema_LongStringTruncated(t *testing.T) {
	long := strings.Repeat("x", 70)
	root, err := ParseJSON([]byte(`"` + long + `"`))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	got := root.Schema()
	if !strings.Contains(got, "…") {
		t.Fatalf("Schema: want truncation marker, got %q", got)
	}
	if strings.Contains(got, long) {
		t.Fatalf("Schema: want truncated string, got full string: %q", got)
	}
}

func TestNode_CompactYAML(t *testing.T) {
	cases := []struct {
		name string
		json string
		want string
	}{
		{
			"bare vs quoted string",
			`{"name":"Alice","note":"needs: quoting"}`,
			"name: Alice\nnote: \"needs: quoting\"",
		},
		{
			"inline scalar array",
			`{"nums":[1,2,3]}`,
			"nums: [1, 2, 3]",
		},
		{
			"block array of objects",
			`{"items":[{"id":1},{"id":2}]}`,
			"items:\n  - id: 1\n  - id: 2",
		},
		{
			"nested object, numeric-looking string quoted",
			`{"address":{"city":"BA","zip":"1000"}}`,
			"address:\n  city: BA\n  zip: \"1000\"",
		},
		{
			"root array of scalars inline",
			`[1,2,3]`,
			"[1, 2, 3]",
		},
		{
			"empty object",
			`{}`,
			"{}",
		},
		{
			"empty array",
			`[]`,
			"[]",
		},
		{
			"big integer preserves exact literal",
			`{"id":9223372036854775807}`,
			"id: 9223372036854775807",
		},
		{
			"float trailing zeros preserved",
			`{"price":10.50000}`,
			"price: 10.50000",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, err := ParseJSON([]byte(tc.json))
			if err != nil {
				t.Fatalf("parse %q: %v", tc.json, err)
			}
			if got := root.CompactYAML(); got != tc.want {
				t.Fatalf("CompactYAML(%q):\ngot:  %q\nwant: %q", tc.json, got, tc.want)
			}
		})
	}
}
