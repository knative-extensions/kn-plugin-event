package configuration

import (
	"knative.dev/kn-plugin-event/internal/cli/ics"
)

// CreateIcs creates the configured ics.App to work with.
func CreateIcs() *ics.App {
	binding := senderBinding()
	return &ics.App{Binding: eventsBinding(binding)}
}
