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

// FS is a fs.FS made for testing only.
type FS struct {
	fs.FS
	path string
}

// New creates a new FS using the tb.TempDir as root.
func New(tb testing.TB) FS {
	tb.Helper()

	path := tb.TempDir()
	tb.Logf("creating testFS at %s", path)
	tmpfs, err := fs.Sub(os.DirFS(path), ".")
	if err != nil {
		tb.Fatalf("failed to create test fs at %s: %s", path, err)
	}
	return FS{
		FS:   tmpfs,
		path: path,
	}
}

// WriteFile writes a file to FS.
func (t FS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filepath.Join(t.path, name), data, perm)
}

// MkdirAll creates the dir and all the necessary parents into FS.
func (t FS) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(filepath.Join(t.path, path), perm)
}
