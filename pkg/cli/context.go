/*
 Copyright 2024 The Knative Authors

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package cli

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"knative.dev/client/pkg/output"
	outlogging "knative.dev/client/pkg/output/logging"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/signals"
)

// Contextual represents a contextual entity that also can serve as an
// output.Printer.
type Contextual interface {
	SetContext(ctx context.Context)
	Context() context.Context
	output.Printer
}

// InitialContext returns the initial context object, so it could be set ahead
// of time the setup is called.
func InitialContext() context.Context {
	return initialCtx
}

// LoggingSetup is a func that sets the logging into the context.
type LoggingSetup func(ctx context.Context) context.Context

// DefaultLoggingSetup is the default logging setup.
func DefaultLoggingSetup(logLevel zapcore.Level) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		ctx = outlogging.WithLogLevel(ctx, logLevel)
		return outlogging.EnsureLogger(ctx)
	}
}

// SimplifiedLoggingSetup is just a production logger to avoid creating
// additional log files.
//
//	 TODO: Remove this after simplified logging is supported in
//		     knative.dev/client/pkg/output/logging package.
func SimplifiedLoggingSetup(logLevel zapcore.Level) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		prtr := output.PrinterFrom(ctx)
		errout := prtr.ErrOrStderr()
		ec := zap.NewProductionEncoderConfig()
		logger := zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(ec),
			zapcore.AddSync(errout),
			logLevel,
		))
		return logging.WithLogger(ctx, logger.Sugar())
	}
}

// SetupContext will set the context commonly for all CLIs.
func SetupContext(ctxual Contextual, loggingSetup LoggingSetup) {
	ctx := ctxual.Context()
	if ctx == initialCtx {
		// TODO: knative.dev/pkg/signals should allow for resetting the
		//       context for testing purposes.
		ctx = signals.NewContext()
	}
	ctx = output.WithContext(ctx, ctxual)
	ctx = loggingSetup(ctx)
	ctxual.SetContext(ctx)
}

var (
	initialCtxKey = struct{}{}         //nolint:gochecknoglobals
	initialCtx    = context.WithValue( //nolint:gochecknoglobals
		context.Background(), initialCtxKey, true)
)
