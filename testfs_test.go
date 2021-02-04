package testfs

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"testing"
	"testing/fstest"
)

var _ fs.FS = TestFS{}

func ExampleNewTestFS() {
	tmpfs := New(&testing.T{})
	_ = tmpfs.MkdirAll("foo/bar", 0o764)
	_ = tmpfs.WriteFile("foo/bar/foobar", []byte("example"), 0o644)
	bts, _ := fs.ReadFile(tmpfs, "foo/bar/foobar")
	fmt.Println(string(bts))
	//output: example
}

func TestTestFS(t *testing.T) {
	tmpfs := New(t)

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

	if err := fstest.TestFS(tmpfs, "foo/bar/foobar"); err != nil {
		t.Fatalf("failed to check fs: %s", err)
	}
}
