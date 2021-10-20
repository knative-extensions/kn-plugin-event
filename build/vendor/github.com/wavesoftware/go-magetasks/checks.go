package magetasks

import (
	"github.com/magefile/mage/mg"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/deps"
	"github.com/wavesoftware/go-magetasks/pkg/tasks"
)

// Check will run all lints checks.
func Check() {
	mg.Deps(deps.Install)
	t := tasks.Start("🔍", "Checking", len(config.Actual().Checks) > 0)
	for _, check := range config.Actual().Checks {
		p := t.Part(check.Name)
		pp := p.Starting()
		pp.Done(check.Operation(pp))
	}
	t.End(nil)
}
