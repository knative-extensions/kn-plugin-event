package internal

import (
	"os"
	"os/exec"
	"path"

	"github.com/wavesoftware/go-ensure"
)

// EnsureBuildDir creates a build directory.
func EnsureBuildDir() {
	d := path.Join(BuildDir(), "bin")
	ensure.NoError(os.MkdirAll(d, os.ModePerm))
}

// DontExists will check if target file dont exist.
func DontExists(file string) bool {
	_, err := os.Stat(file)
	return err != nil && os.IsNotExist(err)
}

// ExecutableAvailable will return true if given executable in available in
// system env.PATH's.
func ExecutableAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
