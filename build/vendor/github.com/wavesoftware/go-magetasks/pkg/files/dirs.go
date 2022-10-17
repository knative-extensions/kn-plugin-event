package files

import (
	"os"
	"path"

	"github.com/magefile/mage/mg"
	"github.com/wavesoftware/go-ensure"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/dotenv"
	"github.com/wavesoftware/go-magetasks/pkg/output"
)

// EnsureBuildDir creates a build directory.
func EnsureBuildDir() {
	mg.Deps(dotenv.Load, output.Setup)
	d := path.Join(BuildDir(), "bin")
	ensure.NoError(os.MkdirAll(d, os.ModePerm))
}

// BuildDir returns project build dir.
func BuildDir() string {
	artifacts := os.Getenv("ARTIFACTS")
	if artifacts != "" {
		return artifacts
	}
	return relativeToProjectRoot(config.Actual().BuildDirPath)
}

// ProjectDir returns project repo directory.
func ProjectDir() string {
	if config.Actual().ProjectDir != "" {
		return config.Actual().ProjectDir
	}
	repoDir, err := os.Getwd()
	ensure.NoError(err)
	return repoDir
}

func relativeToProjectRoot(paths []string) string {
	fullpath := make([]string, len(paths)+1)
	fullpath[0] = ProjectDir()
	for ix, elem := range paths {
		fullpath[ix+1] = elem
	}
	return path.Join(fullpath...)
}
