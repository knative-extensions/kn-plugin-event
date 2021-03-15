package configuration

import (
	"github.com/cardil/kn-event/internal/cli/ics"
)

// CreateIcs creates the configured ics.App to work with.
func CreateIcs() *ics.App {
	binding := senderBinding()
	return &ics.App{Binding: eventsBinding(binding)}
}
