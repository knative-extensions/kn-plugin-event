package configuration

import (
	"github.com/cardil/kn-event/internal/event"
	"github.com/cardil/kn-event/internal/sender"
)

// ConfigureSender will setup a sender factory.
func ConfigureSender() {
	event.SenderFactory = sender.New
}
