// Needs to suppress testpackage lint to be able to call main() func.
package main // nolint:testpackage

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"knative.dev/kn-plugin-event/internal/cli/cmd"
)

func TestMainFunc(t *testing.T) {
	tc := cmd.TestingCmd{
		Cmd: mainCmd,
	}
	buf := bytes.NewBuffer([]byte{})
	tc.Out(buf)
	tc.Args("")
	tc.Exit(func(code int) {
		assert.Equal(t, 0, code)
	})

	main()

	out := buf.String()
	assert.Contains(t, out, "Manage CloudEvents from command line")
}
