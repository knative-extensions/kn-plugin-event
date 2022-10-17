package checks

import (
	"fmt"
	"path"

	"github.com/magefile/mage/sh"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/files"
)

// Staticcheck will configure staticcheck in the build.
func Staticcheck() config.Task {
	return config.Task{
		Name:      "staticcheck",
		Operation: staticcheck,
		Overrides: []config.Configurator{
			config.NewDependencies("honnef.co/go/tools/cmd/staticcheck@latest"),
		},
	}
}

func staticcheck(notifier config.Notifier) error {
	configFile := "staticcheck.conf"
	c := path.Join(files.ProjectDir(), configFile)
	if files.DontExists(c) {
		skipBecauseOfMissingConfig(notifier, configFile)
		return nil
	}
	cmd := fmt.Sprintf("%s/tools/staticcheck", files.BuildDir())
	return sh.RunV(cmd, "-f", "stylish", "./...")
}
