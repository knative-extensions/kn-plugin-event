package cmd

import (
	"github.com/cardil/kn-event/internal/cli"
	"github.com/spf13/cobra"
)

var buildCmd = func() *cobra.Command {
	c := &cobra.Command{
		Use:   "build",
		Short: "Builds a CloudEvent and print it to stdout",
		RunE: func(cmd *cobra.Command, args []string) error {
			options.OutWriter = cmd.OutOrStdout()
			options.ErrWriter = cmd.ErrOrStderr()
			ce, err := cli.CreateWithArgs(eventArgs)
			if err != nil {
				return err
			}
			out, err := cli.PresentWith(ce, options.Output)
			if err != nil {
				return err
			}
			cmd.Println(out)
			return nil
		},
	}
	addBuilderFlags(c)
	return c
}()
