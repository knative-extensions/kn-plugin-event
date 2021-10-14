package magetasks

import (
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/artifact"
)

// Configure will set up a project to be built.
func Configure(cfg config.Config) {
	artifact.ConfigureDefaults()
	cfg = config.FillInDefaultValues(cfg)
	config.Configure(project{
		cfg: &cfg,
	})
}

type project struct {
	cfg *config.Config
}

func (p project) Config() *config.Config {
	return p.cfg
}
