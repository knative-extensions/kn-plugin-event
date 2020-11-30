// +build mage

package main

import (
	"github.com/cardil/kn-event/internal"
	// mage:import
	"github.com/wavesoftware/go-magetasks"
	"github.com/wavesoftware/go-magetasks/config"
)

// Default target is set to binary.
var Default = magetasks.Binary

func init() {
	config.Binaries = append(config.Binaries, config.Binary{
		Name: internal.PluginName,
	})
	config.VersionVariablePath = "github.com/cardil/kn-event/internal.Version"
}
