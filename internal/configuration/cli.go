package configuration

import (
	"github.com/cardil/kn-event/internal/cli"
)

// CreateCli creates the configured cli.App to work with.
func CreateCli() *cli.App {
	binding := senderBinding()
	return &cli.App{Binding: eventsBinding(binding)}
}
