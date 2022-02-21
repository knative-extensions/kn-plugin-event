package logging

import (
	"context"

	"go.uber.org/zap/zaptest"
	"knative.dev/reconciler-test/pkg/environment"
)

// EnvironmentTestLogger setups the test logger for environment test.
func EnvironmentTestLogger(t zaptest.TestingT) environment.EnvOpts {
	return func(ctx context.Context, env environment.Environment) (context.Context, error) {
		return WithTestLogger(ctx, t), nil
	}
}
