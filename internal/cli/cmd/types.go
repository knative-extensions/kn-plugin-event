package cmd

import "github.com/spf13/cobra"

type subcommand interface {
	command() *cobra.Command
}
