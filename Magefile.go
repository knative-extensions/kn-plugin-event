// +build mage

package main

import (
	// mage:import
	"github.com/wavesoftware/go-magetasks"
	"github.com/wavesoftware/go-magetasks/config"
)

// Default target is set to Binary
//goland:noinspection GoUnusedGlobalVariable
var Default = magetasks.Binary

func init() {
	config.Binaries = append(config.Binaries, config.Binary{
		Name: "kn-event",
	})
	config.VersionVariablePath = "github.com/cardil/kn-event/internal.Version"
}
