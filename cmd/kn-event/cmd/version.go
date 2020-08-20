package cmd

import (
	"fmt"

	"github.com/cardil/kn-event/internal"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the kn event plugin version",
		Run: func(cmd *cobra.Command, args []string) {
			switch Output {
			case HumanReadable:
				fmt.Printf(
					"%s version: %s\n",
					internal.PluginName, internal.Version,
				)
			case JSON:
				fmt.Printf(
					`{"name": "%s", "version": "%s"}`,
					internal.PluginName, internal.Version,
				)
			}
		},
	}
)
