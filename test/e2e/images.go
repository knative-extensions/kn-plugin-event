//go:build e2e

package e2e

import (
	"knative.dev/kn-plugin-event/test/images"
	"knative.dev/reconciler-test/pkg/environment"
	"knative.dev/reconciler-test/pkg/eventshub"
)

// ConfigureImages will register packages to be built into test images.
func ConfigureImages() environment.EnvOpts {
	return environment.UnionOpts(
		registerWatholaForwarderImage,
		images.ResolveImages([]images.PackageResolver{
			WatholaForwarderImageFromContext,
			eventshub.ImageFromContext,
		}),
	)
}
