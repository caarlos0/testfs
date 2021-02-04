// Package testfs provides a simple fs.FS which is contained in a test
// (using testing.TB's TempDir) and with a few helper methods.
//
// The temporary FS is auto-cleaned once the test and all its children finish.
package testfs

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// TestFS is a fs.FS made for testing only.
type TestFS struct {
	fs.FS
	path string
}

// NewTestFS creates a new fs.FS using the tb.TempDir as root.
func NewTestFS(tb testing.TB) TestFS {
	tb.Helper()

	path := tb.TempDir()
	tb.Logf("creating testFS at %s", path)
	tmpfs, err := fs.Sub(os.DirFS(path), ".")
	if err != nil {
		tb.Fatalf("failed to create test fs at %s: %s", path, err)
	}
	return TestFS{
		FS:   tmpfs,
		path: path,
	}
}

// WriteFile writes a file to TestFS.
func (t TestFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filepath.Join(t.path, name), data, perm)
}

// MkdirAll creates the dir and all the necessary parents into TestFS.
func (t TestFS) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(filepath.Join(t.path, path), perm)
}
