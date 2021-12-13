package event_test

import (
	"errors"
	"testing"
	"time"

	"gotest.tools/v3/assert"
	"knative.dev/kn-plugin-event/pkg/event"
	"knative.dev/kn-plugin-event/pkg/tests"
)

func TestCreateFromSpec(t *testing.T) {
	id := event.NewID()
	eventType := "org.example.kn.event.ping"
	eventSource := "/k8s/events/ping"
	spec := &event.Spec{
		Type:   eventType,
		ID:     id,
		Source: eventSource,
		Fields: []event.FieldSpec{
			{Path: "person.name", Value: "Chris"},
			{Path: "person.email", Value: "ksuszyns@example.com"},
			{Path: "ping", Value: 123},
			{Path: "active", Value: true},
			{Path: "ref", Value: "321"},
		},
	}
	actual, err := event.CreateFromSpec(spec)
	assert.NilError(t, err)
	assert.Equal(t, eventType, actual.Type())
	assert.Equal(t, id, actual.ID())
	assert.Equal(t, eventSource, actual.Source())
	expectedData := map[string]interface{}{
		"person": map[string]interface{}{
			"name":  "Chris",
			"email": "ksuszyns@example.com",
		},
		"ping":   123.,
		"ref":    "321",
		"active": true,
	}
	actualData, err := tests.UnmarshalCloudEventData(actual.Data())
	assert.NilError(t, err)
	assert.DeepEqual(t, expectedData, actualData)
	assert.Check(t, tests.TimesAlmostEqual(time.Now(), actual.Time()))
}

func TestCreateFromSpecWithInvalidFieldSpec(t *testing.T) {
	spec := &event.Spec{
		Fields: []event.FieldSpec{
			{Path: "person.name", Value: "Chris Suszynski"},
			{Path: "person.name.first", Value: "Chris"},
		},
	}
	_, err := event.CreateFromSpec(spec)
	assert.Check(t, errors.Is(err, event.ErrCantSetField))
	assert.ErrorContains(t, err,
		"\"person.name.first\" path in conflict with value \"Chris Suszynski\"")
}
