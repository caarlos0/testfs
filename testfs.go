// Package testfs provides a simple fs.FS which is contained in a test
// (using testing.TB's TempDir) and with a few helper methods.
//
// The temporary FS is auto-cleaned once the test and all its children finish.
package testfs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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

	path, err := filepath.EvalSymlinks(tb.TempDir())
	if err != nil {
		tb.Fatalf("failed to create testfs: %s", err)
	}
	tb.Logf("creating testFS at %s", path)
	return FS{
		FS:   os.DirFS(path),
		path: path,
	}
}

// Path returns the FS root path.
func (t FS) Path() string {
	return t.path
}

var ErrOutsideFS = fmt.Errorf("path is outside test fs root folder")

// WriteFile writes a file to FS.
func (t FS) WriteFile(name string, data []byte, perm os.FileMode) error {
	if filepath.IsAbs(name) {
		if strings.HasPrefix(name, t.path) {
			return os.WriteFile(name, data, perm)
		}
		return fmt.Errorf("%s: %w", name, ErrOutsideFS)
	}
	return os.WriteFile(filepath.Join(t.path, name), data, perm)
}

// MkdirAll creates the dir and all the necessary parents into FS.
func (t FS) MkdirAll(path string, perm os.FileMode) error {
	if filepath.IsAbs(path) {
		if strings.HasPrefix(path, t.path) {
			return os.MkdirAll(path, perm)
		}
		return fmt.Errorf("%s: %w", path, ErrOutsideFS)
	}
	return os.MkdirAll(filepath.Join(t.path, path), perm)
}
