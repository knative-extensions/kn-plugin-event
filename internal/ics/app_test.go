package ics_test

import (
	"bytes"
	"net/url"
	"strings"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/wavesoftware/go-commandline"
	"gotest.tools/v3/assert"
	internalics "knative.dev/kn-plugin-event/internal/ics"
	"knative.dev/kn-plugin-event/pkg/cli/ics"
	"knative.dev/kn-plugin-event/pkg/tests"
)

func TestApp(t *testing.T) {
	var outBuf bytes.Buffer
	opts := []commandline.Option{
		commandline.WithOutput(&outBuf),
	}
	id := uuid.New().String()
	want := cloudevents.NewEvent()
	want.SetID(id)
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
			return commandline.New(internalics.App{}).Execute(opts...)
		})
	})
	out := outBuf.String()
	assert.NilError(t, err)

	assert.DeepEqual(t, want, *got)
	assert.Check(t, strings.Contains(out, "Event sent"))
	assert.Check(t, strings.Contains(out, id))
}
