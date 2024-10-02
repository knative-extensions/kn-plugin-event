/*
 Copyright 2024 The Knative Authors

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

package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"knative.dev/client/pkg/flags/sink"
	"knative.dev/kn-plugin-event/pkg/binding"
	"knative.dev/kn-plugin-event/pkg/cli"
	"knative.dev/kn-plugin-event/pkg/event"
)

var (
	// ErrSendTargetValidationFailed is returned if a send target can't pass a
	// validation.
	ErrSendTargetValidationFailed = errors.New("send target validation failed")

	// ErrCantSendEvent is returned if event can't be sent.
	ErrCantSendEvent = errors.New("can't send event")
)

type sendCommand struct {
	target *cli.TargetArgs
	event  *cli.EventArgs
	*App
}

func (s *sendCommand) command() *cobra.Command {
	c := &cobra.Command{
		Use:   "send",
		Short: "Builds and sends a CloudEvent to recipient",
		RunE:  s.run,
	}
	addBuilderFlags(s.event, c)
	c.Flags().StringVarP(
		&s.target.Sink, "to", "r", "",
		sink.Usage("to"),
	)
	c.Flags().StringVar(
		&s.target.AddressableURI, "addressable-uri", "",
		`Specify a relative URI of a target addressable resource. If this
option isn't specified target URL will not be changed.`,
	)
	s.SetGlobalFlags(c.Flags())
	s.SetCommandFlags(c.Flags())
	c.PreRunE = func(*cobra.Command, []string) error {
		err := cli.ValidateTarget(s.target)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrSendTargetValidationFailed, err)
		}
		return nil
	}
	return c
}

func (s *sendCommand) run(cmd *cobra.Command, _ []string) error {
	c := binding.CliApp()
	ce, err := c.CreateWithArgs(s.event)
	if err != nil {
		return cantBuildEventError(err)
	}
	err = c.Send(cmd.Context(), *ce, *s.target, &s.Params)
	if err != nil {
		return cantSentEvent(err)
	}
	return nil
}

func cantSentEvent(err error) error {
	if errors.Is(err, event.ErrCantSentEvent) {
		return err
	}
	return fmt.Errorf("%w: %w", event.ErrCantSentEvent, err)
}
