package ics

import (
	"errors"

	"knative.dev/kn-plugin-event/pkg/event"
)

var (
	// ErrCouldntEncode is returned when problem occur while trying to encode an
	// event.
	ErrCouldntEncode = errors.New("couldn't encode an event")
	// ErrCouldntDecode is returned when problem occur while trying to decode an
	// event.
	ErrCouldntDecode = errors.New("couldn't decode an event")
	// ErrCantConfigureICS is returned when problem occur while trying to
	// configure ICS sender.
	ErrCantConfigureICS = errors.New("can't configure ICS sender")
	// ErrCantSendWithICS if can't send with ICS sender.
	ErrCantSendWithICS = errors.New("can't send with ICS sender")
)

// Args holds a list of args for in-cluster-sender.
type Args struct {
	Sink        string
	CeOverrides string
	Event       string
}

// App holds an ICS app binding.
type App struct {
	event.Binding
}
