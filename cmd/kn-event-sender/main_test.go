package main_test

import (
	"net/url"
	"testing"
	"time"

	kes "knative.dev/kn-plugin-event/cmd/kn-event-sender"
	"knative.dev/kn-plugin-event/internal/cli/ics"
	"knative.dev/kn-plugin-event/internal/tests"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestMainSender(t *testing.T) {
	want := cloudevents.NewEvent()
	want.SetID("azxswq")
	want.SetTime(time.Now().UTC())
	want.SetType("example")
	want.SetSource("tests://example")
	kevent, err := ics.Encode(want)
	assert.NoError(t, err)

	got, err := tests.WithCloudEventsServer(func(serverURL url.URL) error {
		env := map[string]string{
			"K_SINK":  serverURL.String(),
			"K_EVENT": kevent,
		}
		return tests.WithEnviron(env, func() error {
			kes.TestMain()
			return nil
		})
	})
	assert.NoError(t, err)

	assert.EqualValues(t, want, *got)
}
