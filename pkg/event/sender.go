package event

import (
	"context"
	"errors"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/pkg/logging"
)

// ErrCantSentEvent if event can't be sent.
var ErrCantSentEvent = errors.New("can't sent event")

// Sender will send event to specified target.
type Sender interface {
	// Send will send cloudevents.Event to configured target, or return an error
	// if one occurs.
	Send(ctx context.Context, ce cloudevents.Event) error
}

// CreateSender creates a Sender.
type CreateSender func(cfg *k8s.Configurator, target *Target) (Sender, error)

// Binding holds injectable dependencies.
type Binding struct {
	CreateSender
	k8s.NewKubeClients
}

// NewSender will create a sender that can send event to cluster.
func (b Binding) NewSender(cfg *k8s.Configurator, target *Target) (Sender, error) {
	sndr, err := b.CreateSender(cfg, target)
	if err != nil {
		return nil, err
	}
	return &sendLogic{Sender: sndr}, nil
}

type sendLogic struct {
	Sender
}

func (l *sendLogic) Send(ctx context.Context, ce cloudevents.Event) error {
	err := l.Sender.Send(ctx, ce)
	log := logging.FromContext(ctx)
	if err == nil {
		log.Infof("Event (ID: %s) have been sent.", ce.ID())
		return nil
	}
	return cantSentEvent(err)
}

func cantSentEvent(err error) error {
	if errors.Is(err, ErrCantSentEvent) {
		return err
	}
	return fmt.Errorf("%w: %w", ErrCantSentEvent, err)
}
