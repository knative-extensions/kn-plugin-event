package overrides

import (
	"fmt"
	"os"
	"path"

	"github.com/magefile/mage/sh"
	"knative.dev/toolbox/magetasks/config"
	"knative.dev/toolbox/magetasks/pkg/artifact"
	"knative.dev/toolbox/magetasks/pkg/ldflags"
	"knative.dev/toolbox/magetasks/pkg/output/color"
)

const prebuiltImageType = "🧩 test-content"

func init() { //nolint:gochecknoinits
	List = append(List, testImages([]testImage{
		"./cmd/kn-event-sender",
		"knative.dev/eventing/test/test_images/wathola-forwarder",
		"knative.dev/reconciler-test/cmd/eventshub",
	}))
}

type testImages []testImage

func (ti testImages) Configure(configurable config.Configurable) {
	if tid := os.Getenv("TEST_IMAGES_DIR"); tid == "" {
		return
	}
	cfg := configurable.Config()
	cfg.Artifacts = make([]config.Artifact, 0, len(ti))
	for _, image := range ti {
		cfg.Artifacts = append(cfg.Artifacts, image)
	}
	cfg.Builders = append(cfg.Builders, testImageBuilder{})
}

type testImage string

func (t testImage) GetType() string {
	return prebuiltImageType
}

func (t testImage) GetName() string {
	return path.Base(t.Path())
}

func (t testImage) Path() string {
	return string(t)
}

type testImageBuilder struct{}

func (p testImageBuilder) Accepts(artifact config.Artifact) bool {
	return artifact.GetType() == prebuiltImageType
}

func (p testImageBuilder) Build(art config.Artifact, notifier config.Notifier) config.Result {
	pi, ok := art.(testImage)
	if !ok {
		return config.Result{Error: artifact.ErrInvalidArtifact}
	}
	name := pi.GetName()
	args := []string{
		"build",
	}
	c := config.Actual()
	if c.Version != nil || len(c.BuildVariables) > 0 {
		builder := ldflags.NewBuilder()
		for key, resolver := range c.BuildVariables {
			builder.Add(key, resolver)
		}
		if c.Version != nil {
			builder.Add(c.Version.Path, c.Version.Resolver.Version)
		}
		args = builder.BuildOnto(args)
	}
	binary := fullBinaryName(name)
	args = append(args, "-o", binary, pi.Path())
	env := map[string]string{"CGO_ENABLED": "0"}
	notifier.Notify(fmt.Sprintf("go build: %s",
		color.Blue(name),
	))
	err := sh.RunWithV(env, "go", args...)
	if err != nil {
		err = fmt.Errorf("%w: %v", artifact.ErrGoBuildFailed, err)
	}
	return config.Result{Error: err, Info: map[string]interface{}{}}
}

func fullBinaryName(name string) string {
	return path.Join(os.Getenv("TEST_IMAGES_DIR"), name)
}
