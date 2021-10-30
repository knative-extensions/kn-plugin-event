package cmd_test

import (
	"bytes"
	"net/url"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
	"knative.dev/kn-plugin-event/internal/cli/cmd"
	"knative.dev/kn-plugin-event/pkg/tests"
)

func TestSendToAddress(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	ce, err := tests.WithCloudEventsServer(func(serverURL url.URL) error {
		tc := cmd.TestingCmd{}
		tc.Args(
			"send",
			"--to-url", serverURL.String(),
			"--id", "654321",
			"--field", "person.name=Chris",
			"--field", "person.email=ksuszyns@example.com",
			"--field", "ping=123",
			"--field", "active=true",
			"--raw-field", "ref=321",
		)
		tc.Out(buf)
		return tc.Execute()
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
