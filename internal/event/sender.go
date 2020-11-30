package event

import (
	"errors"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

var (
	// ErrSenderFactoryUnset will be returned if sender factory isnt set.
	ErrSenderFactoryUnset = errors.New("sender factory is not set")

	// SenderFactory is a factory function that will create a sender.
	SenderFactory func(target *Target) (Sender, error)
)

// NewSender will create a sender that can send event to cluster.
func NewSender(target *Target, props *Properties) (Sender, error) {
	if SenderFactory == nil {
		return nil, ErrSenderFactoryUnset
	}
	s, err := SenderFactory(target)
	if err != nil {
		return nil, err
	}
	return &sendLogic{Sender: s, Properties: props}, nil
}

type sendLogic struct {
	Sender
	*Properties
}

func (l *sendLogic) Send(ce cloudevents.Event) error {
	err := l.Sender.Send(ce)
	if err == nil {
		l.Log.Infof("Event (ID: %s) have been sent.", ce.ID())
	}
	return err
}
