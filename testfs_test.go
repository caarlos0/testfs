package testfs

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"testing"
	"testing/fstest"
)

var _ fs.FS = FS{}

func ExampleNew() {
	tmpfs := New(&testing.T{})
	_ = tmpfs.MkdirAll("foo/bar", 0o764)
	_ = tmpfs.WriteFile("foo/bar/foobar", []byte("example"), 0o644)
	bts, _ := fs.ReadFile(tmpfs, "foo/bar/foobar")
	fmt.Println(string(bts))
	//output: example
}

func TestFS(t *testing.T) {
	tmpfs := New(t)

	if tmpfs.Path() != tmpfs.path {
		t.Fatalf("expected Path to be %s, got %s", tmpfs.path, tmpfs.Path())
	}

	content := "example"
	testfile := "foo/bar/foobar"
	testfile2 := "asd/sada/aaaa"

	if err := tmpfs.MkdirAll("/tmp/asd", 0o700); err == nil {
		t.Fatalf("expected to fail")
	}

	if err := tmpfs.WriteFile("/tmp/asd", []byte(content), 0o700); err == nil {
		t.Fatalf("expected to fail")
	}

	if err := tmpfs.MkdirAll(filepath.Dir(testfile), 0o764); err != nil {
		t.Fatalf("failed to create dir: %s", err)
	}
	if err := tmpfs.WriteFile(testfile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create dir: %s", err)
	}

	if err := tmpfs.MkdirAll(filepath.Dir(filepath.Join(tmpfs.Path(), testfile2)), 0o764); err != nil {
		t.Fatalf("failed to create dir: %s", err)
	}
	if err := tmpfs.WriteFile(filepath.Join(tmpfs.Path(), testfile2), []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create dir: %s", err)
	}

	bts, err := fs.ReadFile(tmpfs, testfile)
	if err != nil {
		t.Fatalf("failed to read %s: %s", testfile, err)
	}
	if string(bts) != content {
		t.Fatalf("invalid %s contents, got %s, want %s", testfile, string(bts), content)
	}

	bts, err = fs.ReadFile(tmpfs, testfile2)
	if err != nil {
		t.Fatalf("failed to read %s: %s", testfile2, err)
	}
	if string(bts) != content {
		t.Fatalf("invalid %s contents, got %s, want %s", testfile2, string(bts), content)
	}

	if err := fstest.TestFS(tmpfs, testfile, testfile2); err != nil {
		t.Fatalf("failed to check fs: %s", err)
	}
}
