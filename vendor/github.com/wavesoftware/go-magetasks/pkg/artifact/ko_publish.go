package artifact

import (
	"fmt"
	"log"
	"os"

	"github.com/google/ko/pkg/build"
	"github.com/google/ko/pkg/commands"
	"github.com/google/ko/pkg/commands/options"
	"github.com/google/ko/pkg/publish"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/output/color"
)

const (
	koPublishResult        = "ko.publish.result"
	koDockerRepo           = "KO_DOCKER_REPO"
	magetasksImageBasename = "IMAGE_BASENAME"
)

// KoPublisherConfigurator is used to configure the publish options for KO.
type KoPublisherConfigurator func(*options.PublishOptions) error

// KoPublisher publishes images with Google's KO.
type KoPublisher struct {
	Configurators []KoPublisherConfigurator
}

func (kp KoPublisher) Accepts(artifact config.Artifact) bool {
	_, ok := artifact.(Image)
	return ok
}

func (kp KoPublisher) Publish(artifact config.Artifact, notifier config.Notifier) config.Result {
	image, ok := artifact.(Image)
	if !ok {
		return config.Result{Error: ErrInvalidArtifact}
	}
	buildResult, ok := config.Actual().Context.Value(BuildKey(image)).(config.Result)
	if !ok || buildResult.Failed() {
		return config.Result{Error: fmt.Errorf(
			"%w: can't find successful KO build result", ErrInvalidArtifact)}
	}
	result, ok := buildResult.Info[koBuildResult].(build.Result)
	if !ok {
		return config.Result{Error: fmt.Errorf(
			"%w: can't find successful KO build result", ErrInvalidArtifact)}
	}
	po, err := kp.publishOptions()
	if err != nil {
		return resultErrKoFailed(err)
	}
	ctx := config.Actual().Context
	publisher, err := commands.NewPublisher(po)
	if err != nil {
		return resultErrKoFailed(err)
	}
	ref, err := publisher.Publish(ctx, result, image.GetName())
	if err != nil {
		return resultErrKoFailed(err)
	}
	notifier.Notify(fmt.Sprintf("pushed image: %s", color.Blue(ref)))
	return config.Result{Info: map[string]interface{}{
		koPublishResult: ref,
	}}
}

func (kp KoPublisher) publishOptions() (*options.PublishOptions, error) {
	if v, ok := os.LookupEnv(magetasksImageBasename); ok {
		if _, ok2 := os.LookupEnv(koDockerRepo); !ok2 {
			if err := os.Setenv(koDockerRepo, v); err != nil {
				return nil, err
			}
		}
	}
	opts := &options.PublishOptions{
		BaseImportPaths: true,
		Push:            true,
	}
	if version := config.Actual().Version; version != nil {
		opts.Tags = []string{version.Resolver()}
	}
	if v, ok := os.LookupEnv(koDockerRepo); ok {
		opts.DockerRepo = v
	}
	for _, configurator := range kp.Configurators {
		if err := configurator(opts); err != nil {
			return nil, err
		}
	}
	return opts, nil
}

func closePublisher(publisher publish.Interface) {
	if err := publisher.Close(); err != nil {
		log.Fatal(err)
	}
}
