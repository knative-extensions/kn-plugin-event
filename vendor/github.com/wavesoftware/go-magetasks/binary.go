package magetasks

import (
	"path"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/internal"
	"github.com/wavesoftware/go-magetasks/pkg/ldflags"
	"github.com/wavesoftware/go-magetasks/pkg/tasks"
)

// Binary will build a binary executable file.
func Binary() {
	mg.Deps(Test, internal.EnsureBuildDir)
	if len(config.Binaries) > 0 {
		t := tasks.Start("🔨", "Building")
		errs := make([]error, 0)
		for _, binary := range config.Binaries {
			args := []string{
				"build",
			}
			args = ldflags.AppendGitVersion(args, t)
			args = append(args, "-o", fullBinaryName(binary), fullBinaryDirectory(binary))
			err := sh.RunV("go", args...)
			errs = append(errs, err)
		}
		t.End(errs...)
	}
}

func fullBinaryName(bin config.Binary) string {
	return path.Join(internal.BuildDir(), "bin", bin.Name)
}

func fullBinaryDirectory(bin config.Binary) string {
	return path.Join(internal.RepoDir(), "cmd", bin.Name)
}
