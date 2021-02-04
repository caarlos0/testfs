# testfs

A simple `fs.FS` which is contained in a test (using `testing.TB`'s `TempDir()`)
and with a few helper methods.

PS: This lib only works on Go 1.16+.

## Example

```go
func TestSomething(t *testing.T) {
	tmpfs := testfs.New(t)
	testfile := "foo/bar/foobar"
	_ = tmpfs.MkdirAll(filepath.Dir(testfile), 0o764)
	_= tmpfs.WriteFile(testfile, []byte("example"), 0o644)

	// you can now use tmpfs as a fs.FS...
	// fs.WalkDir(tmpfs, ".", func(path string, d fs.DirEntry, err error) error { return nil })

	// and read files of course:
	bts, _ := fs.ReadFile(tmpfs, testfile)
	fmt.Println(string(bts))
}
```

## Why

The idea is to able to test code that use a `fs.FS`, without having to,
for example, commit a bunch of files inside `testdata` and without using
in-memory implementation that might not do the same thing as a real FS.

This is a real FS, it only limits itself to a temporary directory and
cleans after itself once the test is done. You also get a couple of helper
methods to create the structure you need.
