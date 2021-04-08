package checks

import (
	"path"

	"github.com/magefile/mage/sh"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/internal"
)

// Revive will configure revive in the build.
func Revive() {
	config.Dependencies = append(config.Dependencies, "github.com/mgechev/revive")
	config.Checks = append(config.Checks, config.CustomTask{
		Name: "revive",
		Task: revive,
	})
}

func revive() error {
	c := path.Join(internal.RepoDir(), "revive.toml")
	if internal.DontExists(c) {
		return nil
	}
	return sh.RunV("revive", "-config", "revive.toml",
		"-formatter", "stylish", "./...")
}
