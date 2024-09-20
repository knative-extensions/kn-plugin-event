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

	"go.uber.org/zap/zapcore"
	"knative.dev/client/pkg/output"
	outlogging "knative.dev/client/pkg/output/logging"
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

// SetupContext will set the context commonly for all CLIs.
func SetupContext(ctxual Contextual, defaultLogLevel zapcore.Level) {
	ctx := ctxual.Context()
	if ctx == initialCtx {
		// TODO: knative.dev/pkg/signals should allow for resetting the
		//       context for testing purposes.
		ctx = signals.NewContext()
	}
	ctx = output.WithContext(ctx, ctxual)
	ctx = outlogging.WithLogLevel(ctx, defaultLogLevel)
	ctx = outlogging.EnsureLogger(ctx)
	ctxual.SetContext(ctx)
}

var (
	initialCtxKey = struct{}{}         //nolint:gochecknoglobals
	initialCtx    = context.WithValue( //nolint:gochecknoglobals
		context.Background(), initialCtxKey, true)
)
