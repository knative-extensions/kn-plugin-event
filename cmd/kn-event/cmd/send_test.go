package cmd

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/cardil/kn-event/internal/event"
	"github.com/cardil/kn-event/internal/tests"
	"github.com/stretchr/testify/assert"
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
