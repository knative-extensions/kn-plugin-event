package cmd

import (
	"hash/crc32"
	"io"
	"os"

	"github.com/cardil/kn-event/internal/cli"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
)

var outputModeIds = map[cli.OutputMode][]string{
	cli.HumanReadable: {"human"},
	cli.JSON:          {"json"},
	cli.YAML:          {"yaml"},
}

var (
	options = &cli.OptionsArgs{}

	rootCmd = &cobra.Command{
		Use:     "event",
		Aliases: []string{"kn event"},
		Short:   "A plugin for operating on CloudEvents",
		Long: `Manage CloudEvents from command line. Perform, easily, tasks like sending,
building, and parsing, all from command line.`,
	}
)

// Execute will execute the application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitFunc(exitCode(err))
	}
}

// SetOut sets output stream to cmd.
func SetOut(newOut io.Writer) {
	rootCmd.SetOut(newOut)
}

var exitFunc = os.Exit

func exitCode(err error) int {
	return int(crc32.ChecksumIEEE([]byte(err.Error())))%254 + 1
}

func init() {
	SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	rootCmd.PersistentFlags().BoolVarP(
		&options.Verbose, "verbose", "v",
		false, "verbose output",
	)
	rootCmd.PersistentFlags().VarP(
		enumflag.New(&options.Output, "output", outputModeIds, enumflag.EnumCaseInsensitive),
		"output", "o",
		"Output format. One of: human|json|yaml.",
	)

	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().StringVar(
		&options.KnConfig, "config", "~/.config/kn/config.yaml",
		"kn configuration file",
	)
	rootCmd.PersistentFlags().StringVar(
		&options.Kubeconfig, "kubeconfig", "~/.kube/config",
		"kubectl configuration file",
	)
	rootCmd.PersistentFlags().BoolVar(
		&options.LogHTTP, "log-http", false,
		"log http traffic",
	)

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		options.OutWriter = cmd.OutOrStdout()
		options.ErrWriter = cmd.ErrOrStderr()
	}
}
