package clipboard

import (
	"errors"
	"testing"
)

func TestCopy(t *testing.T) {
	orig := WriteFunc
	defer func() { WriteFunc = orig }()

	var got string
	WriteFunc = func(s string) error {
		got = s
		return nil
	}

	if err := Copy("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello" {
		t.Fatalf("WriteFunc: got %q, want %q", got, "hello")
	}
}

func TestCopy_error(t *testing.T) {
	orig := WriteFunc
	defer func() { WriteFunc = orig }()

	want := errors.New("no clipboard available")
	WriteFunc = func(s string) error {
		return want
	}

	if err := Copy("hello"); !errors.Is(err, want) {
		t.Fatalf("Copy: got %v, want %v", err, want)
	}
}
