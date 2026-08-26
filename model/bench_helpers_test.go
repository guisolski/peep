package model

import (
	"fmt"
	"strings"
)

// genWideObjectJSON builds a single JSON object with n numeric keys:
// {"key0":0,"key1":1,...}.
func genWideObjectJSON(n int) []byte {
	var sb strings.Builder
	sb.WriteByte('{')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(',')
		}
		fmt.Fprintf(&sb, "%q:%d", fmt.Sprintf("key%d", i), i)
	}
	sb.WriteByte('}')
	return []byte(sb.String())
}

// genDeepObjectJSON builds a single-key object nested depth levels deep:
// {"a":{"a":{...:1}...}}.
func genDeepObjectJSON(depth int) []byte {
	var sb strings.Builder
	for i := 0; i < depth; i++ {
		sb.WriteString(`{"a":`)
	}
	sb.WriteByte('1')
	for i := 0; i < depth; i++ {
		sb.WriteByte('}')
	}
	return []byte(sb.String())
}

// genLargeArrayJSON builds a single JSON array of n increasing integers.
func genLargeArrayJSON(n int) []byte {
	var sb strings.Builder
	sb.WriteByte('[')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(',')
		}
		fmt.Fprintf(&sb, "%d", i)
	}
	sb.WriteByte(']')
	return []byte(sb.String())
}
