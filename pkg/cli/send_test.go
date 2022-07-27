package cli_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"gotest.tools/v3/assert"
	"knative.dev/kn-plugin-event/pkg/cli"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/system"
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
	var (
		outBuf bytes.Buffer
		errBuf bytes.Buffer
	)
	sender := &tests.Sender{}
	app := cli.App{
		Binding: event.Binding{
			CreateSender: func(target *event.Target) (event.Sender, error) {
				return sender, nil
			},
		},
		Environment: system.WithOutputs(&outBuf, &errBuf, nil),
	}
	err := app.Send(
		want,
		cli.TargetArgs{URL: "http://example.org"},
		&cli.Options{
			Output: mode,
		},
	)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(sender.Sent))
	assert.Equal(t, want.ID(), sender.Sent[0].ID())

	outputs := outBuf.String()
	assert.Check(t, strings.Contains(outputs,
		fmt.Sprintf("Event (ID: %s) have been sent.", want.ID()),
	))
}
