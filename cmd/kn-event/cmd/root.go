package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:     "event",
		Aliases: []string{"kn event"},
		Short:   "A plugin for Knative CLI (kn) that operates CloudEvents",
		Long: `Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet.
Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet.
Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet.`,
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
	rootCmd.AddCommand(sendCmd)
}
