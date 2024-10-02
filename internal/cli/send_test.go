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
	"net/url"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/wavesoftware/go-commandline"
	"gotest.tools/v3/assert"
	"knative.dev/kn-plugin-event/pkg/tests"
)

func TestSendToAddress(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	ce, err := tests.WithCloudEventsServer(func(serverURL url.URL) error {
		return testapp().Execute(
			commandline.WithCommand(func(cmd *cobra.Command) {
				cmd.SetOut(buf)
				cmd.SetErr(buf)
				cmd.SetArgs([]string{
					"send",
					"--to", serverURL.String(),
					"--id", "654321",
					"--field", "person.name=Chris",
					"--field", "person.email=ksuszyns@example.com",
					"--field", "ping=123",
					"--field", "active=true",
					"--raw-field", "ref=321",
				})
			}),
		)
	})
	assert.NilError(t, err)
	out := buf.String()
	assert.Check(t, strings.Contains(out, "Event (ID: 654321) have been sent."))
	assert.Check(t, ce != nil)
	assert.Equal(t, "654321", ce.ID())
	payload, err := tests.UnmarshalCloudEventData(ce.Data())
	assert.NilError(t, err)
	assert.DeepEqual(t, map[string]interface{}{
		"person": map[string]interface{}{
			"name":  "Chris",
			"email": "ksuszyns@example.com",
		},
		"ping":   123.,
		"active": true,
		"ref":    "321",
	}, payload)
}
