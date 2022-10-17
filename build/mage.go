//go:build ignored
// +build ignored

package main

import (
	"os"
	"path"
	"runtime"

	"github.com/wavesoftware/go-magetasks/entrypoint"
)

func main() {
	os.Exit(entrypoint.Execute(newContext()))
}

func newContext() entrypoint.Context {
	bd := builddir()
	return entrypoint.Context{
		Directories: entrypoint.Directories{
			BuildDir:   bd,
			ProjectDir: path.Dir(bd),
			CacheDir:   path.Join(bd, "_output"),
		},
	}
}

func builddir() string {
	_, file, _, _ := runtime.Caller(0) //nolint:dogsled
	return path.Dir(file)
}
