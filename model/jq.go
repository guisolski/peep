package model

import (
	"encoding/json/v2"
	"fmt"

	"github.com/itchyny/gojq"
)

// compileJQ parses a jq expression into a query. Pure: same expr always
// yields an equivalent query or the same error.
func compileJQ(expr string) (*gojq.Query, error) {
	q, err := gojq.Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	return q, nil
}

// unmarshalJSON is the sole I/O step shared by FilterModel.Eval and
// EvalAllJQ: decoding the document jq queries run against. Callers that
// evaluate many expressions against the same document (FilterModel) should
// call this once and reuse the result instead of re-decoding per query.
func unmarshalJSON(data []byte) (interface{}, error) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return v, nil
}

// runJQFirst runs q against v and returns only the first output,
// marshaled. Pure: depends only on its arguments.
func runJQFirst(q *gojq.Query, v interface{}) ([]byte, error) {
	iter := q.Run(v)
	val, ok := iter.Next()
	if !ok {
		return []byte("null"), nil
	}
	if err, ok := val.(error); ok {
		return nil, err
	}
	return json.Marshal(val)
}

// JQResult is one output of a jq evaluation: either a marshaled value or,
// if that particular output was a runtime error, the error itself.
type JQResult struct {
	JSON []byte
	Err  error
}

// EvalAllJQ runs expr against data and returns every output the jq
// program produces, in order — unlike FilterModel.Eval, which only takes
// the first result (sufficient for a fast interactive preview), this
// mirrors real jq's behavior of emitting (and printing) every value a
// query yields, including continuing past a per-output runtime error.
func EvalAllJQ(data []byte, expr string) ([]JQResult, error) {
	q, err := compileJQ(expr)
	if err != nil {
		return nil, err
	}
	v, err := unmarshalJSON(data)
	if err != nil {
		return nil, err
	}
	var results []JQResult
	iter := q.Run(v)
	for {
		val, ok := iter.Next()
		if !ok {
			break
		}
		if e, ok := val.(error); ok {
			results = append(results, JQResult{Err: e})
			continue
		}
		b, err := json.Marshal(val)
		if err != nil {
			results = append(results, JQResult{Err: err})
			continue
		}
		results = append(results, JQResult{JSON: b})
	}
	return results, nil
}
