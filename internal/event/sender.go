package event

import (
	"errors"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// ErrCantSentEvent if event can't be sent.
var ErrCantSentEvent = errors.New("can't sent event")

// NewSender will create a sender that can send event to cluster.
func (b Binding) NewSender(target *Target) (Sender, error) {
	sender, err := b.CreateSender(target)
	if err != nil {
		return nil, err
	}
	return &sendLogic{Sender: sender, Properties: target.Properties}, nil
}

type sendLogic struct {
	Sender
	*Properties
}

func (l *sendLogic) Send(ce cloudevents.Event) error {
	err := l.Sender.Send(ce)
	if err == nil {
		l.Log.Infof("Event (ID: %s) have been sent.", ce.ID())
		return nil
	}
	return cantSentEvent(err)
}

func cantSentEvent(err error) error {
	if errors.Is(err, ErrCantSentEvent) {
		return err
	}
	return fmt.Errorf("%w: %v", ErrCantSentEvent, err)
}
