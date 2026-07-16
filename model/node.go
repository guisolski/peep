package model

import (
	"encoding/json"
	"fmt"
	"sort"
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
	Type      NodeType
	Key       string
	StrVal    string
	NumVal    float64
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
		if n.NumVal == float64(int64(n.NumVal)) {
			return fmt.Sprintf("%d", int64(n.NumVal))
		}
		return fmt.Sprintf("%g", n.NumVal)
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

func (n *Node) ToInterface() interface{} {
	switch n.Type {
	case TypeObject:
		m := make(map[string]interface{}, len(n.Children))
		for _, c := range n.Children {
			m[c.Key] = c.ToInterface()
		}
		return m
	case TypeArray:
		arr := make([]interface{}, len(n.Children))
		for i, c := range n.Children {
			arr[i] = c.ToInterface()
		}
		return arr
	case TypeString:
		return n.StrVal
	case TypeNumber:
		return n.NumVal
	case TypeBool:
		return n.BoolVal
	case TypeNull:
		return nil
	}
	return nil
}

func ParseJSON(data []byte) (*Node, error) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return buildNode(v, "", nil, 0, false), nil
}

func buildNode(v interface{}, key string, parent *Node, depth int, isArrItem bool) *Node {
	n := &Node{Key: key, Parent: parent, Depth: depth, IsArrItem: isArrItem}
	switch val := v.(type) {
	case map[string]interface{}:
		n.Type = TypeObject
		for k, child := range val {
			n.Children = append(n.Children, buildNode(child, k, n, depth+1, false))
		}
		sort.Slice(n.Children, func(i, j int) bool {
			return n.Children[i].Key < n.Children[j].Key
		})
	case []interface{}:
		n.Type = TypeArray
		for i, child := range val {
			n.Children = append(n.Children, buildNode(child, fmt.Sprintf("%d", i), n, depth+1, true))
		}
	case string:
		n.Type = TypeString
		n.StrVal = val
	case float64:
		n.Type = TypeNumber
		n.NumVal = val
	case bool:
		n.Type = TypeBool
		n.BoolVal = val
	case nil:
		n.Type = TypeNull
	}
	return n
}
