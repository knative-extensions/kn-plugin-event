package cmd_test

import (
	"bytes"
	"math"
	"testing"

	"github.com/wavesoftware/go-commandline"
	"gotest.tools/v3/assert"
	"knative.dev/kn-plugin-event/internal/cli/cmd"
)

func TestRootInvalidCommand(t *testing.T) {
	retcode := math.MinInt64
	buf := bytes.NewBuffer([]byte{})
	testapp().ExecuteOrDie(
		commandline.WithOutput(buf),
		commandline.WithExit(func(code int) {
			retcode = code
		}),
		commandline.WithArgs("invalid-command"),
	)

	assert.Check(t, retcode != math.MinInt64)
	assert.Check(t, retcode != 0)
}

func testapp() *commandline.App {
	return commandline.New(new(cmd.App))
}
