package cmd_test

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"knative.dev/kn-plugin-event/cmd/kn-event/cmd"
	"knative.dev/kn-plugin-event/internal/tests"
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
	if !assert.NoError(t, err) {
		return
	}
	out := buf.String()
	assert.Contains(t, out, "Event (ID: 654321) have been sent.")
	assert.NotNil(t, ce)
	assert.Equal(t, "654321", ce.ID())
	payload, err := tests.UnmarshalCloudEventData(ce.Data())
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
