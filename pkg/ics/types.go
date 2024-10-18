package ics

import (
	"errors"

	"knative.dev/kn-plugin-event/pkg/event"
)

var (
	// ErrCouldntEncode is returned when the problem occurs while trying to encode
	// an event.
	ErrCouldntEncode = errors.New("couldn't encode an event")
	// ErrCouldntDecode is returned when the problem occurs while trying to
	// decode an event.
	ErrCouldntDecode = errors.New("couldn't decode an event")
	// ErrCantConfigureICS is returned when the problem occurs while trying to
	// configure ICS sender.
	ErrCantConfigureICS = errors.New("can't configure in-cluster sender")
	// ErrICSFailed if the in-cluster sender has failed.
	ErrICSFailed = errors.New("the in-cluster sender failed")
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
