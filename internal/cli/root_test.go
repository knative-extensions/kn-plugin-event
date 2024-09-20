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
	"bytes"
	"context"
	"math"
	"testing"

	"github.com/spf13/cobra"
	"github.com/wavesoftware/go-commandline"
	"gotest.tools/v3/assert"
	"knative.dev/kn-plugin-event/internal/cli"
)

func TestRootInvalidCommand(t *testing.T) {
	retcode := math.MinInt64
	buf := bytes.NewBuffer([]byte{})
	testapp().ExecuteOrDie(
		commandline.WithCommand(func(cmd *cobra.Command) {
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs([]string{"invalid-command"})
		}),
		commandline.WithExit(func(code int) {
			retcode = code
		}),
	)

	assert.Check(t, retcode != math.MinInt64)
	assert.Check(t, retcode != 0)
}

func testapp() *commandline.App {
	return commandline.New(&wrap{new(cli.App)})
}

type wrap struct {
	delagate commandline.CobraProvider
}

func (w *wrap) Command() *cobra.Command {
	c := w.delagate.Command()
	c.SetContext(context.TODO())
	return c
}
