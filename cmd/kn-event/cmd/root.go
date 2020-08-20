package cmd

import (
	"fmt"
	"hash/crc32"
	"os"

	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
)

// OutputMode is type of output to produce
type OutputMode enumflag.Flag

// OutputMode enumeration values.
const (
	HumanReadable OutputMode = iota
	JSON
	YAML
)

var outputModeIds = map[OutputMode][]string{
	HumanReadable: {"human"},
	JSON:          {"json"},
	YAML:          {"yaml"},
}

var (
	// Verbose tells does commands should display additional information about
	// what's happening? Verbose information is printed on stderr.
	Verbose bool
	// Output define type of output commands should be producing
	Output = HumanReadable

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
		fmt.Println(err)
		os.Exit(exitCode(err))
	}
}

func exitCode(err error) int {
	return int(crc32.ChecksumIEEE([]byte(err.Error())))%254 + 1
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(
		&Verbose, "verbose", "v",
		false, "verbose output",
	)
	rootCmd.PersistentFlags().VarP(
		enumflag.New(&Output, "output", outputModeIds, enumflag.EnumCaseInsensitive),
		"output", "o",
		"Output format. One of: human|json|yaml.",
	)

	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(versionCmd)
}
