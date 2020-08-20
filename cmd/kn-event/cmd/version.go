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
				fmt.Printf("kn-event version: %s\n", internal.Version)
			case JSON:
				fmt.Printf("{\"name\": \"kn-event\", \"version\": \"%s\"}\n", internal.Version)
			}
		},
	}
)
