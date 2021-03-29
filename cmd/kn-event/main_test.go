/*
Copyright 2021 The Knative Authors
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

package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"knative.dev/kn-plugin-event/cmd/kn-event/cmd"
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
