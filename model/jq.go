package model

import (
	"encoding/json"
	"fmt"

	"github.com/itchyny/gojq"
)

// parseJQ parses expr and unmarshals data, the common first step shared by
// FilterModel.Eval (first-result-only, for the live filter view) and
// EvalAllJQ (all-results, for the headless --query flag).
func parseJQ(data []byte, expr string) (*gojq.Query, interface{}, error) {
	q, err := gojq.Parse(expr)
	if err != nil {
		return nil, nil, fmt.Errorf("parse: %w", err)
	}
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, nil, fmt.Errorf("unmarshal: %w", err)
	}
	return q, v, nil
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
	q, v, err := parseJQ(data, expr)
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
