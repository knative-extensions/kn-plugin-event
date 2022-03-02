//go:build e2e
// +build e2e

package e2e

import (
	"knative.dev/kn-plugin-event/test/images"
	"knative.dev/reconciler-test/pkg/environment"
)

// ConfigureImages will register packages to be built into test images.
func ConfigureImages(t images.TestingT) {
	environment.RegisterPackage(watholaForwarderPackage)
	images.ResolveImages(t, []string{
		"knative.dev/reconciler-test/cmd/eventshub",
		watholaForwarderPackage,
	})
}
