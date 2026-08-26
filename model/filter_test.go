package model

import (
	"testing"
)

func TestFilterModel_Eval(t *testing.T) {
	data := []byte(`{"a":1,"b":[1,2,3]}`)
	fm := NewFilterModel(data, 80, 20)

	result, err := fm.Eval(".a")
	if err != nil {
		t.Fatalf("Eval .a error: %v", err)
	}
	if string(result) != "1" {
		t.Fatalf("Eval .a: got %q, want \"1\"", result)
	}
}

func TestFilterModel_EvalArray(t *testing.T) {
	data := []byte(`[1,2,3]`)
	fm := NewFilterModel(data, 80, 20)

	result, err := fm.Eval(".[1]")
	if err != nil {
		t.Fatalf("Eval .[1] error: %v", err)
	}
	if string(result) != "2" {
		t.Fatalf("Eval .[1]: got %q, want \"2\"", result)
	}
}

func TestFilterModel_EvalSyntaxError(t *testing.T) {
	data := []byte(`{"a":1}`)
	fm := NewFilterModel(data, 80, 20)
	_, err := fm.Eval("???invalid???")
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestFilterModel_Eval_InvalidDocument(t *testing.T) {
	fm := NewFilterModel([]byte(`not json`), 80, 20)
	_, err := fm.Eval(".a")
	if err == nil {
		t.Fatal("expected error for invalid document, got nil")
	}
}
