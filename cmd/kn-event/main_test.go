// Needs to suppress testpackage lint to be able to call main() func.
package main

import (
	"bytes"
	"math"
	"strings"
	"testing"

	"github.com/wavesoftware/go-commandline"
	"gotest.tools/v3/assert"
	"knative.dev/kn-plugin-event/internal/cli/cmd"
)

func TestMainFunc(t *testing.T) {
	retcode := math.MinInt64
	defer func() {
		cmd.Options = nil
	}()
	var buf bytes.Buffer
	cmd.Options = []commandline.Option{
		commandline.WithExit(func(code int) {
			retcode = code
		}),
		commandline.WithOutput(&buf),
		commandline.WithArgs(""),
	}

	main()

	out := buf.String()
	assert.Check(t, strings.Contains(out, "Manage CloudEvents from command line"))
	assert.Check(t, retcode == math.MinInt64)
}
