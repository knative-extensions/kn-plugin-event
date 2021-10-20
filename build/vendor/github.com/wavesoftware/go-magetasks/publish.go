package magetasks

import (
	"context"
	"errors"
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/artifact"
	"github.com/wavesoftware/go-magetasks/pkg/tasks"
)

// ErrNoPublisherForArtifact when no publisher for artifact is found.
var ErrNoPublisherForArtifact = errors.New("no publisher for artifact found")

// Publish will publish built artifacts to remote site.
func Publish() {
	mg.Deps(Build)
	artifacts := config.Actual().Artifacts
	t := tasks.Start("📤", "Publishing", len(artifacts) > 0)
	for _, art := range artifacts {
		p := t.Part(fmt.Sprintf("%s %s", art.GetType(), art.GetName()))
		pp := p.Starting()

		publishArtifact(art, pp)
	}
	t.End()
}

func publishArtifact(art config.Artifact, pp tasks.PartProcessing) {
	found := false
	for _, publisher := range config.Actual().Publishers {
		if !publisher.Accepts(art) {
			continue
		}
		found = true
		result := publisher.Publish(art, pp)
		if result.Failed() {
			pp.Done(result.Error)
			return
		}
		config.WithContext(func(ctx context.Context) context.Context {
			return context.WithValue(ctx, artifact.PublishKey(art), result)
		})
	}
	var err error
	if !found {
		err = ErrNoPublisherForArtifact
	}
	pp.Done(err)
}
