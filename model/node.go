package model

import (
	"bytes"
	"encoding/json/jsontext"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type NodeType int

const (
	TypeObject NodeType = iota
	TypeArray
	TypeString
	TypeNumber
	TypeBool
	TypeNull
)

type Node struct {
	Type   NodeType
	Key    string
	StrVal string
	// NumVal is a best-effort float64 conversion of a numeric value, kept
	// for internal numeric use. NumRaw holds the original JSON literal text
	// verbatim and is the source of truth for display/export, since NumVal
	// silently loses precision for integers beyond 2^53.
	NumVal    float64
	NumRaw    string
	BoolVal   bool
	Children  []*Node
	Parent    *Node
	Collapsed bool
	Depth     int
	IsArrItem bool // true when parent is TypeArray
}

func (n *Node) Path() string {
	if n.Parent == nil {
		return "."
	}
	parentPath := n.Parent.Path()
	if n.IsArrItem {
		base := parentPath
		if base == "." {
			base = ""
		}
		return base + fmt.Sprintf("[%s]", n.Key)
	}
	if parentPath == "." {
		return "." + n.Key
	}
	return parentPath + "." + n.Key
}

func (n *Node) Summary() string {
	switch n.Type {
	case TypeObject:
		return fmt.Sprintf("{%d}", len(n.Children))
	case TypeArray:
		return fmt.Sprintf("[%d]", len(n.Children))
	case TypeString:
		return fmt.Sprintf("%q", n.StrVal)
	case TypeNumber:
		return n.NumRaw
	case TypeBool:
		if n.BoolVal {
			return "true"
		}
		return "false"
	case TypeNull:
		return "null"
	}
	return ""
}

// MarshalJSONTo implements encoding/json/v2's MarshalerTo, writing n's
// subtree directly to enc — preserving Children order and NumRaw's exact
// numeric literal text, so json.Marshal(n, ...) round-trips big integers
// exactly instead of going through a float64-lossy intermediate value.
func (n *Node) MarshalJSONTo(enc *jsontext.Encoder) error {
	switch n.Type {
	case TypeObject:
		if err := enc.WriteToken(jsontext.BeginObject); err != nil {
			return err
		}
		for _, c := range n.Children {
			if err := enc.WriteToken(jsontext.String(c.Key)); err != nil {
				return err
			}
			if err := c.MarshalJSONTo(enc); err != nil {
				return err
			}
		}
		return enc.WriteToken(jsontext.EndObject)
	case TypeArray:
		if err := enc.WriteToken(jsontext.BeginArray); err != nil {
			return err
		}
		for _, c := range n.Children {
			if err := c.MarshalJSONTo(enc); err != nil {
				return err
			}
		}
		return enc.WriteToken(jsontext.EndArray)
	case TypeString:
		return enc.WriteToken(jsontext.String(n.StrVal))
	case TypeNumber:
		return enc.WriteValue(jsontext.Value(n.NumRaw))
	case TypeBool:
		return enc.WriteToken(jsontext.Bool(n.BoolVal))
	case TypeNull:
		return enc.WriteToken(jsontext.Null)
	}
	return fmt.Errorf("unknown node type %v", n.Type)
}

// ParseJSON decodes data in a single streaming pass directly into a *Node
// tree — no intermediate interface{} representation. Object key order
// matches the source document (never reordered), and NumRaw preserves each
// number's exact original literal text.
func ParseJSON(data []byte) (*Node, error) {
	dec := jsontext.NewDecoder(bytes.NewReader(data))
	root, err := decodeNode(dec, "", nil, 0, false)
	if err != nil {
		return nil, err
	}
	// jsontext.Decoder supports reading a stream of top-level values, but
	// peep documents must be exactly one JSON value — reject trailing
	// content the way json.Unmarshal already did.
	if tok, err := dec.ReadToken(); err == nil {
		return nil, fmt.Errorf("unexpected trailing data after top-level JSON value: %s", tok.String())
	} else if !errors.Is(err, io.EOF) {
		return nil, err
	}
	return root, nil
}

func decodeNode(dec *jsontext.Decoder, key string, parent *Node, depth int, isArrItem bool) (*Node, error) {
	tok, err := dec.ReadToken()
	if err != nil {
		return nil, err
	}
	n := &Node{Key: key, Parent: parent, Depth: depth, IsArrItem: isArrItem}
	switch tok.Kind() {
	case jsontext.KindBeginObject:
		n.Type = TypeObject
		for dec.PeekKind() != jsontext.KindEndObject {
			nameTok, err := dec.ReadToken()
			if err != nil {
				return nil, err
			}
			child, err := decodeNode(dec, nameTok.String(), n, depth+1, false)
			if err != nil {
				return nil, err
			}
			n.Children = append(n.Children, child)
		}
		if _, err := dec.ReadToken(); err != nil { // consume '}'
			return nil, err
		}
	case jsontext.KindBeginArray:
		n.Type = TypeArray
		for i := 0; dec.PeekKind() != jsontext.KindEndArray; i++ {
			child, err := decodeNode(dec, strconv.Itoa(i), n, depth+1, true)
			if err != nil {
				return nil, err
			}
			n.Children = append(n.Children, child)
		}
		if _, err := dec.ReadToken(); err != nil { // consume ']'
			return nil, err
		}
	case jsontext.KindString:
		n.Type = TypeString
		n.StrVal = tok.String()
	case jsontext.KindNumber:
		n.Type = TypeNumber
		n.NumRaw = tok.String()
		n.NumVal, _ = tok.Float() // best-effort; NumRaw is the source of truth
	case jsontext.KindTrue:
		n.Type = TypeBool
		n.BoolVal = true
	case jsontext.KindFalse:
		n.Type = TypeBool
		n.BoolVal = false
	case jsontext.KindNull:
		n.Type = TypeNull
	default:
		return nil, fmt.Errorf("unexpected token kind %v", tok.Kind())
	}
	return n, nil
}
