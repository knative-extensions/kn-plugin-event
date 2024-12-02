//go:build mage

package main

import (
	"os"

	"knative.dev/kn-plugin-event/build/overrides"
	"knative.dev/kn-plugin-event/pkg/metadata"

	// mage:import
	"knative.dev/toolbox/magetasks"
	"knative.dev/toolbox/magetasks/config"
	"knative.dev/toolbox/magetasks/config/buildvars"
	"knative.dev/toolbox/magetasks/pkg/artifact"
	"knative.dev/toolbox/magetasks/pkg/artifact/platform"
	"knative.dev/toolbox/magetasks/pkg/checks"
	"knative.dev/toolbox/magetasks/pkg/git"
	"knative.dev/toolbox/magetasks/pkg/image"
	"knative.dev/toolbox/magetasks/pkg/knative"
)

// Default target is set to binary.
//
// goland:noinspection GoUnusedGlobalVariable
var Default = magetasks.Build //nolint:deadcode,gochecknoglobals

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
						URL: "https://github.com/knative-extensions/kn-plugin-event.git",
					}),
				),
			),
		},
		Artifacts: []config.Artifact{sender, cli},
		Checks: []config.Task{checks.GolangCiLint(func(o *checks.GolangCiLintOptions) {
			o.Version = "v1.62.2"
		})},
		BuildVariables: map[string]config.Resolver{
			metadata.ImageBasenamePath(): imageBasenameFromEnv,
		},
		Overrides: overrides.List,
	})
}

func cliBuildVariables(sender artifact.Image) config.BuildVariables {
	return buildvars.Assemble([]buildvars.Operator{
		image.InfluenceableReference{
			Path:        metadata.ImagePath(),
			EnvVariable: "KN_PLUGIN_EVENT_SENDER_IMAGE",
			Image:       sender,
		},
	})
}

func imageBasenameFromEnv() string {
	return env("KO_DOCKER_REPO", "IMAGE_BASENAME")
}

func env(keys ...string) string {
	for _, key := range keys {
		if val, ok := os.LookupEnv(key); ok {
			return val
		}
	}
	return ""
}
