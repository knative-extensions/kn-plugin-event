package event_test

import (
	"testing"
	"time"

	"github.com/cardil/kn-event/internal/event"
	"github.com/stretchr/testify/assert"
)

func TestCreateWithArgs(t *testing.T) {
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
	actualData, err := event.UnmarshalData(actual.Data())
	assert.NoError(t, err)
	assert.EqualValues(t, expectedData, actualData)
	delta := 1_000_000.
	assert.InDelta(t, time.Now().UnixNano(), actual.Time().UnixNano(), delta)
}
