package magetasks

import (
	"github.com/magefile/mage/mg"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/internal"
	"github.com/wavesoftware/go-magetasks/pkg/tasks"
)

// Check will run all lints checks.
func Check() {
	mg.Deps(internal.BuildDeps)
	t := tasks.StartMultiline("🔍", "Checking")
	for _, check := range config.Checks {
		p := t.Part(check.Name)
		ps := p.Starting()
		ps.Done(check.Task())
	}
	t.End(nil)
}
