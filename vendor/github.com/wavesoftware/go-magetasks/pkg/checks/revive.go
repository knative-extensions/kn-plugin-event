package checks

import (
	"path"

	"github.com/magefile/mage/sh"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/files"
)

// Revive will configure revive in the build.
func Revive() config.Task {
	return config.Task{
		Name:      "revive",
		Operation: revive,
		Overrides: []config.Configurator{
			config.NewDependencies("github.com/mgechev/revive"),
		},
	}
}

func revive(notifier config.Notifier) error {
	configFile := ".revive.toml"
	c := path.Join(files.ProjectDir(), configFile)
	if files.DontExists(c) {
		skipBecauseOfMissingConfig(notifier, configFile)
		return nil
	}
	return sh.RunV("revive", "-config", configFile,
		"-formatter", "stylish", "./...")
}
