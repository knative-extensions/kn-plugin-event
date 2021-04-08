package container

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/internal"
	"github.com/wavesoftware/go-magetasks/pkg/tasks"
)

// Publish publishes built images to a remote registry.
func Publish() {
	mg.Deps(Images)

	if len(config.Binaries) > 0 {
		t := tasks.StartMultiline("📤", "Publishing OCI images")
		errs := make([]error, 0)
		for _, binary := range config.Binaries {
			p := t.Part(binary.Name)
			cf := containerFile(binary)
			im := imageName(binary)
			if internal.DontExists(cf) {
				p.Skip(fmt.Sprintf("no container image for %s", im))
				continue
			}
			ps := p.Starting()
			args := []string{"push", im}
			err := sh.RunV(containerEngine(), args...)
			ps.Done(err)
		}
		t.End(errs...)
	}
}
