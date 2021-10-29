package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"knative.dev/kn-plugin-event/pkg/cli"
	"knative.dev/kn-plugin-event/pkg/configuration"
)

// ErrCantBePresented is returned if data can't be presented.
var ErrCantBePresented = errors.New("can't be presented")

type buildCommand struct {
	*Cmd
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
	b.options.OutWriter = cmd.OutOrStdout()
	b.options.ErrWriter = cmd.ErrOrStderr()
	c := configuration.CreateCli()
	ce, err := c.CreateWithArgs(b.event)
	if err != nil {
		return cantBuildEventError(err)
	}
	out, err := c.PresentWith(ce, b.options.Output)
	if err != nil {
		return fmt.Errorf("event %w: %v", ErrCantBePresented, err)
	}
	cmd.Println(out)
	return nil
}

func cantBuildEventError(err error) error {
	if errors.Is(err, cli.ErrCantBuildEvent) {
		return err
	}
	return fmt.Errorf("%w: %v", cli.ErrCantBuildEvent, err)
}
