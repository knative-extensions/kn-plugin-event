package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/cardil/kn-event/cmd/kn-event/cmd"
	"github.com/stretchr/testify/assert"
)

func TestMainFunc(t *testing.T) {
	r, w, _ := os.Pipe()
	cmd.SetOut(w)

	main()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		assert.NoError(t, err)
		outC <- buf.String()
	}()

	// back to normal state
	assert.NoError(t, w.Close())
	out := <-outC

	assert.Contains(t, out, "Manage CloudEvents from command line")
}
