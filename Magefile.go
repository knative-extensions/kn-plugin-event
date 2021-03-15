// +build mage

package main

import (
	"fmt"
	"log"

	"github.com/cardil/kn-event/internal"
	"github.com/joho/godotenv"

	// mage:import
	"github.com/wavesoftware/go-magetasks"
	"github.com/wavesoftware/go-magetasks/config"

	// mage:import
	_ "github.com/wavesoftware/go-magetasks/container"
	"github.com/wavesoftware/go-magetasks/pkg/checks"
)

// Default target is set to binary.
//goland:noinspection GoUnusedGlobalVariable
var Default = magetasks.Binary // nolint:deadcode,gochecknoglobals

func init() { //nolint:gochecknoinits
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}
	bins := []string{
		internal.PluginName,
		fmt.Sprintf("%s-sender", internal.PluginName),
	}
	for _, bin := range bins {
		config.Binaries = append(config.Binaries, config.Binary{Name: bin})
	}
	config.VersionVariablePath = "github.com/cardil/kn-event/internal.Version"
	checks.GolangCiLintWithOptions(checks.GolangCiLintOptions{
		New: true,
		Fix: true,
	})
}
