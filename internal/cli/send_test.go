package cli_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/cardil/kn-event/internal/cli"
	"github.com/cardil/kn-event/internal/event"
	"github.com/cardil/kn-event/internal/tests"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
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
	var buf bytes.Buffer
	sender := &tests.Sender{}
	app := cli.App{
		Binding: event.Binding{
			CreateSender: func(target *event.Target) (event.Sender, error) {
				return sender, nil
			},
		},
	}
	err := app.Send(
		want,
		&cli.TargetArgs{URL: "http://example.org"},
		&cli.OptionsArgs{
			Output:    mode,
			OutWriter: &buf,
		},
	)
	assert.NoError(t, err)
	out := buf.String()
	assert.Equal(t, 1, len(sender.Sent))
	assert.Equal(t, want.ID(), sender.Sent[0].ID())

	assert.Contains(t, out,
		fmt.Sprintf("Event (ID: %s) have been sent.", want.ID()),
	)
}
