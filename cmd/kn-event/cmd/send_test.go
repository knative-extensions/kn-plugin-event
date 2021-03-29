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

package cmd

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"knative.dev/kn-plugin-event/internal/event"
	"knative.dev/kn-plugin-event/internal/tests"
)

func TestSendToAddress(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	ce, err := tests.WithCloudEventsServer(func(serverURL url.URL) error {
		rootCmd.SetArgs([]string{
			"send",
			"--to-url", serverURL.String(),
			"--id", "654321",
			"--field", "person.name=Chris",
			"--field", "person.email=ksuszyns@example.com",
			"--field", "ping=123",
			"--field", "active=true",
			"--raw-field", "ref=321",
		})
		rootCmd.SetOut(buf)
		return rootCmd.Execute()
	})
	assert.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "Event (ID: 654321) have been sent.")
	assert.NotNil(t, ce)
	assert.Equal(t, "654321", ce.ID())
	payload, err := event.UnmarshalData(ce.Data())
	assert.NoError(t, err)
	assert.EqualValues(t, map[string]interface{}{
		"person": map[string]interface{}{
			"name":  "Chris",
			"email": "ksuszyns@example.com",
		},
		"ping":   123.,
		"active": true,
		"ref":    "321",
	}, payload)
}
