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
	"os"
	"path"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"github.com/wavesoftware/go-commandline"
	"gotest.tools/v3/assert"
	outlogging "knative.dev/client/pkg/output/logging"
	"knative.dev/kn-plugin-event/internal/cli"
	"knative.dev/kn-plugin-event/pkg/errors"
)

var errFoo = errors.New("foo")

func TestErrorHandler(t *testing.T) {
	t.Setenv("FORCE_COLOR", "yes")
	cmd := &cobra.Command{
		Use:           "example",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(_ *cobra.Command, _ []string) error {
			return errors.Wrap(errFoo, cli.ErrCantBePresented)
		},
	}
	logfilePath := path.Join(os.TempDir(), "log.jsonl")
	const fileMode = 0o600
	logfile, err := os.OpenFile(logfilePath, os.O_CREATE, fileMode)
	require.NoError(t, err)
	defer func(logfile *os.File) {
		_ = logfile.Close()
	}(logfile)
	ctx := outlogging.WithLogFile(context.TODO(), logfile)
	cmd.SetContext(ctx)
	errBuf := bytes.NewBufferString("")
	cmd.SetErr(errBuf)
	var gotCode *int
	opts := append(
		cli.EffectiveOptions(),
		commandline.WithExit(func(code int) {
			gotCode = &code
		}),
	)
	commandline.New(app{cmd}).ExecuteOrDie(opts...)
	assert.Check(t, gotCode != nil)
	assert.Equal(t, *gotCode, 227)
	wantErrorOutput := `
🔥 Error: can't be presented
  └─ caused by: foo

🌟 Hint: The execution logs could help debug the failure.
         Consider, taking a look at the log file: `
	wantErrorOutput += logfilePath + "\n"
	wantErrorOutput = strings.TrimPrefix(wantErrorOutput, "\n")
	assert.Equal(t, wantErrorOutput, errBuf.String())
}

type app struct {
	cmd *cobra.Command
}

func (a app) Command() *cobra.Command {
	return a.cmd
}
