package cmd

import (
	"fmt"
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
)

// OutputModeIds maps enumeration values to their textual representations
var OutputModeIds = map[OutputMode][]string{
	HumanReadable: {"human"},
	JSON:          {"json"},
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
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(
		&Verbose, "verbose", "v",
		false, "verbose output",
	)
	rootCmd.PersistentFlags().VarP(
		enumflag.New(&Output, "output", OutputModeIds, enumflag.EnumCaseInsensitive),
		"output", "o",
		"Output format. One of: human|json.",
	)

	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(versionCmd)
}
