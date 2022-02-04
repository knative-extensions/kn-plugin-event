package logging

import (
	"context"

	"go.uber.org/zap/zaptest"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/signals"
)

// NewContext creates a new context.Context with development logger, and based
// on OS signals. It should be used only in TestMain() functions.
func NewContext() context.Context {
	ctx := signals.NewContext()
	return logging.WithLogger(ctx, contextLogger(ctx))
}

// WithTestLogger adds a test logger to given context.Context.
func WithTestLogger(ctx context.Context, t zaptest.TestingT) context.Context {
	return logging.WithLogger(ctx, testLogger(t))
}
