package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"gotest.tools/v3/assert"
	"knative.dev/kn-plugin-event/pkg/cli"
	"knative.dev/kn-plugin-event/pkg/event"
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

func assertWithOutputMode(tb testing.TB, want cloudevents.Event, mode cli.OutputMode) {
	tb.Helper()
	ctx := context.TODO()
	var buf bytes.Buffer
	sender := &tests.Sender{}
	app := cli.App{
		Binding: event.Binding{
			CreateSender: func(target *event.Target) (event.Sender, error) {
				return sender, nil
			},
		},
	}
	err := app.Send(ctx, want, cli.TargetArgs{URL: "http://example.org"}, &cli.Options{
		Output:    mode,
		OutWriter: &buf,
	})
	assert.NilError(tb, err)
	out := buf.String()
	assert.Equal(tb, 1, len(sender.Sent))
	assert.Equal(tb, want.ID(), sender.Sent[0].ID())

	assert.Check(tb, strings.Contains(out,
		fmt.Sprintf("Event (ID: %s) have been sent.", want.ID()),
	))
}
