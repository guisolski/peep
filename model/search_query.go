package model

import (
	"strconv"
	"strings"
	"unicode"
)

// Comparator / clause kinds for English-friendly slash search.
type cmpOp int

const (
	cmpContains cmpOp = iota // field eq via : — value text contains rhs
	cmpEq
	cmpGt
	cmpGte
	cmpLt
	cmpLte
)

type clauseKind int

const (
	clauseTerm clauseKind = iota
	clauseField
)

type searchClause struct {
	kind  clauseKind
	text  string // term
	key   string // field name fragment
	op    cmpOp
	value string
}

// English filler words stripped before parsing (locale: English only).
// Short tokens like "a"/"an"/"to" are intentionally omitted so single-letter
// and common key searches still work.
var searchFillers = map[string]bool{
	"which": true, "what": true, "who": true, "where": true, "that": true,
	"have": true, "has": true, "with": true, "the": true,
	"show": true, "find": true, "get": true, "me": true, "please": true,
	"than": true,
}

// Longer phrases first so "greater than or equal" wins over "greater than".
var cmpPhrases = []struct {
	phrase string
	sym    string
}{
	{"greater than or equal", ">="},
	{"less than or equal", "<="},
	{"at least", ">="},
	{"at most", "<="},
	{"greater than", ">"},
	{"more than", ">"},
	{"bigger than", ">"},
	{"less than", "<"},
	{"smaller than", "<"},
	{"equal to", "="},
	{"equals", "="},
}

// normalizeSearchQuery lowercases, strips "?", maps English comparator
// phrases to symbols, and removes English filler words.
func normalizeSearchQuery(q string) string {
	s := strings.ToLower(strings.TrimSpace(q))
	s = strings.TrimRight(s, "?!.")
	for _, p := range cmpPhrases {
		s = strings.ReplaceAll(s, p.phrase, " "+p.sym+" ")
	}
	fields := strings.Fields(s)
	out := make([]string, 0, len(fields))
	for _, w := range fields {
		if searchFillers[w] {
			continue
		}
		out = append(out, w)
	}
	return strings.Join(out, " ")
}

// parseSearchQuery turns a raw query into AND clauses.
func parseSearchQuery(q string) []searchClause {
	s := normalizeSearchQuery(q)
	if s == "" {
		return nil
	}
	phrases, rest := extractQuotedPhrases(s)
	var clauses []searchClause
	for _, p := range phrases {
		if p == "" {
			continue
		}
		clauses = append(clauses, searchClause{kind: clauseTerm, text: p})
	}
	tokens := strings.Fields(rest)
	for i := 0; i < len(tokens); {
		if c, n := parseFieldTokens(tokens[i:]); n > 0 {
			clauses = append(clauses, c)
			i += n
			continue
		}
		clauses = append(clauses, searchClause{kind: clauseTerm, text: tokens[i]})
		i++
	}
	return clauses
}

func extractQuotedPhrases(s string) (phrases []string, rest string) {
	var b strings.Builder
	inQuote := false
	var phrase strings.Builder
	for _, r := range s {
		if r == '"' {
			if inQuote {
				phrases = append(phrases, phrase.String())
				phrase.Reset()
				inQuote = false
			} else {
				inQuote = true
			}
			continue
		}
		if inQuote {
			phrase.WriteRune(r)
		} else {
			b.WriteRune(r)
		}
	}
	if inQuote {
		// Unclosed quote: remainder is one phrase.
		phrases = append(phrases, phrase.String())
	}
	return phrases, b.String()
}

func parseFieldTokens(tokens []string) (searchClause, int) {
	if len(tokens) == 0 {
		return searchClause{}, 0
	}
	t0 := tokens[0]

	// Single token: key:value / key>value / key>=value / ...
	if c, ok := parseFieldToken(t0); ok {
		return c, 1
	}

	// key : value  or  key > value  (op as its own token)
	if len(tokens) >= 3 {
		if op, ok := parseOpToken(tokens[1]); ok && isIdent(tokens[0]) {
			return searchClause{
				kind:  clauseField,
				key:   tokens[0],
				op:    op,
				value: tokens[2],
			}, 3
		}
	}

	// key: value (colon glued to key)
	if len(tokens) >= 2 && strings.HasSuffix(t0, ":") && len(t0) > 1 {
		key := strings.TrimSuffix(t0, ":")
		if isIdent(key) {
			return searchClause{
				kind:  clauseField,
				key:   key,
				op:    cmpContains,
				value: tokens[1],
			}, 2
		}
	}

	return searchClause{}, 0
}

