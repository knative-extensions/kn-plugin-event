package configuration

import (
	"knative.dev/kn-plugin-event/internal/cli"
)

// CreateCli creates the configured cli.App to work with.
func CreateCli() *cli.App {
	binding := senderBinding()
	return &cli.App{Binding: eventsBinding(binding)}
}
