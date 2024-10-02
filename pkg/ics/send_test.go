package ics_test

import (
	"context"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"gotest.tools/v3/assert"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/ics"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/tests"
	"knative.dev/pkg/logging"
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
	t.Setenv("K_SINK", "http://cosmos.custer.local")
	t.Setenv("K_EVENT", kevent)
	app := ics.App{
		Binding: event.Binding{
			CreateSender: func(*k8s.Configurator, *event.Target) (event.Sender, error) {
				return sender, nil
			},
		},
	}
	cfg := &k8s.Configurator{}
	log := zaptest.NewLogger(t, zaptest.Level(zap.WarnLevel))
	ctx := logging.WithLogger(context.TODO(), log.Sugar())
	err = app.SendFromEnv(ctx, cfg)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(sender.Sent))
	got := sender.Sent[0]
	assert.DeepEqual(t, want, got)
}
