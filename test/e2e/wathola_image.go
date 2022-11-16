//go:build e2e

package e2e

import (
	"context"

	"knative.dev/reconciler-test/pkg/environment"
)

const watholaForwarderPackage = "knative.dev/eventing/test/test_images/wathola-forwarder"

type watholaForwarderImageKey struct{}

// WatholaForwarderImageFromContext gets the wathola forwarder image from context.
func WatholaForwarderImageFromContext(ctx context.Context) string {
	if e, ok := ctx.Value(watholaForwarderImageKey{}).(string); ok {
		return e
	}
	return "ko://" + watholaForwarderPackage
}

// WithCustomWatholaForwarderImage allows you to specify a custom wathola
// forwarder image to be used when invoking watholaForwarder.step.
func WithCustomWatholaForwarderImage(image string) environment.EnvOpts {
	return func(ctx context.Context, env environment.Environment) (context.Context, error) {
		return context.WithValue(ctx, watholaForwarderImageKey{}, image), nil
	}
}

func registerWatholaForwarderImage(
	ctx context.Context,
	env environment.Environment,
) (context.Context, error) {
	pkg := WatholaForwarderImageFromContext(ctx)
	return environment.RegisterPackage(pkg)(ctx, env)
}
