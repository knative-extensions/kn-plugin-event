/*
Copyright 2021 The Knative Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cli_test

import (
	"bytes"
	"fmt"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"knative.dev/kn-plugin-event/internal/cli"
	"knative.dev/kn-plugin-event/internal/event"
	"knative.dev/kn-plugin-event/internal/tests"
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
