// +build mage

/*
Copyright 2021 The Knative Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"knative.dev/kn-plugin-event/internal"

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
	config.VersionVariablePath = "knative.dev/kn-plugin-event/internal.Version"
	checks.GolangCiLintWithOptions(checks.GolangCiLintOptions{
		New: true,
		Fix: true,
	})
}
