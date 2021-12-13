package magetasks

import (
	"context"
	"errors"
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/artifact"
	"github.com/wavesoftware/go-magetasks/pkg/files"
	"github.com/wavesoftware/go-magetasks/pkg/tasks"
)

// ErrNoBuilderForArtifact when no builder for artifact is found.
var ErrNoBuilderForArtifact = errors.New("no builder for artifact found")

// Build will build project artifacts, binaries and images.
func Build() {
	mg.Deps(Test, files.EnsureBuildDir)
	t := tasks.Start("🔨", "Building", len(config.Actual().Artifacts) > 0)
	for _, art := range config.Actual().Artifacts {
		p := t.Part(fmt.Sprintf("%s %s", art.GetType(), art.GetName()))
		pp := p.Starting()

		buildArtifact(art, pp)
	}
	t.End()
}

func buildArtifact(art config.Artifact, pp tasks.PartProcessing) {
	found := false
	for _, builder := range config.Actual().Builders {
		if !builder.Accepts(art) {
			continue
		}
		found = true
		result := builder.Build(art, pp)
		if result.Failed() {
			pp.Done(result.Error)
			return
		}
		config.WithContext(func(ctx context.Context) context.Context {
			return context.WithValue(ctx, artifact.BuildKey(art), result)
		})
	}
	var err error
	if !found {
		err = ErrNoBuilderForArtifact
	}
	pp.Done(err)
}
