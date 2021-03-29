/*
Copyright 2021 The Knative Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ics

import (
	"fmt"
	"net/url"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"knative.dev/kn-plugin-event/internal/event"
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
