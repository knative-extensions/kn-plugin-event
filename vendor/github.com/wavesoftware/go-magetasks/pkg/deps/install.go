package deps

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/dotenv"
	"github.com/wavesoftware/go-magetasks/pkg/files"
	"github.com/wavesoftware/go-magetasks/pkg/output"
)

// Install install build dependencies.
func Install() error {
	mg.Deps(dotenv.Load, output.Setup, files.EnsureBuildDir)
	for _, dep := range config.Actual().Dependencies.Installs() {
		env := map[string]string{
			"GOBIN": fmt.Sprintf("%s/tools", files.BuildDir()),
		}
		err := sh.RunWith(env, "go", "install", dep)
		if err != nil {
			return err
		}
	}
	return nil
}
