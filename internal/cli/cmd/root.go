package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
	"github.com/wavesoftware/go-commandline"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // for kubeconfig auth plugins to work correctly see issue #24 .
	"knative.dev/kn-plugin-event/pkg/cli"
	"knative.dev/kn-plugin-event/pkg/metadata"
)

// Options to override the commandline for testing purposes.
var Options []commandline.Option //nolint:gochecknoglobals

type App struct {
	cli.Options
}

func (a *App) Command() *cobra.Command {
	c := &cobra.Command{
		Use:          metadata.PluginUse,
		Aliases:      []string{fmt.Sprintf("kn %s", metadata.PluginUse)},
		Short:        metadata.PluginDescription,
		Long:         metadata.PluginLongDescription,
		SilenceUsage: true,
	}
	c.PersistentFlags().BoolVarP(
		&a.Verbose, "verbose", "v",
		false, "verbose output",
	)
	c.PersistentFlags().VarP(
		enumflag.New(&a.Output, "output", outputModeIds(), enumflag.EnumCaseInsensitive),
		"output", "o",
		"Output format. One of: human|json|yaml.",
	)

	eventArgs := &cli.EventArgs{}
	targetArgs := &cli.TargetArgs{}
	commands := []subcommand{
		&buildCommand{App: a, event: eventArgs},
		&sendCommand{App: a, event: eventArgs, target: targetArgs},
		&versionCommand{App: a},
	}
	for _, each := range commands {
		c.AddCommand(each.command())
	}

	c.PersistentFlags().StringVar(
		&a.KubeconfigOptions.Path, "kubeconfig", "",
		"kubectl configuration file (default: ~/.kube/config)",
	)
	c.PersistentFlags().StringVar(
		&a.KubeconfigOptions.Context, "context", "",
		"name of the kubeconfig context to use",
	)
	c.PersistentFlags().StringVar(
		&a.KubeconfigOptions.Cluster, "cluster", "",
		"name of the kubeconfig cluster to use",
	)

	return c
}

var _ commandline.CobraProvider = new(App)

func outputModeIds() map[cli.OutputMode][]string {
	return map[cli.OutputMode][]string{
		cli.HumanReadable: {"human"},
		cli.JSON:          {"json"},
		cli.YAML:          {"yaml"},
	}
}
