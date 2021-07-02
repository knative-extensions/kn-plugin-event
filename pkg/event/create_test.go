package event_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	assert.EqualValues(t, expectedData, actualData)
	delta := 1_000_000.
	assert.InDelta(t, time.Now().UnixNano(), actual.Time().UnixNano(), delta)
}

func TestCreateFromSpecWithInvalidFieldSpec(t *testing.T) {
	spec := &event.Spec{
		Fields: []event.FieldSpec{
			{Path: "person.name", Value: "Chris Suszynski"},
			{Path: "person.name.first", Value: "Chris"},
		},
	}
	_, err := event.CreateFromSpec(spec)
	assert.True(t, errors.Is(err, event.ErrCantSetField))
	assert.Contains(t, err.Error(),
		"\"person.name.first\" path in conflict with value \"Chris Suszynski\"")
}
