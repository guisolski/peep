package model

import "testing"

func TestEvalAllJQ_single(t *testing.T) {
	results, err := EvalAllJQ([]byte(`{"a":1}`), ".a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("results: got %d, want 1", len(results))
	}
	if string(results[0].JSON) != "1" {
		t.Fatalf("results[0]: got %q, want \"1\"", results[0].JSON)
	}
}

func TestEvalAllJQ_multi(t *testing.T) {
	results, err := EvalAllJQ([]byte(`[1,2,3]`), ".[]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("results: got %d, want 3", len(results))
	}
	for i, want := range []string{"1", "2", "3"} {
		if results[i].Err != nil {
			t.Fatalf("results[%d]: unexpected err %v", i, results[i].Err)
		}
		if string(results[i].JSON) != want {
			t.Fatalf("results[%d]: got %q, want %q", i, results[i].JSON, want)
		}
	}
}

func TestEvalAllJQ_empty(t *testing.T) {
	results, err := EvalAllJQ([]byte(`{"a":1}`), "empty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("results: got %d, want 0", len(results))
	}
}

func TestEvalAllJQ_parseError(t *testing.T) {
	_, err := EvalAllJQ([]byte(`{"a":1}`), "???invalid???")
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

// TestEvalAllJQ_runtimeErrorContinues establishes gojq's actual behavior
// when one output in a stream errors: does iteration stop, or continue to
// later outputs? This test documents whatever gojq v0.12.19 actually does.
func TestEvalAllJQ_runtimeErrorContinues(t *testing.T) {
	results, err := EvalAllJQ([]byte(`[1,2,3]`), `.[] | if . == 2 then error("bad") else . end`)
	if err != nil {
		t.Fatalf("unexpected top-level error: %v", err)
	}
	var gotErr bool
	var okCount int
	for _, r := range results {
		if r.Err != nil {
			gotErr = true
			continue
		}
		okCount++
	}
	if !gotErr {
		t.Fatalf("results: want at least one error output, got %+v", results)
	}
	t.Logf("gojq runtime-error-mid-stream behavior: %d results total, %d ok, error present=%v", len(results), okCount, gotErr)
}
