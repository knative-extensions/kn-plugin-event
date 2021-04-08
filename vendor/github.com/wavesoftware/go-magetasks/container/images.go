package container

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/wavesoftware/go-ensure"
	"github.com/wavesoftware/go-magetasks"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/internal"
	"github.com/wavesoftware/go-magetasks/pkg/tasks"
)

// Images builds a OCI images of binaries.
func Images() {
	mg.Deps(magetasks.Binary)

	if len(config.Binaries) > 0 {
		t := tasks.StartMultiline("📦", "Packaging OCI images")
		for _, binary := range config.Binaries {
			p := t.Part(binary.Name)
			cf := containerFile(binary)
			if internal.DontExists(cf) {
				p.Skip(fmt.Sprintf("containerfile %s don't exist", cf))
				continue
			}
			st := p.Starting()
			args := []string{
				"build", "-f", cf, "-t", imageName(binary),
			}
			for key, val := range binary.ImageArgs {
				args = append(args, "--build-arg",
					fmt.Sprintf("%s=%s", key, val))
			}
			args = append(args, ".")
			err := sh.RunV(containerEngine(), args...)
			st.Done(err)
		}
		t.End()
	}
}

func containerEngine() string {
	ce, err := resolveContainerEngine()
	ensure.NoError(err)
	return ce
}
