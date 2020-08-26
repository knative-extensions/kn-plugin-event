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
	// Verbose tells does commands should display additional information about
	// what's happening? Verbose information is printed on stderr.
	Verbose bool
	// Output define type of output commands should be producing
	Output = cli.HumanReadable

	rootCmd = &cobra.Command{
		Use:     "event",
		Aliases: []string{"kn event"},
		Short:   "A plugin for operating on CloudEvents",
		Long: `Manage CloudEvents from command line. Perform, easily, tasks like sending,
building, and parsing, all from command line.`,
	}
)

// Execute will execute the application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitFunc(exitCode(err))
	}
}

// SetOut sets output stream to cmd
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
		&Verbose, "verbose", "v",
		false, "verbose output",
	)
	rootCmd.PersistentFlags().VarP(
		enumflag.New(&Output, "output", outputModeIds, enumflag.EnumCaseInsensitive),
		"output", "o",
		"Output format. One of: human|json|yaml.",
	)

	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(versionCmd)
}
