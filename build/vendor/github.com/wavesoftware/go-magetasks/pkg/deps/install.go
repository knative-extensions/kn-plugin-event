package deps

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/dotenv"
	"github.com/wavesoftware/go-magetasks/pkg/output"
)

// Install install build dependencies.
func Install() error {
	mg.Deps(dotenv.Load, output.Setup)
	for _, dep := range config.Actual().Dependencies.Installs() {
		err := sh.RunWith(map[string]string{"GO111MODULE": "off"}, "go", "get", dep)
		if err != nil {
			return err
		}
	}
	return nil
}
