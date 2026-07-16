package model

import (
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
	// children sorted alphabetically: active, age, name, score
	if root.Children[0].Key != "active" {
		t.Fatalf("first child key: got %q, want \"active\"", root.Children[0].Key)
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
	_, err := ParseJSON([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
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
		json string
		want string
	}{
		{`"hello"`, `"hello"`},
		{`42`, "42"},
		{`3.14`, "3.14"},
		{`true`, "true"},
		{`false`, "false"},
		{`null`, "null"},
	}
	for _, tc := range cases {
		root, err := ParseJSON([]byte(tc.json))
		if err != nil {
			t.Fatalf("parse %q: %v", tc.json, err)
		}
		if got := root.Summary(); got != tc.want {
			t.Fatalf("Summary(%q): got %q, want %q", tc.json, got, tc.want)
		}
	}
}

func TestNode_ToInterface(t *testing.T) {
	data := []byte(`{"a":1,"b":[true,null]}`)
	root, err := ParseJSON(data)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	v := root.ToInterface()
	m, ok := v.(map[string]interface{})
	if !ok {
		t.Fatalf("ToInterface: got %T, want map", v)
	}
	if m["a"] != float64(1) {
		t.Fatalf("a: got %v", m["a"])
	}
	arr, ok := m["b"].([]interface{})
	if !ok || len(arr) != 2 {
		t.Fatalf("b: got %v", m["b"])
	}
}
