package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"knative.dev/kn-plugin-event/internal/cli"
	"knative.dev/kn-plugin-event/internal/configuration"
)

var (
	// ErrCantBuildEvent is returned if an event can't be built.
	ErrCantBuildEvent = errors.New("can't build event")

	// ErrCantBePresented is returned if data can't be presented.
	ErrCantBePresented = errors.New("can't be presented")
)

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
		return fmt.Errorf("%w: %v", ErrCantBuildEvent, err)
	}
	out, err := c.PresentWith(ce, b.options.Output)
	if err != nil {
		return fmt.Errorf("event %w: %v", ErrCantBePresented, err)
	}
	cmd.Println(out)
	return nil
}
