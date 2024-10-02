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
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
	"github.com/wavesoftware/go-commandline"
	"go.uber.org/zap/zapcore"
	outlogging "knative.dev/client/pkg/output/logging"
	"knative.dev/kn-plugin-event/pkg/cli"
	"knative.dev/kn-plugin-event/pkg/metadata"
)

// Options to override the commandline for testing purposes.
var Options []commandline.Option //nolint:gochecknoglobals

type App struct {
	cli.Params
}

func (a *App) Command() *cobra.Command {
	c := &cobra.Command{
		Use:          metadata.PluginUse,
		Aliases:      []string{"kn " + metadata.PluginUse},
		Short:        metadata.PluginDescription,
		Long:         metadata.PluginLongDescription,
		SilenceUsage: true,
	}

	eventArgs := &cli.EventArgs{}
	targetArgs := &cli.TargetArgs{}
	subcommands := []subcommand{
		&buildCommand{App: a, event: eventArgs},
		&sendCommand{App: a, event: eventArgs, target: targetArgs},
		&versionCommand{App: a},
	}
	for _, each := range subcommands {
		c.AddCommand(each.command())
	}
	c.SetContext(cli.InitialContext())
	c.PersistentPreRun = func(cmd *cobra.Command, _ []string) {
		lvl := zapcore.InfoLevel
		if a.Verbose {
			lvl = zapcore.DebugLevel
		}
		cli.SetupOutput(cmd, cli.DefaultLoggingSetup(lvl))
	}
	c.PersistentPostRunE = func(cmd *cobra.Command, _ []string) error {
		closer := outlogging.LogFileCloserFrom(cmd.Context())
		// ensure to close the log file
		return closer()
	}
	a.setGlobalFlags(c)

	return c
}

func (a *App) setGlobalFlags(c *cobra.Command) {
	c.PersistentFlags().BoolVarP(
		&a.Verbose, "verbose", "v",
		false, "verbose output",
	)
	c.PersistentFlags().VarP(
		enumflag.New(&a.OutputMode, "output", outputModeIDs(), enumflag.EnumCaseInsensitive),
		"output", "o",
		"OutputMode format. One of: human|json|yaml.",
	)
}

var _ commandline.CobraProvider = new(App)

func outputModeIDs() map[cli.OutputMode][]string {
	return map[cli.OutputMode][]string{
		cli.HumanReadable: {"human"},
		cli.JSON:          {"json"},
		cli.YAML:          {"yaml"},
	}
}
