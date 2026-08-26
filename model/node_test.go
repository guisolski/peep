package model

import (
	"encoding/json/v2"
	"testing"
)

func TestParseJSON_object(t *testing.T) {
	data := []byte(`{"name":"Alice","age":30,"active":true,"score":null}`)
	root, err := ParseJSON(data)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if root.Type != TypeObject {
		t.Fatalf("root type: got %v, want TypeObject", root.Type)
	}
	if len(root.Children) != 4 {
		t.Fatalf("children: got %d, want 4", len(root.Children))
	}
	// children preserve source document order: name, age, active, score
	want := []string{"name", "age", "active", "score"}
	for i, k := range want {
		if root.Children[i].Key != k {
			t.Fatalf("child %d key: got %q, want %q", i, root.Children[i].Key, k)
		}
	}
}

func TestParseJSON_array(t *testing.T) {
	data := []byte(`[1,2,3]`)
	root, err := ParseJSON(data)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if root.Type != TypeArray {
		t.Fatalf("root type: got %v, want TypeArray", root.Type)
	}
	if len(root.Children) != 3 {
		t.Fatalf("children: got %d, want 3", len(root.Children))
	}
	if root.Children[0].NumVal != 1 {
		t.Fatalf("first child value: got %v, want 1", root.Children[0].NumVal)
	}
}

func TestParseJSON_invalid(t *testing.T) {
	cases := []struct {
		name string
		json string
	}{
		{"not json at all", `not json`},
		{"trailing garbage after object", `{"a":1} garbage`},
		{"trailing second top-level value", `{"a":1}{"b":2}`},
		{"unterminated object", `{"a":1`},
		{"unterminated array", `[1,2,3`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseJSON([]byte(tc.json))
			if err == nil {
				t.Fatalf("expected error for %q, got nil", tc.json)
			}
		})
	}
}

func TestParseJSON_NumRaw(t *testing.T) {
	cases := []struct {
		name       string
		json       string
		wantNumRaw string
		wantNumVal float64
	}{
		{"small int", `42`, "42", 42},
		{"negative", `-17`, "-17", -17},
		{"float trailing zeros preserved", `3.14000`, "3.14000", 3.14},
		{"scientific notation preserved", `1e10`, "1e10", 1e10},
		{"big integer beyond float64 precision", `9223372036854775807`, "9223372036854775807", 9223372036854775807},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, err := ParseJSON([]byte(tc.json))
			if err != nil {
				t.Fatalf("parse %q: %v", tc.json, err)
			}
			if root.NumRaw != tc.wantNumRaw {
				t.Fatalf("NumRaw: got %q, want %q", root.NumRaw, tc.wantNumRaw)
			}
			if root.NumVal != tc.wantNumVal {
				t.Fatalf("NumVal: got %v, want %v", root.NumVal, tc.wantNumVal)
			}
		})
	}
}

func TestNode_Path(t *testing.T) {
	data := []byte(`{"address":{"city":"BA"}}`)
	root, err := ParseJSON(data)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	addr := root.Children[0] // "address"
	if addr.Path() != ".address" {
		t.Fatalf("address path: got %q, want \".address\"", addr.Path())
	}
	city := addr.Children[0] // "city"
	if city.Path() != ".address.city" {
		t.Fatalf("city path: got %q, want \".address.city\"", city.Path())
	}
}

func TestNode_PathArray(t *testing.T) {
	data := []byte(`{"items":[10,20]}`)
	root, err := ParseJSON(data)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	items := root.Children[0] // "items"
	first := items.Children[0]
	if first.Path() != ".items[0]" {
		t.Fatalf("array path: got %q, want \".items[0]\"", first.Path())
	}
}

func TestNode_Summary(t *testing.T) {
	cases := []struct {
		name string
		json string
		want string
	}{
		{"string", `"hello"`, `"hello"`},
		{"int", `42`, "42"},
		{"float", `3.14`, "3.14"},
		{"true", `true`, "true"},
		{"false", `false`, "false"},
		{"null", `null`, "null"},
		{"trailing zeros preserved", `10.50000`, "10.50000"},
		{"big integer preserved", `9223372036854775807`, "9223372036854775807"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, err := ParseJSON([]byte(tc.json))
			if err != nil {
				t.Fatalf("parse %q: %v", tc.json, err)
			}
			if got := root.Summary(); got != tc.want {
				t.Fatalf("Summary(%q): got %q, want %q", tc.json, got, tc.want)
			}
		})
	}
}

func TestNode_MarshalJSONTo(t *testing.T) {
	cases := []struct {
		name string
		json string
		want string
	}{
		{"object preserves input key order", `{"b":1,"a":2}`, `{"b":1,"a":2}`},
		{"nested array and object", `{"a":1,"b":[true,null,"x"]}`, `{"a":1,"b":[true,null,"x"]}`},
		{"big integer round-trips exactly", `{"id":9223372036854775807}`, `{"id":9223372036854775807}`},
		{"float literal preserved verbatim", `{"f":3.14000}`, `{"f":3.14000}`},
		{"empty object and array", `{"o":{},"a":[]}`, `{"o":{},"a":[]}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, err := ParseJSON([]byte(tc.json))
			if err != nil {
				t.Fatalf("parse %q: %v", tc.json, err)
			}
			b, err := json.Marshal(root)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if got := string(b); got != tc.want {
				t.Fatalf("MarshalJSONTo(%q):\ngot:  %q\nwant: %q", tc.json, got, tc.want)
			}
		})
	}
}
