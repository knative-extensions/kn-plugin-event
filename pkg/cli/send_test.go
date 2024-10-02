package cli_test

import (
	"context"
	"fmt"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"gotest.tools/v3/assert"
	outlogging "knative.dev/client/pkg/output/logging"
	"knative.dev/kn-plugin-event/pkg/cli"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/k8s"
	"knative.dev/kn-plugin-event/pkg/tests"
)

func TestSendInCli(t *testing.T) {
	outputModes := []cli.OutputMode{cli.HumanReadable, cli.JSON, cli.YAML}
	for i := range outputModes {
		mode := outputModes[i]
		t.Run(fmt.Sprint("OutputMode==", mode), func(t *testing.T) {
			ce := createExampleEvent()
			assertWithOutputMode(t, ce, mode)
		})
	}
}

func createExampleEvent() cloudevents.Event {
	ce := cloudevents.NewEvent(cloudevents.VersionV1)
	ce.SetID("543")
	ce.SetType("type")
	ce.SetSource("source")
	return ce
}

func assertWithOutputMode(t *testing.T, want cloudevents.Event, mode cli.OutputMode) {
	t.Helper()
	c, logs := observer.New(zapcore.DebugLevel)
	log := &outlogging.ZapLogger{SugaredLogger: zap.New(c).Sugar()}
	ctx := outlogging.WithLogger(context.TODO(), log)
	sender := &tests.Sender{}
	app := cli.App{
		Binding: event.Binding{
			CreateSender: func(*k8s.Configurator, *event.Target) (event.Sender, error) {
				return sender, nil
			},
		},
	}
	err := app.Send(
		ctx,
		want,
		cli.TargetArgs{Sink: "https://example.org"},
		&cli.Params{OutputMode: mode},
	)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(sender.Sent))
	assert.Equal(t, want.ID(), sender.Sent[0].ID())

	msg := fmt.Sprintf("Event (ID: %s) have been sent.", want.ID())
	assert.Equal(t, 1, logs.FilterMessage(msg).Len())
}
