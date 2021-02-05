// Package testfs provides a simple fs.FS which is contained in a test
// (using testing.TB's TempDir) and with a few helper methods.
//
// The temporary FS is auto-cleaned once the test and all its children finish.
package testfs

import (
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

	path := tb.TempDir()
	tb.Logf("creating testFS at %s", path)
	return FS{
		FS:   dirFS(path),
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

func (t FS) Symlink(oldname, newname string) error {
	return os.Symlink(filepath.Join(t.path, oldname), filepath.Join(t.path, newname))
}

// dirFS copies os.DirFS but prevents reading links to files outside the FS.
type dirFS string

func (dir dirFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrInvalid}
	}
	path := string(dir) + "/" + name
	info, err := os.Lstat(path)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: name, Err: err}
	}
	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		lpath, err := os.Readlink(path)
		if err != nil {
			return nil, &os.PathError{Op: "open", Path: name, Err: err}
		}
		if !strings.HasPrefix(lpath, string(dir)) {
			return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrInvalid}
		}
	}
	return os.Open(path)
}
