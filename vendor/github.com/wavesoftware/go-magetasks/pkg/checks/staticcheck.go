package checks

import (
	"path"

	"github.com/magefile/mage/sh"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/internal"
)

// Staticcheck will configure staticcheck in the build.
func Staticcheck() {
	config.Dependencies = append(config.Dependencies,
		"honnef.co/go/tools/cmd/staticcheck")
	config.Checks = append(config.Checks, config.CustomTask{
		Name: "staticcheck",
		Task: staticcheck,
	})
}

func staticcheck() error {
	c := path.Join(internal.RepoDir(), "staticcheck.conf")
	if internal.DontExists(c) {
		return nil
	}
	return sh.RunV("staticcheck", "-f", "stylish", "./...")
}
