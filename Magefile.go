// +build mage

package main

import (
	"os"

	"github.com/fatih/color"
	// mage:import
	"github.com/wavesoftware/go-magetasks"
	"github.com/wavesoftware/go-magetasks/config"
)

// Default target is set to binary
var Default = magetasks.Binary

func init() {
	if val, envset := os.LookupEnv("FORCE_COLOR"); envset && val == "true" {
		color.NoColor = false
	}
	config.Binaries = append(config.Binaries, config.Binary{
		Name: "kn-event",
	})
	config.VersionVariablePath = "github.com/cardil/kn-event/internal.Version"
}
