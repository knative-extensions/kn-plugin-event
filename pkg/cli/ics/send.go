package ics

import (
	"fmt"
	"net/url"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/pkg/logging"
)

// SendFromEnv will send an event based on a values stored in environmental
// variables.
func (app *App) SendFromEnv() error {
	c, err := app.configure()
	if err != nil {
		return err
	}
	err = c.sender.Send(*c.ce)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCantSendWithICS, err)
	}
	log := logging.FromContext(app.Context())
	log.Infow("Event sent", zap.String("ce-id", c.ce.ID()))
	return nil
}

func (app *App) configure() (config, error) {
	args := &Args{
		Sink: "localhost",
	}
	err := envconfig.Process("K", args)
	if err != nil {
		return config{}, fmt.Errorf("%w: %w", ErrCantConfigureICS, err)
	}
	u, err := url.Parse(args.Sink)
	if err != nil {
		return config{}, fmt.Errorf("%w: %w", ErrCantConfigureICS, err)
	}
	target := &event.Target{
		Type:   event.TargetTypeReachable,
		URLVal: u,
	}
	s, err := app.Binding.CreateSender(target)
	if err != nil {
		return config{}, fmt.Errorf("%w: %w", ErrCantConfigureICS, err)
	}
	ce, err := Decode(args.Event)
	if err != nil {
		return config{}, err
	}
	return config{
		sender: s,
		ce:     ce,
	}, nil
}

type config struct {
	sender event.Sender
	ce     *cloudevents.Event
}
