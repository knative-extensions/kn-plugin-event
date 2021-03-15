package ics

import (
	"fmt"
	"net/url"

	"github.com/cardil/kn-event/internal/event"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
)

// SendFromEnv will send an event based on a values stored in environmental
// variables.
func (app *App) SendFromEnv() error {
	c, err := app.configure()
	if err != nil {
		return err
	}
	return c.sender.Send(*c.ce)
}

func (app *App) configure() (config, error) {
	args := &Args{
		Sink: "localhost",
	}
	err := envconfig.Process("K", args)
	if err != nil {
		return config{}, fmt.Errorf("%w: %v", ErrCantConfigureICS, err)
	}
	u, err := url.Parse(args.Sink)
	if err != nil {
		return config{}, fmt.Errorf("%w: %v", ErrCantConfigureICS, err)
	}
	target := &event.Target{
		Type:   event.TargetTypeReachable,
		URLVal: u,
	}
	s, err := app.Binding.CreateSender(target)
	if err != nil {
		return config{}, fmt.Errorf("%w: %v", ErrCantConfigureICS, err)
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
