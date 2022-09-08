package configuration

import (
	"knative.dev/kn-plugin-event/pkg/cli/ics"
	"knative.dev/kn-plugin-event/pkg/system"
)

// CreateIcs creates the configured ics.App to work with.
func CreateIcs(env system.Environment) *ics.App {
	binding := senderBinding()
	return &ics.App{
		Binding:     eventsBinding(binding),
		Environment: env,
	}
}
