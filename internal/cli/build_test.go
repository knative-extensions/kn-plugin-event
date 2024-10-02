/*
 Copyright 2024 The Knative Authors

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
	"encoding/json"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/spf13/cobra"
	"github.com/wavesoftware/go-commandline"
	"gotest.tools/v3/assert"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/tests"
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
			assert.NilError(t, e.SetData(cloudevents.ApplicationJSON, map[string]interface{}{
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

func performTestsOnBuildSubCommand(t *testing.T, args cmdArgs, preparers ...eventPreparer) {
	t.Helper()
	buf := bytes.NewBuffer([]byte{})
	assert.NilError(t, testapp().Execute(
		commandline.WithCommand(func(cmd *cobra.Command) {
			cmd.SetOut(buf)
			cmd.SetArgs(args.args)
		}),
	))
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
	assert.NilError(ec.t, json.Unmarshal(out, &actual))

	assert.NilError(t, actual.Validate())
	assert.Equal(t, expected.Type(), actual.Type())
	assert.Equal(t, expected.DataContentType(), actual.DataContentType())
	assert.DeepEqual(t, ec.unmarshalData(expected.Data()), ec.unmarshalData(actual.Data()))
	assert.Equal(t, expected.Source(), actual.Source())
	assert.Check(t, tests.TimesAlmostEqual(expected.Time(), actual.Time()))
}

func (ec eventChecks) unmarshalData(bytes []byte) map[string]interface{} {
	m, err := tests.UnmarshalCloudEventData(bytes)
	assert.NilError(ec.t, err)
	return m
}

func newEventChecks(t *testing.T) eventChecks {
	t.Helper()
	return eventChecks{t: t, event: event.NewDefault()}
}

type eventChecks struct {
	t     *testing.T
	event *cloudevents.Event
}
