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
	"knative.dev/client/pkg/output"
	outlogging "knative.dev/client/pkg/output/logging"
	"knative.dev/kn-plugin-event/pkg/binding"
	"knative.dev/kn-plugin-event/pkg/cli"
)

// ErrCantBePresented is returned if data can't be presented.
var ErrCantBePresented = errors.New("can't be presented")

type buildCommand struct {
	*App
	event *cli.EventArgs
}

func (b *buildCommand) command() *cobra.Command {
	c := &cobra.Command{
		Use:   "build",
		Short: "Builds a CloudEvent and print it to stdout",
		RunE:  b.run,
	}
	addBuilderFlags(b.event, c)
	return c
}

func (b *buildCommand) run(cmd *cobra.Command, _ []string) error {
	c := binding.CliApp()
	ctx := cmd.Context()
	ce, err := c.CreateWithArgs(b.event)
	if err != nil {
		return cantBuildEventError(err)
	}
	outlogging.LoggerFrom(ctx).Debugf("Event: %#v", ce)
	out, err := c.PresentWith(ce, b.OutputMode)
	if err != nil {
		return fmt.Errorf("event %w: %w", ErrCantBePresented, err)
	}
	output.PrinterFrom(ctx).Println(out)
	return nil
}

func cantBuildEventError(err error) error {
	if errors.Is(err, cli.ErrCantBuildEvent) {
		return err
	}
	return fmt.Errorf("%w: %w", cli.ErrCantBuildEvent, err)
}
