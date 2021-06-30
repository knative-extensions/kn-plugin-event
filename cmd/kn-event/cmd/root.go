package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
	"knative.dev/kn-plugin-event/internal/cli"
	"knative.dev/kn-plugin-event/internal/cli/retcode"
	"knative.dev/kn-plugin-event/internal/event"
)

// Cmd represents a command line application entrypoint.
type Cmd struct {
	options *cli.Options
	root    *cobra.Command
	exit    func(code int)
}

// Execute will execute the application.
func (c *Cmd) Execute() {
	if err := c.execute(); err != nil {
		c.exit(retcode.Calc(err))
	}
}

func (c *Cmd) execute() error {
	c.init()
	// cobra.Command should pass our own errors, no need to wrap them.
	return c.root.Execute() //nolint:wrapcheck
}

func (c *Cmd) init() {
	if c.root != nil {
		return
	}
	c.exit = os.Exit
	c.options = &cli.Options{}
	c.root = &cobra.Command{
		Use:     "event",
		Aliases: []string{"kn event"},
		Short:   "A plugin for operating on CloudEvents",
		Long: `Manage CloudEvents from command line. Perform, easily, tasks like sending,
building, and parsing, all from command line.`,
	}
	c.root.SetOut(os.Stdout)
	c.root.SetErr(os.Stderr)
	c.root.PersistentFlags().BoolVarP(
		&c.options.Verbose, "verbose", "v",
		false, "verbose output",
	)
	c.root.PersistentFlags().VarP(
		enumflag.New(&c.options.Output, "output", outputModeIds(), enumflag.EnumCaseInsensitive),
		"output", "o",
		"Output format. One of: human|json|yaml.",
	)

	eventArgs := &cli.EventArgs{}
	targetArgs := &cli.TargetArgs{}
	commands := []subcommand{
		&buildCommand{Cmd: c, event: eventArgs},
		&sendCommand{Cmd: c, event: eventArgs, target: targetArgs},
		&versionCommand{Cmd: c},
	}
	for _, each := range commands {
		c.root.AddCommand(each.command())
	}

	c.root.PersistentFlags().StringVar(
		&c.options.KnConfig, "config", "~/.config/kn/config.yaml",
		"kn configuration file",
	)
	c.root.PersistentFlags().StringVar(
		&c.options.Kubeconfig, "kubeconfig", event.DefaultKubeconfig,
		"kubectl configuration file",
	)
	c.root.PersistentFlags().BoolVar(
		&c.options.LogHTTP, "log-http", false,
		"log http traffic",
	)

	c.root.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		c.options.OutWriter = cmd.OutOrStdout()
		c.options.ErrWriter = cmd.ErrOrStderr()
	}
}

func outputModeIds() map[cli.OutputMode][]string {
	return map[cli.OutputMode][]string{
		cli.HumanReadable: {"human"},
		cli.JSON:          {"json"},
		cli.YAML:          {"yaml"},
	}
}
