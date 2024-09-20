package cli

import (
	"context"
	"errors"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"knative.dev/kn-plugin-event/pkg/event"
)

// Send will send CloudEvent to target.
func (a *App) Send(ctx context.Context, ce cloudevents.Event, tArgs TargetArgs, params *Params) error {
	target, err := a.createTarget(tArgs, params)
	if err != nil {
		return err
	}
	var sender event.Sender
	sender, err = a.NewSender(params.Parse(), target)
	if err != nil {
		return cantSentEvent(err)
	}
	err = sender.Send(ctx, ce)
	if err == nil {
		return nil
	}
	return cantSentEvent(err)
}

func cantSentEvent(err error) error {
	if errors.Is(err, event.ErrCantSentEvent) {
		return err
	}
	return fmt.Errorf("%w: %w", event.ErrCantSentEvent, err)
}
