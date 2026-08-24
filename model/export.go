package model

import (
	"encoding/json/v2"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const schemaMaxStringLen = 60

// Schema returns a compact, LLM-friendly description of n's shape: field
// names, inferred types, and one example value per scalar field — no full
// data. It is not meant to be parsed back; it's a shape hint for prompting
// an LLM about a document without sending the whole payload.
func (n *Node) Schema() string {
	var sb strings.Builder
	n.writeSchema(&sb, 0)
	return sb.String()
}

func (n *Node) writeSchema(sb *strings.Builder, depth int) {
	switch n.Type {
	case TypeObject:
		if len(n.Children) == 0 {
			sb.WriteString("{}")
			return
		}
		sb.WriteString("{\n")
		for _, c := range n.Children {
			writeIndent(sb, depth+1)
			sb.WriteString(c.Key)
			sb.WriteString(": ")
			c.writeSchema(sb, depth+1)
			sb.WriteString("\n")
		}
		writeIndent(sb, depth)
		sb.WriteString("}")
	case TypeArray:
		if len(n.Children) == 0 {
			sb.WriteString("[]")
			return
		}
		types := distinctTypeNames(n.Children)
		if len(types) == 1 {
			t := n.Children[0].Type
			fmt.Fprintf(sb, "array<%s> (%d items)", typeName(t), len(n.Children))
			if t == TypeObject || t == TypeArray {
				sb.WriteString(" ")
				n.Children[0].writeSchema(sb, depth)
			}
			return
		}
		fmt.Fprintf(sb, "array<mixed: %s> (%d items)", strings.Join(types, "|"), len(n.Children))
	default:
		sb.WriteString(typeName(n.Type))
		if example := n.exampleValue(); example != "" {
			sb.WriteString(" = ")
			sb.WriteString(example)
		}
	}
}

func (n *Node) exampleValue() string {
	switch n.Type {
	case TypeString:
		r := []rune(n.StrVal)
		s := n.StrVal
		if len(r) > schemaMaxStringLen {
			s = string(r[:schemaMaxStringLen]) + "…"
		}
		return fmt.Sprintf("%q", s)
	case TypeNumber, TypeBool:
		return n.Summary()
	}
	return ""
}

func typeName(t NodeType) string {
	switch t {
	case TypeObject:
		return "object"
	case TypeArray:
		return "array"
	case TypeString:
		return "string"
	case TypeNumber:
		return "number"
	case TypeBool:
		return "bool"
	case TypeNull:
		return "null"
	}
	return "unknown"
}

func distinctTypeNames(children []*Node) []string {
	set := make(map[string]bool, len(children))
	for _, c := range children {
		set[typeName(c.Type)] = true
	}
	names := make([]string, 0, len(set))
	for name := range set {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func writeIndent(sb *strings.Builder, depth int) {
	sb.WriteString(strings.Repeat("  ", depth))
}

// CompactYAML renders n's actual values in a dense, YAML-like format: bare
// (unquoted) scalars where unambiguous, 2-space indents, no braces/commas,
// and inline [a, b, c] for scalar-only arrays. It is a one-way,
// non-round-trippable export optimized purely for LLM prompt tokens — use
// Node.ToInterface() with json.Marshal if you need real, parseable JSON.
func (n *Node) CompactYAML() string {
	var sb strings.Builder
	switch n.Type {
	case TypeObject:
		n.writeCompact(&sb, 0)
	case TypeArray:
		writeCompactArray(&sb, n, 0)
	default:
		sb.WriteString(scalarRepr(n))
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

// writeCompact writes an object's key: value lines at the given depth.
func (n *Node) writeCompact(sb *strings.Builder, depth int) {
	if len(n.Children) == 0 {
		sb.WriteString("{}\n")
		return
	}
	for _, c := range n.Children {
		writeIndent(sb, depth)
		sb.WriteString(c.Key)
		sb.WriteString(":")
		c.writeCompactValue(sb, depth)
	}
}

// writeCompactValue writes what follows "key:" for c — a bare/quoted
// scalar, an inline scalar array, or a nested block on following lines.
func (c *Node) writeCompactValue(sb *strings.Builder, parentDepth int) {
	switch c.Type {
	case TypeObject:
		if len(c.Children) == 0 {
			sb.WriteString(" {}\n")
			return
		}
		sb.WriteString("\n")
		c.writeCompact(sb, parentDepth+1)
	case TypeArray:
		if len(c.Children) == 0 {
			sb.WriteString(" []\n")
			return
		}
		if allScalars(c.Children) {
			sb.WriteString(" ")
			sb.WriteString(inlineArray(c.Children))
			sb.WriteString("\n")
			return
		}
		sb.WriteString("\n")
		writeArrayBlock(sb, c, parentDepth+1)
	default:
		sb.WriteString(" ")
		sb.WriteString(scalarRepr(c))
		sb.WriteString("\n")
	}
}

// writeCompactArray writes arr (a top-level or nested array) as either an
// inline scalar list or a "- " block, at the given depth.
func writeCompactArray(sb *strings.Builder, arr *Node, depth int) {
	if len(arr.Children) == 0 {
		sb.WriteString("[]\n")
		return
	}
	if allScalars(arr.Children) {
		sb.WriteString(inlineArray(arr.Children))
		sb.WriteString("\n")
		return
	}
	writeArrayBlock(sb, arr, depth)
}

func writeArrayBlock(sb *strings.Builder, arr *Node, depth int) {
	for _, item := range arr.Children {
		writeIndent(sb, depth)
		sb.WriteString("- ")
		item.writeArrayItem(sb, depth)
	}
}

// writeArrayItem writes one "- " list item. dashDepth is the indent depth
// of the dash itself; nested fields indent one level deeper than that.
func (item *Node) writeArrayItem(sb *strings.Builder, dashDepth int) {
	switch item.Type {
	case TypeObject:
		if len(item.Children) == 0 {
			sb.WriteString("{}\n")
			return
		}
		for i, c := range item.Children {
			if i > 0 {
				writeIndent(sb, dashDepth+1)
			}
			sb.WriteString(c.Key)
			sb.WriteString(":")
			c.writeCompactValue(sb, dashDepth+1)
		}
	case TypeArray:
		if len(item.Children) == 0 {
			sb.WriteString("[]\n")
			return
		}
		if allScalars(item.Children) {
			sb.WriteString(inlineArray(item.Children))
			sb.WriteString("\n")
			return
		}
		sb.WriteString("\n")
		writeArrayBlock(sb, item, dashDepth+1)
	default:
		sb.WriteString(scalarRepr(item))
		sb.WriteString("\n")
	}
}

func allScalars(children []*Node) bool {
	for _, c := range children {
		if c.Type == TypeObject || c.Type == TypeArray {
			return false
		}
	}
	return true
}

func inlineArray(children []*Node) string {
	parts := make([]string, len(children))
	for i, c := range children {
		parts[i] = scalarRepr(c)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func scalarRepr(n *Node) string {
	switch n.Type {
	case TypeString:
		return yamlString(n.StrVal)
	case TypeNumber, TypeBool, TypeNull:
		return n.Summary()
	}
	return ""
}

// yamlString renders s bare when unambiguous, or JSON-quoted (reusing
// Go's correct escaping) when bareness would be ambiguous in a YAML-like
// reader's eyes.
func yamlString(s string) string {
	if needsQuoting(s) {
		b, _ := json.Marshal(s)
		return string(b)
	}
	return s
}

func needsQuoting(s string) bool {
	if s == "" || strings.TrimSpace(s) != s {
		return true
	}
	if strings.ContainsAny(s, "\n:#") {
		return true
	}
	if looksLikeYAMLScalar(s) {
		return true
	}
	switch s[0] {
	case '-', '[', '{', '"', '\'', '@', '`', '|', '>', '%', '&', '*', '!':
		return true
	}
	return false
}

func looksLikeYAMLScalar(s string) bool {
	switch strings.ToLower(s) {
	case "true", "false", "null", "yes", "no", "~":
		return true
	}
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
