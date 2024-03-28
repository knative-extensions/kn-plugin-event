package cli

import (
	"errors"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"knative.dev/kn-plugin-event/pkg/event"
)

// Send will send CloudEvent to target.
func (a *App) Send(ce cloudevents.Event, target TargetArgs, options *Options) error {
	props, err := options.WithLogger(a)
	if err != nil {
		return err
	}
	t, err := a.createTarget(target, props)
	if err != nil {
		return err
	}
	s, err := a.Binding.NewSender(t)
	if err != nil {
		return cantSentEvent(err)
	}
	err = s.Send(ce)
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
