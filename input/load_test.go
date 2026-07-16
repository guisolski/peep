package input

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")
	want := []byte(`{"a":1}`)
	if err := os.WriteFile(path, want, 0644); err != nil {
		t.Fatal(err)
	}

	res, err := Load([]string{path}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(res.Data) != string(want) {
		t.Fatalf("got %q, want %q", res.Data, want)
	}
	if res.Source != SourceFile {
		t.Fatalf("got source %v, want SourceFile", res.Source)
	}
	if res.Name != path {
		t.Fatalf("got name %q, want %q", res.Name, path)
	}
}

func TestLoadFileMissing(t *testing.T) {
	_, err := Load([]string{"/nonexistent/file.json"}, false)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
