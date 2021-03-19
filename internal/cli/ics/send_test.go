package ics_test

import (
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"knative.dev/kn-plugin-event/internal/cli/ics"
	"knative.dev/kn-plugin-event/internal/event"
	"knative.dev/kn-plugin-event/internal/tests"
)

func TestSendFromEnv(t *testing.T) {
	want := cloudevents.NewEvent()
	want.SetID("123=456=123")
	want.SetTime(time.Now().UTC())
	want.SetType("example")
	want.SetSource("tests://example")
	kevent, err := ics.Encode(want)
	assert.NoError(t, err)
	sender := &tests.Sender{}
	env := map[string]string{
		"K_SINK":  "http://cosmos.custer.local",
		"K_EVENT": kevent,
	}
	app := ics.App{Binding: event.Binding{
		CreateSender: func(target *event.Target) (event.Sender, error) {
			return sender, nil
		},
	}}
	err = tests.WithEnviron(env, app.SendFromEnv)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(sender.Sent))
	got := sender.Sent[0]
	assert.EqualValues(t, want, got)
}
