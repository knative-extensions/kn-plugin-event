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

package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"knative.dev/kn-plugin-event/internal/event"
)

func TestBuildSubCommandWithNoOptions(t *testing.T) {
	performTestsOnBuildSubCommand(t, newCmdArgs("build", "--output", "json"))
}

func TestBuildSubCommandWithComplexOptions(t *testing.T) {
	performTestsOnBuildSubCommand(
		t, newCmdArgs("build",
			"--output", "json",
			"--type", "org.example.ping",
			"--id", "71830",
			"--source", "/api/v1/ping",
			"--field", "person.name=Chris",
			"--field", "person.email=ksuszyns@example.com",
			"--field", "ping=123",
			"--field", "active=true",
			"--raw-field", "ref=321",
		),
		func(e *cloudevents.Event) {
			e.SetType("org.example.ping")
			e.SetID("71830")
			e.SetSource("/api/v1/ping")
			assert.NoError(t, e.SetData(cloudevents.ApplicationJSON, map[string]interface{}{
				"person": map[string]interface{}{
					"name":  "Chris",
					"email": "ksuszyns@example.com",
				},
				"ping":   123,
				"active": true,
				"ref":    "321",
			}))
		},
	)
}

type eventPreparer func(*cloudevents.Event)

type cmdArgs struct {
	args []string
}

func newCmdArgs(args ...string) cmdArgs {
	return cmdArgs{
		args: args,
	}
}

func performTestsOnBuildSubCommand(t *testing.T, cmd cmdArgs, preparers ...eventPreparer) {
	rootCmd.SetArgs(cmd.args)
	buf := bytes.NewBuffer([]byte{})
	rootCmd.SetOut(buf)
	assert.NoError(t, rootCmd.Execute())
	output := buf.Bytes()
	ec := newEventChecks(t)
	for _, preparer := range preparers {
		preparer(ec.event)
	}
	ec.performEventChecks(output)
}

func (ec eventChecks) performEventChecks(out []byte) {
	actual := cloudevents.NewEvent()
	expected := ec.event
	t := ec.t
	assert.NoError(ec.t, json.Unmarshal(out, &actual))

	assert.NoError(t, actual.Validate())
	assert.Equal(t, expected.Type(), actual.Type())
	assert.Equal(t, expected.DataContentType(), actual.DataContentType())
	assert.Equal(t, ec.unmarshalData(expected.Data()), ec.unmarshalData(actual.Data()))
	assert.Equal(t, expected.Source(), actual.Source())
	delta := 1_000_000.
	assert.InDelta(t, expected.Time().UnixNano(), actual.Time().UnixNano(), delta)
}

func (ec eventChecks) unmarshalData(bytes []byte) map[string]interface{} {
	m, err := event.UnmarshalData(bytes)
	assert.NoError(ec.t, err)
	return m
}

func newEventChecks(t *testing.T) eventChecks {
	return eventChecks{t: t, event: event.NewDefault()}
}

type eventChecks struct {
	t     *testing.T
	event *cloudevents.Event
}
