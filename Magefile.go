//go:build mage
// +build mage

package main

import (
	"os"
	"strings"

	// mage:import
	"github.com/wavesoftware/go-magetasks"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/artifact"
	"github.com/wavesoftware/go-magetasks/pkg/artifact/platform"
	"github.com/wavesoftware/go-magetasks/pkg/checks"
	"github.com/wavesoftware/go-magetasks/pkg/git"
	"knative.dev/kn-plugin-event/overrides"
	"knative.dev/kn-plugin-event/pkg/metadata"
)

// Default target is set to binary.
//goland:noinspection GoUnusedGlobalVariable
var Default = magetasks.Build // nolint:deadcode,gochecknoglobals

func init() { //nolint:gochecknoinits
	sender := artifact.Image{
		Metadata: config.Metadata{Name: "kn-event-sender"},
		Architectures: []platform.Architecture{
			platform.AMD64, platform.ARM64, platform.S390X, platform.PPC64LE,
		},
	}
	cli := artifact.Binary{
		Metadata: config.Metadata{
			Name: "kn-event",
			BuildVariables: config.NewBuildVariablesBuilder().
				ConditionallyAdd(referenceImageByDigest, metadata.ImagePath(), artifact.ImageReferenceOf(sender)).
				Build(),
		},
		Platforms: []artifact.Platform{
			{OS: platform.Linux, Architecture: platform.AMD64},
			{OS: platform.Linux, Architecture: platform.ARM64},
			{OS: platform.Linux, Architecture: platform.PPC64LE},
			{OS: platform.Linux, Architecture: platform.S390X},
			{OS: platform.Mac, Architecture: platform.AMD64},
			{OS: platform.Mac, Architecture: platform.ARM64},
			{OS: platform.Windows, Architecture: platform.AMD64},
		},
	}
	magetasks.Configure(config.Config{
		Version:   &config.Version{Path: metadata.VersionPath(), Resolver: git.Version},
		Artifacts: []config.Artifact{sender, cli},
		Checks:    []config.Task{checks.GolangCiLint()},
		BuildVariables: map[string]config.Resolver{
			metadata.ImageBasenamePath(): env("IMAGE_BASENAME"),
		},
		Overrides: overrides.List,
	})
}

func env(key string) func() string {
	return func() string {
		return os.Getenv(key)
	}
}

func skipImageReference() bool {
	if val, ok := os.LookupEnv("SKIP_IMAGE_REFERENCE"); ok {
		return strings.ToLower(val) == "true"
	}
	return false
}

func referenceImageByDigest() bool {
	return !skipImageReference()
}
