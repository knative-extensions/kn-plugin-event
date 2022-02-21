//go:build mage
// +build mage

package main

import (
	// mage:import
	"github.com/wavesoftware/go-magetasks"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/config/buildvars"
	"github.com/wavesoftware/go-magetasks/pkg/artifact"
	artifactimage "github.com/wavesoftware/go-magetasks/pkg/artifact/image"
	"github.com/wavesoftware/go-magetasks/pkg/artifact/platform"
	"github.com/wavesoftware/go-magetasks/pkg/checks"
	"github.com/wavesoftware/go-magetasks/pkg/git"
	"github.com/wavesoftware/go-magetasks/pkg/image"
	"github.com/wavesoftware/go-magetasks/pkg/knative"
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
			Name:           "kn-event",
			BuildVariables: cliBuildVariables(sender),
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
		Version: &config.Version{
			Path: metadata.VersionPath(),
			Resolver: knative.NewVersionResolver(
				knative.WithGit(
					git.WithRemote(git.Remote{
						URL: "https://github.com/knative-sandbox/kn-plugin-event.git",
					}),
				),
			),
		},
		Artifacts: []config.Artifact{sender, cli},
		Checks:    []config.Task{checks.GolangCiLint()},
		BuildVariables: map[string]config.Resolver{
			metadata.ImageBasenamePath():          artifactimage.BaseName,
			metadata.ImageBasenameSeparatorPath(): artifactimage.BaseNameSeparator,
		},
		Overrides: overrides.List,
	})
}

func cliBuildVariables(sender artifact.Image) config.BuildVariables {
	return buildvars.Assemble([]buildvars.Operator{
		image.InfluenceableReference{
			Path:        metadata.ImagePath(),
			EnvVariable: "KN_PLUGIN_EVENT_IMAGE_SENDER",
			Image:       sender,
		},
	})
}