func parseFieldToken(tok string) (searchClause, bool) {
	for _, sym := range []string{">=", "<=", ">", "<", "=", ":"} {
		if i := strings.Index(tok, sym); i > 0 {
			key := tok[:i]
			val := tok[i+len(sym):]
			if !isIdent(key) || val == "" {
				continue
			}
			return searchClause{
				kind:  clauseField,
				key:   key,
				op:    opFromSym(sym),
				value: val,
			}, true
		}
	}
	return searchClause{}, false
}

func parseOpToken(tok string) (cmpOp, bool) {
	switch tok {
	case ":", "=":
		return cmpContains, true
	case ">":
		return cmpGt, true
	case ">=":
		return cmpGte, true
	case "<":
		return cmpLt, true
	case "<=":
		return cmpLte, true
	}
	return 0, false
}

func opFromSym(sym string) cmpOp {
	switch sym {
	case ">":
		return cmpGt
	case ">=":
		return cmpGte
	case "<":
		return cmpLt
	case "<=":
		return cmpLte
	case "=":
		return cmpEq
	default: // ":"
		return cmpContains
	}
}

func isIdent(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			continue
		}
		return false
	}
	return true
}

// nodeMatchesClauses reports whether n satisfies every clause (AND).
func nodeMatchesClauses(n *Node, clauses []searchClause) bool {
	if len(clauses) == 0 {
		return false
	}
	for _, c := range clauses {
		if !nodeMatchesClause(n, c) {
			return false
		}
	}
	return true
}

func nodeMatchesClause(n *Node, c searchClause) bool {
	switch c.kind {
	case clauseTerm:
		return nodeContainsText(n, c.text)
	case clauseField:
		if !strings.Contains(strings.ToLower(n.Key), strings.ToLower(c.key)) {
			return false
		}
		return compareNodeValue(n, c.op, c.value)
	}
	return false
}

func nodeContainsText(n *Node, text string) bool {
	if text == "" {
		return false
	}
	t := strings.ToLower(text)
	if strings.Contains(strings.ToLower(n.Key), t) {
		return true
	}
	return strings.Contains(strings.ToLower(nodeValueText(n)), t)
}

func nodeValueText(n *Node) string {
	switch n.Type {
	case TypeString:
		return n.StrVal
	case TypeNumber:
		return n.NumRaw
	default:
		return n.Summary()
	}
}

// compareNodeValue applies op against n's value. Inequalities require TypeNumber.
func compareNodeValue(n *Node, op cmpOp, raw string) bool {
	raw = strings.TrimSpace(raw)
	switch op {
	case cmpContains, cmpEq:
		return strings.Contains(strings.ToLower(nodeValueText(n)), strings.ToLower(raw))
	case cmpGt, cmpGte, cmpLt, cmpLte:
		if n.Type != TypeNumber {
			return false
		}
		rhs, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return false
		}
		lhs := n.NumVal
		switch op {
		case cmpGt:
			return lhs > rhs
		case cmpGte:
			return lhs >= rhs
		case cmpLt:
			return lhs < rhs
		case cmpLte:
			return lhs <= rhs
		}
	}
	return false
}

// expandAncestors expands every collapsed ancestor so n can appear in the
// visible tree flat list.
func expandAncestors(n *Node) {
	for p := n.Parent; p != nil; p = p.Parent {
		p.Collapsed = false
	}
}

func visibleIndex(tree *TreeModel, n *Node) int {
	for i, x := range tree.flat {
		if x == n {
			return i
		}
	}
	return -1
}

func indexInNodes(nodes []*Node, n *Node) int {
	for i, x := range nodes {
		if x == n {
			return i
		}
	}
	return -1
}
