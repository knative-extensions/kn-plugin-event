package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendNetImplementedYet(t *testing.T) {
	rootCmd.SetArgs([]string{"send"})
	buf := bytes.NewBuffer([]byte{})
	rootCmd.SetOut(buf)
	assert.EqualError(t, rootCmd.Execute(), "not yet implemented")
}
