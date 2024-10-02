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

package cli_test

import (
	"testing"

	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
	"gotest.tools/v3/assert"
	outlogging "knative.dev/client/pkg/output/logging"
	"knative.dev/kn-plugin-event/pkg/cli"
)

func TestSetupContext(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.SetContext(cli.InitialContext())
	cli.SetupOutput(cmd, cli.SimplifiedLoggingSetup(zapcore.InvalidLevel))
	cli.SetupOutput(cmd, cli.DefaultLoggingSetup(zapcore.InvalidLevel))
	ctx := cmd.Context()
	assert.Equal(t, zapcore.InvalidLevel, outlogging.LogLevelFromContext(ctx))
}
