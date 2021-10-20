package config

import (
	"context"

	"github.com/fatih/color"
)

var (
	// DefaultBuilders is a list of default builders.
	DefaultBuilders = make([]Builder, 0)
	// DefaultPublishers is a list of default publishers.
	DefaultPublishers = make([]Publisher, 0)
)

// FillInDefaultValues in provided config and returns a filled one.
func FillInDefaultValues(cfg Config) Config {
	if len(cfg.BuildDirPath) == 0 {
		cfg.BuildDirPath = []string{"build", "_output"}
	}
	empty := &MageTag{}
	if cfg.MageTag == *empty {
		cfg.MageTag = MageTag{
			Color: color.FgCyan,
			Label: "[MAGE]",
		}
	}
	if cfg.Dependencies == nil {
		cfg.Dependencies = NewDependencies("github.com/kyoh86/richgo")
	}
	if cfg.Context == nil {
		cfg.Context = context.TODO()
	}
	if cfg.Artifacts == nil {
		cfg.Artifacts = make([]Artifact, 0)
	}
	if len(cfg.Builders) == 0 {
		cfg.Builders = append(cfg.Builders, DefaultBuilders...)
	}
	if len(cfg.Publishers) == 0 {
		cfg.Publishers = append(cfg.Publishers, DefaultPublishers...)
	}
	return cfg
}
