package ics

import (
	"context"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"knative.dev/client/pkg/flags/sink"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/logging"
)

// SendFromEnv will send an event based on a values stored in environmental
// variables.
func (app *App) SendFromEnv(ctx context.Context, cfg *k8s.Configurator) error {
	c, err := app.configure(cfg)
	if err != nil {
		return err
	}
	err = c.sender.Send(ctx, *c.ce)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCantSendWithICS, err)
	}
	log := logging.FromContext(ctx)
	log.Infow("Event sent", zap.String("ce-id", c.ce.ID()))
	return nil
}

func (app *App) configure(cfg *k8s.Configurator) (config, error) {
	args := &Args{
		Sink: "localhost",
	}
	err := envconfig.Process("K", args)
	if err != nil {
		return config{}, fmt.Errorf("%w: %w", ErrCantConfigureICS, err)
	}
	u, err := apis.ParseURL(args.Sink)
	if err != nil {
		return config{}, fmt.Errorf("%w: %w", ErrCantConfigureICS, err)
	}
	target := &event.Target{
		Reference: &sink.Reference{URL: u},
	}
	s, err := app.Binding.CreateSender(cfg, target)
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
