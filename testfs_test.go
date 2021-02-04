package testfs

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"testing"
)

var _ fs.FS = TestFS{}

func ExampleNewTestFS() {
	tmpfs := NewTestFS(&testing.T{})
	_ = tmpfs.MkdirAll("foo/bar", 0o764)
	_ = tmpfs.WriteFile("foo/bar/foobar", []byte("example"), 0o644)
	bts, _ := fs.ReadFile(tmpfs, "foo/bar/foobar")
	fmt.Println(string(bts))
	//output: example
}

func TestTestFS(t *testing.T) {
	tmpfs := NewTestFS(t)

	testfile := "foo/bar/foobar"

	if err := tmpfs.MkdirAll(filepath.Dir(testfile), 0o764); err != nil {
		t.Fatalf("failed to create dir: %s", err)
	}
	if err := tmpfs.WriteFile(testfile, []byte("example"), 0o644); err != nil {
		t.Fatalf("failed to create dir: %s", err)
	}

	bts, err := fs.ReadFile(tmpfs, testfile)
	if err != nil {
		t.Fatalf("failed to read %s: %s", testfile, err)
	}
	if string(bts) != "example" {
		t.Fatalf("invalid %s contents, got %s, want %s", testfile, string(bts), "example")
	}

	var paths []string
	if err := fs.WalkDir(tmpfs, ".", func(path string, d fs.DirEntry, err error) error {
		paths = append(paths, path)
		return nil
	}); err != nil {
		t.Fatalf("failed to walk fs: %s", err)
	}

	expectedPaths := []string{".", "foo", "foo/bar", "foo/bar/foobar"}
	if len(paths) != len(expectedPaths) {
		t.Fatalf("expected %d paths, got %d. Paths: %s", len(expectedPaths), len(paths), paths)
	}
	for i := range expectedPaths {
		if paths[i] != expectedPaths[i] {
			t.Fatalf("expected paths[%d] to be %s, was %s", i, expectedPaths[i], paths[i])
		}
	}
}
