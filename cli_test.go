package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestHeadlessMode_none(t *testing.T) {
	mode, expr, err := headlessMode("", false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mode != "" || expr != "" {
		t.Fatalf("got mode=%q expr=%q, want both empty", mode, expr)
	}
}

func TestHeadlessMode_mutuallyExclusive(t *testing.T) {
	if _, _, err := headlessMode(".a", true, false); err == nil {
		t.Fatal("expected error when --query and --schema are both set, got nil")
	}
	if _, _, err := headlessMode("", true, true); err == nil {
		t.Fatal("expected error when --schema and --llm are both set, got nil")
	}
}

func TestRunHeadless_query(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := runHeadless([]byte(`[1,2,3]`), "query", ".[]", &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code: got %d, want 0 (stderr: %s)", code, stderr.String())
	}
	if stdout.String() != "1\n2\n3\n" {
		t.Fatalf("stdout: got %q", stdout.String())
	}
}

func TestRunHeadless_queryRuntimeError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := runHeadless([]byte(`[1,2,3]`), "query", `.[] | if . == 2 then error("bad") else . end`, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("exit code: got %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "bad") {
		t.Fatalf("stderr: got %q, want mention of \"bad\"", stderr.String())
	}
}

func TestRunHeadless_schema(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := runHeadless([]byte(`{"a":1}`), "schema", "", &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code: got %d, want 0 (stderr: %s)", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "a: number = 1") {
		t.Fatalf("stdout: got %q", stdout.String())
	}
}

func TestRunHeadless_llm(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := runHeadless([]byte(`{"a":1}`), "llm", "", &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit code: got %d, want 0 (stderr: %s)", code, stderr.String())
	}
	if strings.TrimSpace(stdout.String()) != "a: 1" {
		t.Fatalf("stdout: got %q", stdout.String())
	}
}
