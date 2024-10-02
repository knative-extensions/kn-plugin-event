package ics_test

import (
	"bytes"
	"context"
	"net/url"
	"strings"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/wavesoftware/go-commandline"
	"gotest.tools/v3/assert"
	internalics "knative.dev/kn-plugin-event/internal/ics"
	"knative.dev/kn-plugin-event/pkg/ics"
	"knative.dev/kn-plugin-event/pkg/tests"
)

func TestApp(t *testing.T) {
	var outBuf, errBuf bytes.Buffer
	opts := []commandline.Option{
		commandline.WithCommand(func(cmd *cobra.Command) {
			cmd.SetOut(&outBuf)
			cmd.SetErr(&errBuf)
			cmd.SetContext(context.TODO())
		}),
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
	assert.NilError(t, err)

	assert.DeepEqual(t, want, *got)
	assert.Check(t, strings.Contains(errBuf.String(), "Event sent"))
	assert.Check(t, strings.Contains(errBuf.String(), id))
	assert.Equal(t, "", outBuf.String())
}
