package cli_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/cardil/kn-event/internal/cli"
	"github.com/cardil/kn-event/internal/event"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestSendInCli(t *testing.T) {
	outputModes := []cli.OutputMode{cli.HumanReadable, cli.JSON, cli.YAML}
	for _, mode := range outputModes {
		t.Run(fmt.Sprint("OutputMode==", mode), func(t *testing.T) {
			assertWithOutputMode(t, mode)
		})
	}
}

func assertWithOutputMode(t *testing.T, mode cli.OutputMode) {
	var buf bytes.Buffer
	sender := &mockSender{}
	ce := cloudevents.NewEvent(cloudevents.VersionV1)
	ce.SetID("543")
	ce.SetType("type")
	ce.SetSource("source")
	target := &cli.TargetArgs{
		URL: "http://example.org",
	}
	opts := &cli.OptionsArgs{
		Output:    mode,
		OutWriter: &buf,
	}
	err := withSender(sender, func() error {
		return cli.Send(ce, target, opts)
	})
	assert.NoError(t, err)
	assert.NotNil(t, ce)
	assert.Equal(t, "543", ce.ID())
	out := buf.String()
	assert.Equal(t, 1, len(sender.sent))

	assert.Contains(t, out, "Event (ID: 543) have been sent.")
}

type mockSender struct {
	sent []cloudevents.Event
}

func (m *mockSender) Send(ce cloudevents.Event) error {
	m.sent = append(m.sent, ce)
	return nil
}

func withSender(sender event.Sender, body func() error) error {
	old := event.SenderFactory
	defer func() {
		event.SenderFactory = old
	}()
	event.SenderFactory = func(target *event.Target) (event.Sender, error) {
		return sender, nil
	}
	return body()
}
