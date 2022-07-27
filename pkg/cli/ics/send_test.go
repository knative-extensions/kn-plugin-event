package ics_test

import (
	"context"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"gotest.tools/v3/assert"
	"knative.dev/kn-plugin-event/pkg/cli/ics"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/system"
	"knative.dev/kn-plugin-event/pkg/tests"
	"knative.dev/reconciler-test/pkg/logging"
)

func TestSendFromEnv(t *testing.T) {
	want := cloudevents.NewEvent()
	want.SetID("123=456=123")
	want.SetTime(time.Now().UTC())
	want.SetType("example")
	want.SetSource("tests://example")
	kevent, err := ics.Encode(want)
	assert.NilError(t, err)
	sender := &tests.Sender{}
	env := map[string]string{
		"K_SINK":  "http://cosmos.custer.local",
		"K_EVENT": kevent,
	}
	app := ics.App{
		Binding: event.Binding{
			CreateSender: func(target *event.Target) (event.Sender, error) {
				return sender, nil
			},
		},
		Environment: system.WithContext(logging.WithTestLogger(context.TODO(), t), nil),
	}
	err = tests.WithEnviron(env, app.SendFromEnv)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(sender.Sent))
	got := sender.Sent[0]
	assert.DeepEqual(t, want, got)
}
