package configuration

import (
	"knative.dev/kn-plugin-event/pkg/cli"
	"knative.dev/kn-plugin-event/pkg/system"
)

// CreateCli creates the configured cli.App to work with.
func CreateCli(env system.Environment) *cli.App {
	binding := senderBinding()
	return &cli.App{
		Binding:     eventsBinding(binding),
		Environment: env,
	}
}
