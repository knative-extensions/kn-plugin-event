package cmd_test

import (
	"bytes"
	"encoding/json"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"gotest.tools/v3/assert"
	"knative.dev/kn-plugin-event/internal/cli/cmd"
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
	c := cmd.TestingCmd{}
	c.Out(buf)
	c.Args(args.args...)
	assert.NilError(t, c.Execute())
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
