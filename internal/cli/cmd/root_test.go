package cmd_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"knative.dev/kn-plugin-event/internal/cli/cmd"
)

func TestRootInvalidCommand(t *testing.T) {
	called := false
	c := cmd.TestingCmd{}
	c.Exit(func(code int) {
		t.Logf("exit code received: %d", code)
		called = true
	})
	c.Args("invalid-command")
	buf := bytes.NewBuffer([]byte{})
	c.Out(buf)
	c.ExecuteOrFail()

	assert.True(t, called)
}
