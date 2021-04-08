package internal

import (
	"github.com/magefile/mage/sh"
	"github.com/wavesoftware/go-magetasks/config"
)

// BuildDeps install build dependencies.
func BuildDeps() error {
	for _, dep := range config.Dependencies {
		err := sh.RunWith(map[string]string{"GO111MODULE": "off"}, "go", "get", dep)
		if err != nil {
			return err
		}
	}
	return nil
}
