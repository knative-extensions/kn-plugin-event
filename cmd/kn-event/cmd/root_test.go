package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootInvalidCommand(t *testing.T) {
	called := false
	exitFunc = func(code int) {
		t.Logf("exit code received: %d", code)
		called = true
	}
	rootCmd.SetArgs([]string{"invalid-command"})
	buf := bytes.NewBuffer([]byte{})
	rootCmd.SetOut(buf)
	Execute()

	assert.True(t, called)
}
