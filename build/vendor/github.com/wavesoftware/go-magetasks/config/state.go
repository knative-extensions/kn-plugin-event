package config

import (
	"context"
	"log"
)

var state Configurable //nolint:gochecknoglobals

// Actual returns actual config.
func Actual() Config {
	verifyState()
	cfg := state.Config()
	return *cfg
}

// WithContext allows to use context.Context and to amend it.
func WithContext(mutator func(ctx context.Context) context.Context) {
	verifyState()
	ctx := mutator(state.Config().Context)
	state.Config().Context = ctx
}

// Configure with provided Configurable.
func Configure(c Configurable) {
	state = c
	for _, override := range collectConfigurators(c) {
		override.Configure(state)
	}
}

func verifyState() {
	if state == nil {
		log.Fatal("Not configured yet, execute config.Configure() first!")
	}
}

func collectConfigurators(c Configurable) []Configurator {
	cnfrs := make([]Configurator, 0, len(c.Config().Overrides))
	cnfrs = append(cnfrs, c.Config().Overrides...)
	for _, task := range c.Config().Cleaning {
		cnfrs = append(cnfrs, task.Overrides...)
	}
	for _, task := range c.Config().Checks {
		cnfrs = append(cnfrs, task.Overrides...)
	}
	return cnfrs
}
