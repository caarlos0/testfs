package testfs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
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

	testfile := "foo/bar/foobar"
	testfile2 := "foo/bar/link-to-foobar"

	if err := tmpfs.MkdirAll(filepath.Dir(testfile), 0o764); err != nil {
		t.Fatalf("failed to create dir: %s", err)
	}
	if err := tmpfs.WriteFile(testfile, []byte("example"), 0o644); err != nil {
		t.Fatalf("failed to create dir: %s", err)
	}

	if err := tmpfs.Symlink(testfile, testfile2); err != nil {
		t.Fatalf("failed to create symlink: %s", err)
	}

	bts, err := fs.ReadFile(tmpfs, testfile)
	if err != nil {
		t.Fatalf("failed to read %s: %s", testfile, err)
	}
	if string(bts) != "example" {
		t.Fatalf("invalid %s contents, got %s, want %s", testfile, string(bts), "example")
	}

	//XXX
	//if err := fstest.TestFS(tmpfs, testfile, testfile2); err != nil {
	//	t.Fatalf("failed to check fs: %s", err)
	//}
}

func TestFSLinkProtection(t *testing.T) {
	tmpfs := New(t)

	testfile := "link-to-the-outside"

	if err := os.Symlink("./README.md", filepath.Join(tmpfs.path, testfile)); err != nil {
		t.Fatalf("failed to create symlink: %s", err)
	}

	if _, err := fs.ReadFile(tmpfs, testfile); err == nil {
		t.Fatalf("should have failed to read link to outside %s: %s", testfile, err)
	}
}

func TestInvalidPathName(t *testing.T) {
	_, err := New(t).Open("/asddas")
	pe := &os.PathError{}
	if err == nil || !errors.As(err, &pe) {
		t.Fatalf("expected a path error, got %s", err)
	}
}
