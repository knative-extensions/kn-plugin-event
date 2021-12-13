package main_test

import (
	"net/url"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"gotest.tools/v3/assert"
	kes "knative.dev/kn-plugin-event/cmd/kn-event-sender"
	"knative.dev/kn-plugin-event/pkg/cli/ics"
	"knative.dev/kn-plugin-event/pkg/tests"
)

func TestMainSender(t *testing.T) {
	want := cloudevents.NewEvent()
	want.SetID("azxswq")
	want.SetTime(time.Now().UTC())
	want.SetType("example")
	want.SetSource("tests://example")
	kevent, err := ics.Encode(want)
	assert.NilError(t, err)

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
	assert.NilError(t, err)

	assert.DeepEqual(t, want, *got)
}
